package evel

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"ves-agent/config"
	"ves-agent/evel/schema"

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
func NewEvel(collector *config.CollectorConfiguration, event *config.EventConfiguration, cacert string) (*Evel, error) {
	log.Info("Initializing evel")
	var tlsConfig tls.Config
	var httpScheme string

	topic := collector.Topic
	if topic != "" && !strings.HasPrefix(topic, "/") {
		topic = "/" + topic
	}
	if len(cacert) == 0 {
		httpScheme = "http"
		log.Infof("no certificate present: Initializing evel communication to collector without tls ")
		//client = NewVESClient(baseURL, nil, schema.V2841(), event.MaxSize)
	} else {
		caBytes := []byte(cacert)
		rootCa := x509.NewCertPool()
		if !rootCa.AppendCertsFromPEM(caBytes) {
			log.Error("Cannot load root CA. PEM not valid: non tls connection will be configured")
			httpScheme = "http"
			//client = NewVESClient(baseURL, nil, schema.V2841(), event.MaxSize)
			log.Infof("no certificate valid: Initializing evel communication to collector without tls ")
		} else {
			httpScheme = "https"
			tlsConfig.RootCAs = rootCa
			//client = NewVESClient(baseURL, &tlsConfig, schema.V2841(), event.MaxSize)
			log.Infof("Initializing evel communication to collector with tls ")
		}
	}

	// check if fpm_password installed, means that password is encrypted
	fpmPassword := "/usr/bin/fpm-password"
	vesPassword := collector.Password
	vesPassPhrase := collector.PassPhrase
	if _, err := os.Stat(fpmPassword); os.IsNotExist(err) {
		log.Infof("fpm-password not installed: uses clear password")
	} else {
		log.Infof("fpm-password exists: uses encrypted password")
		out, err := exec.Command(fpmPassword, "de", vesPassword, vesPassPhrase).Output()
		if err != nil {
			log.Error("Cannot decrypt ves password. not possible to configure VES client")
			return nil, err
		}
		vesPassword = strings.TrimSuffix(string(out), "\n")
	}

	// baseURL := fmt.Sprintf("http://%s:%s@%s:%d/api/eventListener/v5", user, pass, addr, port)
	baseURL := url.URL{
		Scheme: httpScheme,
		Host:   fmt.Sprintf("%s:%d", collector.FQDN, collector.Port),
		Path:   "/api/eventListener/v5",
		//TODO: User and/or password may be optional ?
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
