/*
	Copyright 2019 Nokia

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package metrics

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"text/template"
	"time"
	"github.com/nokia/onap-vespa/ves-agent/config"
	"github.com/nokia/onap-vespa/govel"

	"github.com/Masterminds/sprig"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/common/model"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

// CollectorState handles the collector internal state
type CollectorState interface {
	// NextMeasurementIndex return the next event index and increments it
	NextMeasurementIndex() (int64, error)
}

type inMemState struct {
	index int64
}

func (mem *inMemState) NextMeasurementIndex() (int64, error) {
	i := mem.index
	mem.index++
	return i, nil
}

// Collector is an utility to collect metrics from
// a PrometheusServer
type Collector struct {
	state       CollectorState                // Measurement state
	rules       config.MetricRules            // Rules for querying data and building the VES events
	api         v1.API                        // Prometheus API
	max         time.Duration                 // Max collection timeframe duration
	domainAbr   string                        // Domain abbreviation for measurements
	evtCfg      *govel.EventConfiguration    // Generals event configuration
	templates   map[string]*template.Template // Cache for templates from rules (to avoid parsing them each time)
	namingCodes map[string]string             // Cache for VnfcNamingCode from VnfcName
}

// NewCollectorWithState creates a new Prometheus Metrics collector from provided configuration
func NewCollectorWithState(cfg *config.MeasurementConfiguration, evtCfg *govel.EventConfiguration, namingCodes map[string]string, state CollectorState) (*Collector, error) {
	log.Info("Initializing Prometheus Measurement Collector to ", cfg.Prometheus.Address)
	clientCfg := api.Config{
		Address: cfg.Prometheus.Address,
		RoundTripper: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   cfg.Prometheus.Timeout,
				KeepAlive: cfg.Prometheus.KeepAlive,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	client, err := api.NewClient(clientCfg)
	if err != nil {
		return nil, err
	}
	return &Collector{
		api:         v1.NewAPI(client),
		rules:       cfg.Prometheus.Rules,
		max:         cfg.MaxBufferingDuration,
		domainAbr:   cfg.DomainAbbreviation,
		evtCfg:      evtCfg,
		templates:   make(map[string]*template.Template),
		state:       state,
		namingCodes: namingCodes,
	}, nil
}

// NewCollector creates a new Prometheus Metrics collector from provided configuration
func NewCollector(cfg *config.MeasurementConfiguration, evtCfg *govel.EventConfiguration, namingCodes map[string]string) (*Collector, error) {
	return NewCollectorWithState(cfg, evtCfg, namingCodes, &inMemState{index: 0})
}

// adjustCollectionStartTime If there's a maximum buffering timeframe set,
// and if current buffering is higher then update "start" time to get the most recent metrics
// fitting in this max timeframe.
//
// Returns the new start time
func (col *Collector) adjustCollectionStartTime(start time.Time, end time.Time, interval time.Duration) time.Time {
	if col.max > 0 && end.Sub(start) > col.max {
		log.Debugf("Rounding collect timeframe to %s", col.max.String())
		return end.Add(-col.max).Truncate(interval)
	}
	return start
}

func (col *Collector) parseTemplate(s string, strict bool) (*template.Template, error) {
	if col.templates == nil {
		col.templates = make(map[string]*template.Template)
	}
	if tmpl, ok := col.templates[s]; ok {
		return tmpl, nil
	}
	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(s)
	if err != nil {
		return nil, err
	}
	if strict {
		tmpl.Option("missingkey=error")
	}
	col.templates[s] = tmpl
	return tmpl, nil
}

func (col *Collector) execTemplate(s string, data interface{}, strict bool) (string, error) {
	tmpl, err := col.parseTemplate(s, strict)
	if err != nil {
		return "", fmt.Errorf("Bad expression template: %s (%s)", s, err.Error())
	}

	buf := bytes.Buffer{}
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("Cannot execute expression %s (%s)", s, err.Error())
	}
	return buf.String(), nil
}

// Run the metric collection, providing it timing information oin which range to query
func (col *Collector) Run(from, to time.Time, interval time.Duration) (interface{}, error) {
	meas, err := col.CollectMetrics(from, to, interval)
	return meas, err
}

// CollectMetrics perform the metric collection, providing it timing information on which range to query
func (col *Collector) CollectMetrics(from, to time.Time, interval time.Duration) (EventMeasurementSet, error) {
	from = col.adjustCollectionStartTime(from, to, interval)
	// Create new measurement set builder
	VNFName := col.evtCfg.VNFName
	nfNamingCode := col.evtCfg.NfNamingCode
	evtName := col.domainAbr + "_" + nfNamingCode + "_Measurements"
	metrics := NewEventMeasurementSetBuilder(MeasurementFactoryFunc(func(vmID string, timestamp time.Time) (*govel.EventMeasurements, error) {
		if vmID == "" {
			vmID = VNFName
		}
		id, err := col.state.NextMeasurementIndex()
		if err != nil {
			return nil, err
		}
		meas := govel.NewMeasurements(evtName, fmt.Sprintf("Measurements%.10d", id), vmID, interval, timestamp.Add(-interval), timestamp)
		meas.NfNamingCode = nfNamingCode
		meas.NfcNamingCode = col.namingCodes[vmID]
		return meas, nil
	}))

	// Initialize a range with interval (prometheus API)
	rng := v1.Range{
		Start: from,
		End:   to,
		Step:  interval,
	}

	log.Info("Starting metrics collection")
	start := time.Now()
	for _, rule := range col.rules.Metrics {
		// Iterate over rules, query prometheus, convert and collect results
		// into the measurement set builder
		if err := col.collectFromRule(&metrics, rule.WithDefaults(col.rules.DefaultValues), rng); err != nil {
			return nil, err
		}
	}
	log.Infof("Metrics collection completed in %s", time.Since(start))
	// Return the built measurement set
	return metrics.Measurements(), nil
}

func (col *Collector) collectFromRule(metrics *EventMeasurementSetBuilder, rule config.MetricRule, rng v1.Range) error {
	data := map[string]interface{}{"interval": int(rng.Step.Seconds())}
	expr, err := col.execTemplate(rule.Expr, data, true)
	if err != nil {
		return err
	}
	// Query prometheus server
	res, err := col.getMatrix(expr, rng)
	if err != nil {
		return err
	}

	for _, meas := range res {
		labels := map[string]string{}
		for lab, val := range meas.Metric {
			labels[string(lab)] = string(val)
		}
		// data := map[string]interface{}{"labels": labels, "interval": rng.Step}
		data["labels"] = labels
		vnfc, err := col.execTemplate(rule.VMIDLabel, data, true)
		if err != nil {
			return fmt.Errorf("Cannot evaluate vmID: %s", err.Error())
		}
		data["vmId"] = vnfc
		target, err := col.execTemplate(rule.Target, data, true)
		if err != nil || target == "" {
			// return fmt.Errorf("Cannot evaluate target: %s", err.Error())
			// Ignore metrics not having valid a target defined
			log.Warnf("Cannot evaluate target: %s", err.Error())
			continue
		}

		// Create the composite VES metric key
		keys := MeasKeys{}
		lbls := rule.Labels
		if target == "AdditionalObjects" {
			lbls = rule.ObjectKeys
		}
		for _, label := range lbls {

			v, err := col.execTemplate(label.Expr, data, false)
			if err != nil {
				return err
			}
			keys[label.Name] = v
		}
		for _, val := range meas.Values {
			// Insert result into measurement set
			timestamp := val.Timestamp.Time()
			log.Debugf("Got metric %s{%s} => time: %s, VNFC: %s, value: %f", target, keys.String(), timestamp.String(), vnfc, val.Value)
			var err error
			if target == "AdditionalObjects" {
				err = metrics.SetAdditionalObject(vnfc, rule.ObjectName, rule.ObjectInstance, timestamp, float64(val.Value), keys)
			} else {
				err = metrics.Set(target, vnfc, timestamp, float64(val.Value), keys)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (col *Collector) getMatrix(query string, r v1.Range) (model.Matrix, error) {
	log.Debugf("Prometheus query : %s", query)
	result, err := col.api.QueryRange(context.Background(), query, r)
	if err != nil {
		return nil, err
	}
	matrix, ok := result.(model.Matrix)
	if !ok {
		return nil, errors.New("Query result cannot be converted into a matrix")
	}
	return matrix, nil
}

// type GetMatrixResult struct {
// 	Matrix model.Matrix
// 	Err    error
// }

// func (col *Collector) getMatrixAsync(query *config.MetricsQuery, r v1.Range) <-chan GetMatrixResult {
// 	c := make(chan GetMatrixResult)
// 	go func() {
// 		mat, err := col.getMatrix(query, r)
// 		c <- GetMatrixResult{Matrix: mat, Err: err}
// 		close(c)
// 	}()
// 	return c
// }
