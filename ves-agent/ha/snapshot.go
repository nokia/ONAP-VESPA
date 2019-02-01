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

package ha

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/raft"
	log "github.com/sirupsen/logrus"
)

// SnapshotableAgentState represent an agent state that supports
// snapshots and restore from a snapshot
type SnapshotableAgentState interface {
	AgentState
	// Snapshot creates a new state snapshot
	Snapshot() *AgentStateSnapshot
	// Restore state from the provided snapshot
	Restore(*AgentStateSnapshot)
}

// SchedulerStateSnapshot is a snapshiot of a scheduler's state
type SchedulerStateSnapshot struct {
	Interval time.Duration `json:"interval"`
	Next     time.Time     `json:"time"`
}

// AlertInfosStateSnapShot is a snapshot of an alert info
type AlertInfosStateSnapShot struct {
	Sn    int64 `json:"sn"`
	Epoch int64 `json:"epoch"`
}

// AgentStateSnapshot holds a serializable copy of agent state
type AgentStateSnapshot struct {
	MeasIdx      int64                             `json:"meas_idx"`
	HbIdx        int64                             `json:"hb_idx"`
	Schedulers   map[string]SchedulerStateSnapshot `json:"schedulers"`
	FaultIdx     int32                             `json:"fault_idx"`
	AlertInfos   map[int32]AlertInfosStateSnapShot `json:"alertInfos"`
	StorageFault map[string]int32                  `json:"storageFault"`
}

// Persist serialize the snapshot to the given output sink
func (snap *AgentStateSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := json.NewEncoder(sink).Encode(snap); err != nil {
		if err2 := sink.Cancel(); err2 != nil {
			log.Panic(err2)
		}
		return err
	}
	return sink.Close()
}

// Release realeases resources used by snapshot.
// Currently does nothing for AgentStateSnapshot
func (snap *AgentStateSnapshot) Release() {}
