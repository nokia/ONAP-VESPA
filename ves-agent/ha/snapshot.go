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
