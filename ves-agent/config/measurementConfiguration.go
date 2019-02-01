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

package config

import (
	"time"
)

// Label represents a VES field by it's name, with an expression
// for getting its value
type Label struct {
	Name string `mapstructure:"name"`
	Expr string `mapstructure:"expr"`
}

// MetricRule defines how to retrieve metrics and map them
// into a list of evel.EventMeasurement struct
type MetricRule struct {
	Target         string  `mapstructure:"target"`          // Target VES event field
	Expr           string  `mapstructure:"expr"`            // Prometheus query expression
	VMIDLabel      string  `mapstructure:"vmId"`            // Metric label holding the VNF ID
	Labels         []Label `mapstructure:"labels"`          // Set of VES fields to map to values of given label
	ObjectName     string  `mapstructure:"object_name"`     // JSON Object Name
	ObjectInstance string  `mapstructure:"object_instance"` // JSON Object instance
	ObjectKeys     []Label `mapstructure:"object_keys"`     // JSON Object keys
}

func (rule MetricRule) hasLabel(name string) bool {
	for _, label := range rule.Labels {
		if label.Name == name {
			return true
		}
	}
	return false
}

// WithDefaults applies default values from `def` to `rule`, and return a new one
func (rule MetricRule) WithDefaults(def *MetricRule) MetricRule {
	if def == nil {
		return rule
	}
	if rule.Target == "" {
		rule.Target = def.Target
	}
	// if rule.Expr == "" {
	// 	rule.Expr = def.Expr
	// }
	if rule.VMIDLabel == "" {
		rule.VMIDLabel = def.VMIDLabel
	}
	labels := make([]Label, len(rule.Labels))
	copy(labels, rule.Labels)
	rule.Labels = labels
	for _, l := range def.Labels {
		if !rule.hasLabel(l.Name) {
			rule.Labels = append(rule.Labels, l)
		}
	}
	return rule
}

// MetricRules defines a list of rules, and defaults values for them
type MetricRules struct {
	DefaultValues *MetricRule  `mapstructure:"defaults"` // Default rules to apply (except for expr), labels are merged
	Metrics       []MetricRule `mapstructure:"metrics"`  // List of query and mapping of rules
}

// PrometheusConfig parameters
type PrometheusConfig struct {
	Address   string        `mapstructure:"address"`   // Base URL to prometheus API
	Timeout   time.Duration `mapstructure:"timeout"`   // API request timeout
	KeepAlive time.Duration `mapstructure:"keepalive"` // HTTP Keep-Alive
	Rules     MetricRules   `mapstructure:"rules"`     // Querying rules
}

// MeasurementConfiguration parameters
type MeasurementConfiguration struct {
	DomainAbbreviation   string           `mapstructure:"domainAbbreviation"`   // "Measurement" or "Mfvs"
	DefaultInterval      time.Duration    `mapstructure:"defaultInterval"`      // Default measurement interval
	MaxBufferingDuration time.Duration    `mapstructure:"maxBufferingDuration"` // Maximum timeframe size of buffering
	Prometheus           PrometheusConfig `mapstructure:"prometheus"`           // Prometheus configuration
}
