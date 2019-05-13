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

package govel

import (
	"errors"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"github.com/nokia/onap-vespa/govel/schema"

	log "github.com/sirupsen/logrus"
)

type postEventRequest struct {
	Event Event `json:"event"`
}

type postBatchRequest struct {
	EventList Batch `json:"eventList"`
}

// Evel is the main VES collector API entry point
type Evel struct {
	baseURL             url.URL
	topic               string
	heartbeatInterval   time.Duration
	measurementInterval time.Duration
	reportingEntityName string
	reportingEntityID   string
	client              *VESClient
	mutex               sync.RWMutex
	measIntCh           []chan time.Duration
	hbIntCh             []chan time.Duration
}

// NewEvel creates and initialize a new connection to VES collector
func NewEvel(collector *CollectorConfiguration, event *EventConfiguration, cacert string) (*Evel, error) {
	log.Info("Initializing evel")
	var tlsConfig tls.Config
	var httpScheme string

	topic := strings.TrimLeft(collector.Topic, "/")
	if topic != "" && !strings.HasPrefix(topic, "/") {
		topic = "/" + topic
	}
	
	if collector.Secure {
		log.Infof("Secure VES link using HTTPS")
		httpScheme = "https"
	} else {
		log.Warnf("Insecure VES link using HTTP")
		httpScheme = "http"
	} 
	if collector.Secure && len(cacert) > 0 {
		log.Info("Using provided CA certificate")
		caBytes := []byte(cacert)
		rootCa := x509.NewCertPool()
		if !rootCa.AppendCertsFromPEM(caBytes) {
			log.Error("Cannot load root CA. PEM not valid")
		} else {
			tlsConfig.RootCAs = rootCa
		}
	}

	// check if fpm_password installed, means that password is encrypted
	fpmPassword := "/usr/bin/fpm-password"
	vesPassword := collector.Password
	vesPassPhrase := collector.PassPhrase
	if _, err := os.Stat(fpmPassword); os.IsNotExist(err) {
		log.Debugf("Use clear password")
	} else {
		log.Debugf("Use encrypted password")
		out, errFPM := exec.Command(fpmPassword, "de", vesPassword, vesPassPhrase).Output()
		if errFPM != nil {
			log.Warn("Failed to decrypt ves password.")
			return nil, errFPM
		}
		vesPassword = strings.TrimSuffix(string(out), "\n")
		if vesPassword == ""  {
			log.Warn("Failed to decrypt ves password.")
			return nil, errors.New("Failed to decrypt ves password")
		}
	} 
	

	path := strings.TrimLeft(collector.ServerRoot, "/")
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = strings.TrimRight(path, "/") + "/eventListener/v5"
	baseURL := url.URL{
		Scheme: httpScheme,
		Host:   fmt.Sprintf("%s:%d", collector.FQDN, collector.Port),
		Path:   path,
		User: url.UserPassword(collector.User, vesPassword),
	}

	client := NewVESClient(baseURL, &tlsConfig, schema.V2841(), event.MaxSize)
	return &Evel{
		baseURL:             baseURL,
		topic:               topic,
		client:              client,
		reportingEntityName: event.ReportingEntityName,
		reportingEntityID:   event.ReportingEntityID,
		measIntCh:           make([]chan time.Duration, 0),
		hbIntCh:             make([]chan time.Duration, 0),
	}, nil
}

// GetMeasurementInterval returns the measurment interval asked by VES server
// or 0 if agent's default interval should be used
func (evel *Evel) GetMeasurementInterval() time.Duration {
	evel.mutex.RLock()
	defer evel.mutex.RUnlock()
	return evel.measurementInterval
}

// GetHeartbeatInterval returns the heartbeat interval asked by VES server
// or 0 if agent's default interval should be used
func (evel *Evel) GetHeartbeatInterval() time.Duration {
	evel.mutex.RLock()
	defer evel.mutex.RUnlock()
	return evel.heartbeatInterval
}

