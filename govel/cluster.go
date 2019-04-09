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
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// VESCollectorIf is the interface for the VES collector API
type VESCollectorIf interface {
	PostEvent(evt Event) error
	PostBatch(batch Batch) error
	GetMeasurementInterval() time.Duration
	GetHeartbeatInterval() time.Duration
	NotifyMeasurementIntervalChanged(ch chan time.Duration) <-chan time.Duration
	NotifyHeartbeatIntervalChanged(ch chan time.Duration) <-chan time.Duration
}

// Cluster manage switches between collectors:
// primary and backup VES collector
// activ VES collector
type Cluster struct {
	activVES              *Evel
	primaryVES, backupVES *Evel
	maxMissed             int
	retryInterval         time.Duration
	mutex                 sync.RWMutex
}

func (cluster *Cluster) isPrimaryActive() bool {
	return cluster.activVES == cluster.primaryVES
}

// NewCluster initilizes the primary and backup ves collectors.
func NewCluster(prim *CollectorConfiguration, back *CollectorConfiguration, event *EventConfiguration, cacert string) (*Cluster, error) {
	primary, errP := NewEvel(prim, event, cacert)
	var backup *Evel
	var errB error
	if back.FQDN != "" {
		backup, errB = NewEvel(back, event, cacert)
	} else {
		log.Debugf("Ignore empty backup collector.")
	}
	activ := primary
	if errP != nil || primary == nil {
		if errB != nil || backup == nil {
			return &Cluster{}, errors.New("Cannot initialize any of the VES connection")
		}
		log.Warn("Cannot initialize primary VES connection.")
		activ = backup
	} else {
		if errB != nil || backup == nil {
			log.Warn("Cannot initialize backup VES connection.")
		}
	}
	return CreateCluster(activ, primary, backup, event.MaxMissed, event.RetryInterval)
}

// CreateCluster creates cluster from existing collectors.
func CreateCluster(activ, primary, backup *Evel, max int, retry time.Duration) (*Cluster, error) {
	return &Cluster{
		activVES:      activ,
		primaryVES:    primary,
		backupVES:     backup,
		maxMissed:     max,
		retryInterval: retry,
	}, nil
}

// GetMeasurementInterval returns the heartbeat measurement of the activ VES collector
// or 0 if agent's default interval should be used
func (cluster *Cluster) GetMeasurementInterval() time.Duration {
	cluster.mutex.RLock()
	defer cluster.mutex.RUnlock()
	return cluster.activVES.measurementInterval
}

// GetHeartbeatInterval returns the heartbeat interval of the activ VES collector
// or 0 if agent's default interval should be used
func (cluster *Cluster) GetHeartbeatInterval() time.Duration {
	cluster.mutex.RLock()
	defer cluster.mutex.RUnlock()
	return cluster.activVES.heartbeatInterval
}

// NotifyMeasurementIntervalChanged subscribe a channel to receive new measurement interval
// when it changes, from active and backup collector.
// The channel must be buffered or aggressively consumed.
// If the channel cannot be written, it won't receive events (writes are non blocking)
func (cluster *Cluster) NotifyMeasurementIntervalChanged(ch chan time.Duration) <-chan time.Duration {
	cluster.mutex.RLock()
	defer cluster.mutex.RUnlock()
	if cluster.activVES != nil {
		cluster.activVES.NotifyMeasurementIntervalChanged(ch)
	}
	if cluster.backupVES != nil {
		cluster.backupVES.NotifyMeasurementIntervalChanged(ch)
	}
	return ch
}

// NotifyHeartbeatIntervalChanged subscribe a channel to receive new heartbeat interval
// when it changes, from active and backup collector.
// The channel must be buffered or aggressively consumed.
// If the channel cannot be written, it won't receive events (writes are non blocking)
func (cluster *Cluster) NotifyHeartbeatIntervalChanged(ch chan time.Duration) <-chan time.Duration {
	cluster.mutex.RLock()
	defer cluster.mutex.RUnlock()
	if cluster.activVES != nil {
		cluster.activVES.NotifyHeartbeatIntervalChanged(ch)
	}
	if cluster.backupVES != nil {
		cluster.backupVES.NotifyHeartbeatIntervalChanged(ch)
	}
	return ch
}

// PostEvent sends an event to the activ VES collector
func (cluster *Cluster) PostEvent(evt Event) error {
	return cluster.perform("event", func(ves *Evel) error { return ves.PostEvent(evt) })
}

// PostBatch sends a list of events to VES collector in a single
// request using the batch interface
func (cluster *Cluster) PostBatch(batch Batch) error {
	return cluster.perform("batch", func(ves *Evel) error { return ves.PostBatch(batch) })
}

func (cluster *Cluster) perform(info string, f func(ves *Evel) error) error {
	cluster.mutex.RLock()
	defer cluster.mutex.RUnlock()
	var err error
	for nbRetry := 0; nbRetry <= cluster.maxMissed; nbRetry++ {
		if err = f(cluster.activVES); err != nil {
			log.Errorf("Cannot post %s: %s", info, err.Error())
			if nbRetry == cluster.maxMissed {
				log.Errorf("VES collector unreachable, switch.")
				cluster.switchCollector()
			} else {
				log.Infof("Retry post %s in %s", info, cluster.retryInterval.String())
				time.Sleep(cluster.retryInterval)
			}
		} else {
			log.Debugf("Post %s succesfull.", info)
			return nil
		}
	}
	return err
}

// SwitchCollector switch the activ VES server
func (cluster *Cluster) switchCollector() {
	if cluster.isPrimaryActive() || cluster.primaryVES == nil {
		if cluster.backupVES != nil {
			cluster.activVES = cluster.backupVES
			log.Debugf("Use backup collector.")
		} else {
			log.Debugf("No backup collector stay on primary.")
		}
	} else {
		cluster.activVES = cluster.primaryVES
		log.Debugf("Use primary collector.")
	}
}
