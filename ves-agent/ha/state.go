package ha

import (
	"errors"
	"time"
	"github.com/nokia/onap-vespa/ves-agent/convert"
	"github.com/nokia/onap-vespa/ves-agent/heartbeat"
	"github.com/nokia/onap-vespa/ves-agent/metrics"
	"github.com/nokia/onap-vespa/ves-agent/scheduler"

	log "github.com/sirupsen/logrus"
)

// AgentState is the interface that abstract manipulation
// of agent internal state
type AgentState interface {
	scheduler.State
	heartbeat.MonitorState
	metrics.CollectorState
	convert.FaultManagerState
}

type schedulerState struct {
	interval time.Duration
	next     time.Time
}

type inMemState struct {
	measIdx    int64
	hbIdx      int64
	schedulers map[string]*schedulerState
	faultIdx   int32
	alertInfos map[int32]*convert.AlertInfos
	storage    map[string]int32
}

// NewInMemState creates a new snapshotable state stored in memory
func NewInMemState() SnapshotableAgentState {
	return &inMemState{
		schedulers: make(map[string]*schedulerState),
		alertInfos: make(map[int32]*convert.AlertInfos),
		storage:    make(map[string]int32),
	}
}

func (state *inMemState) getOrCreateScheduler(name string) *schedulerState {
	scheduler, ok := state.schedulers[name]
	if !ok {
		scheduler = &schedulerState{}
		state.schedulers[name] = scheduler
	}
	return scheduler
}

func (state *inMemState) NextMeasurementIndex() (int64, error) {
	idx := state.measIdx
	state.measIdx++
	return idx, nil
}

func (state *inMemState) NextHeartbeatIndex() (int64, error) {
	idx := state.hbIdx
	state.hbIdx++
	return idx, nil
}

func (state *inMemState) NextRun(sched string) time.Time {
	sch, ok := state.schedulers[sched]
	if !ok {
		return time.Time{}
	}
	return sch.next
}

func (state *inMemState) UpdateNextRun(sched string, next time.Time) error {
	state.getOrCreateScheduler(sched).next = next
	return nil
}

func (state *inMemState) Interval(sched string) time.Duration {
	sch, ok := state.schedulers[sched]
	if !ok {
		return time.Duration(0)
	}
	return sch.interval
}

func (state *inMemState) UpdateInterval(sched string, interval time.Duration) error {
	state.getOrCreateScheduler(sched).interval = interval
	return nil
}

func (state *inMemState) UpdateScheduler(sched string, interval time.Duration, next time.Time) error {
	schd := state.getOrCreateScheduler(sched)
	schd.interval = interval
	schd.next = next
	return nil
}

// IncrementFaultIdx return the new FaultId to use (FaultManagerState implementation)
func (state *inMemState) NextFaultIndex() (int32, error) {
	state.faultIdx++
	idx := state.faultIdx
	log.Debugf("NextFaultIndex for fault: %d\n", idx)
	return idx, nil
}

func (state *inMemState) InitAlertInfos(faultID int32) error {
	log.Debugf("InitAlertInfos for fault: %010d\n", faultID)
	state.alertInfos[faultID] = &convert.AlertInfos{Sequence: 1, StartEpoch: 0}
	return nil
}

// GetFaultInStorage checks if faultName already associated to an index
func (state *inMemState) GetFaultInStorage(faultName string) int32 {
	if val, ok := state.storage[faultName]; ok {
		return val
	}
	return 0
}

// StoreFaultInStorage stores the index associated to the faultName
func (state *inMemState) StoreFaultInStorage(faultName string, faultID int32) error {
	log.Debugf("state StoreFaultInStorage for fault %s with index %010d", faultName, faultID)
	state.storage[faultName] = faultID
	return nil
}

// DeleteFaultInStorage delete the storage and alertInfos associated to the faultName
func (state *inMemState) DeleteFaultInStorage(faultName string) error {
	log.Debugf("state DeleteFaultInStorage for fault %s", faultName)
	if id, ok := state.storage[faultName]; ok {
		delete(state.alertInfos, id)
	}
	delete(state.storage, faultName)
	return nil
}

// GetFaultSequence return the sequence Number of the faultID index (FaultManagerState implementation)
func (state *inMemState) GetFaultSn(faultID int32) int64 {
	if fault, ok := state.alertInfos[faultID]; ok {
		return fault.Sequence
	}
	return 0
}

// IncrementFaultSn increment the sequence Number of the faultID index (FaultManagerState implementation)
func (state *inMemState) IncrementFaultSn(faultID int32) error {
	log.Debugf("state IncrementFaultSn for fault index %010d", faultID)
	if fault, ok := state.alertInfos[faultID]; ok {
		fault.Sequence++
		return nil
	}
	return errors.New("Fault does not exist")
}

// GetFaultStartEpoch return the startEpoch value of the faultID index (FaultManagerState implementation)
func (state *inMemState) GetFaultStartEpoch(faultID int32) int64 {
	if fault, ok := state.alertInfos[faultID]; ok {
		return fault.StartEpoch
	}
	return 0
}

// SetFaultStartEpoch set the value epoch to the alert faultID (FaultManagerState implementation)
func (state *inMemState) SetFaultStartEpoch(faultID int32, epoch int64) error {
	log.Debugf("state SetFaultStartEpoch for fault: %010d and epoch: %d", faultID, epoch)
	_, ok := state.alertInfos[faultID]
	if !ok {
		log.Errorf("state SetFaultStartEpoch create alertInfos for fault index %010d", faultID)
		if err := state.InitAlertInfos(faultID); err != nil {
			return err
		}
	}
	state.alertInfos[faultID].StartEpoch = epoch
	return nil
}

func (state *inMemState) Snapshot() *AgentStateSnapshot {
	snapshot := new(AgentStateSnapshot)
	snapshot.HbIdx = state.hbIdx
	snapshot.MeasIdx = state.measIdx
	snapshot.FaultIdx = state.faultIdx
	snapshot.Schedulers = make(map[string]SchedulerStateSnapshot)
	for k, v := range state.schedulers {
		snapshot.Schedulers[k] = SchedulerStateSnapshot{
			Interval: v.interval,
			Next:     v.next.UTC(),
		}
	}
	snapshot.AlertInfos = make(map[int32]AlertInfosStateSnapShot)
	for k, v := range state.alertInfos {
		snapshot.AlertInfos[k] = AlertInfosStateSnapShot{
			Sn:    v.Sequence,
			Epoch: v.StartEpoch,
		}
	}
	snapshot.StorageFault = make(map[string]int32)
	for k, v := range state.storage {
		snapshot.StorageFault[k] = v
	}
	return snapshot
}

func (state *inMemState) Restore(snapshot *AgentStateSnapshot) {
	state.hbIdx = snapshot.HbIdx
	state.measIdx = snapshot.MeasIdx
	state.faultIdx = snapshot.FaultIdx
	state.schedulers = make(map[string]*schedulerState)
	for k, v := range snapshot.Schedulers {
		state.schedulers[k] = &schedulerState{
			interval: v.Interval,
			next:     v.Next.Local(),
		}
	}
	for k, v := range snapshot.AlertInfos {
		state.alertInfos[k] = &convert.AlertInfos{
			Sequence:   v.Sn,
			StartEpoch: v.Epoch,
		}
	}
	for k, v := range snapshot.StorageFault {
		state.storage[k] = v
	}
}