// NotifyMeasurementIntervalChanged subscribe a channel to receive new measurement interval
// when it changes. The channel must be buffered or aggressively consumed.
// If the channel cannot be written, it won't receive events (writes are non blocking)
func (evel *Evel) NotifyMeasurementIntervalChanged(ch chan time.Duration) <-chan time.Duration {
	evel.mutex.Lock()
	defer evel.mutex.Unlock()
	evel.measIntCh = append(evel.measIntCh, ch)
	return ch
}

// NotifyHeartbeatIntervalChanged subscribe a channel to receive new heartbeat interval
// when it changes. The channel must be buffered or aggressively consumed.
// If the channel cannot be written, it won't receive events (writes are non blocking)
func (evel *Evel) NotifyHeartbeatIntervalChanged(ch chan time.Duration) <-chan time.Duration {
	evel.mutex.Lock()
	defer evel.mutex.Unlock()
	evel.hbIntCh = append(evel.hbIntCh, ch)
	return ch
}

// PostEvent sends an event to VES collector
func (evel *Evel) PostEvent(evt Event) error {
	if evt.Header().ReportingEntityName == "" {
		evt.Header().ReportingEntityName = evel.reportingEntityName
	}
	if evt.Header().ReportingEntityID == "" {
		evt.Header().ReportingEntityID = evel.reportingEntityID
	}

	log.Debugf("Posting event: %+v", evt)
	req := postEventRequest{Event: evt}
	return evel.doPost(evel.topic, req)
}

// PostBatch sends a list of events to VES collector in a single
// request using the batch interface
func (evel *Evel) PostBatch(batch Batch) error {
	if batch.Len() == 0 {
		return nil
	}
	log.Debugf("Posting a batch of events: %#v", batch)
	batch.UpdateReportingEntityName(evel.reportingEntityName)
	batch.UpdateReportingEntityID(evel.reportingEntityID)
	req := postBatchRequest{EventList: batch}
	var err error
	if err = evel.doPost("eventBatch", req); err == ErrBodyTooLarge {
		b1, b2 := batch.Split()
		if b1.Len() == 0 || b2.Len() == 0 {
			log.Panic("Cannot split batch more. Event is bigger than maximum authorized")
		}
		if err = evel.PostBatch(b1); err != nil {
			return err
		}
		return evel.PostBatch(b2)
	}
	return err
}

func (evel *Evel) doPost(queryPath string, req interface{}) error {
	vesResp, err := evel.client.PostJSON(queryPath, req)
	if err != nil {
		return err
	}
	evel.processCommands(vesResp.CommandList)
	return nil
}

func (evel *Evel) processCommands(commandList []Command) {
	if len(commandList) == 0 {
		//Early return to avoid useless locking
		return
	}
	evel.mutex.Lock()
	defer evel.mutex.Unlock()
	for _, command := range commandList {
		switch command.CommandType {
		case CommandHeartbeatIntervalChange:
			heartbeatInterval := time.Duration(command.HeartbeatInterval) * time.Second
			if heartbeatInterval != evel.heartbeatInterval {
				log.Infof("Heartbeat interval changed from %s to %s ", evel.heartbeatInterval.String(), heartbeatInterval.String())
				evel.heartbeatInterval = heartbeatInterval
				for _, ch := range evel.hbIntCh {
					// Non blocking write, to avoid a dead lock situation
					select {
					case ch <- heartbeatInterval:
					default:
						log.Warnf("Heartbeat interval change (to %s) could not be sent to a channel", heartbeatInterval.String())
					}
				}
			}
		case CommandMeasurementIntervalChange:
			measurementInterval := time.Duration(command.MeasurementInterval) * time.Second
			if measurementInterval != evel.measurementInterval {
				log.Infof("Measurement interval changed from %s to %s", evel.measurementInterval.String(), measurementInterval.String())
				evel.measurementInterval = measurementInterval
				for _, ch := range evel.measIntCh {
					// Non blocking write, to avoid a dead lock situation
					select {
					case ch <- measurementInterval:
					default:
						log.Warnf("Measurement interval change (to %s) could not be sent to a channel", measurementInterval.String())
					}

				}
			}
		default:
			log.Warn("Unsupported command type: ", command.CommandType)
		}
	}
}
