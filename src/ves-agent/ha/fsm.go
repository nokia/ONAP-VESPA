package ha

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/raft"
	log "github.com/sirupsen/logrus"
)

// FSM is the internal Finite State Machine handling
// internal state changes across raft cluster
type FSM struct {
	debug bool
	state SnapshotableAgentState
}

// NewFSM creates a new finite state machine wrapping a snapshotable state, and handling
// state changes across the raft cluster
// If `debug` is true, Commit logs are displayed in log output
func NewFSM(state SnapshotableAgentState, debug bool) *FSM {
	return &FSM{
		debug: debug,
		state: state,
	}
}

// Interval returns the scheduler exceution interval
func (fsm *FSM) Interval(sched string) time.Duration {
	return fsm.state.Interval(sched)
}

// NextRun returns the time at which next execution should occure
func (fsm *FSM) NextRun(sched string) time.Time {
	return fsm.state.NextRun(sched)
}

// GetFaultSn return the fault sequence number
func (fsm *FSM) GetFaultSn(fault int32) int64 {
	return fsm.state.GetFaultSn(fault)
}

// GetFaultStartEpoch return the startEpoch
func (fsm *FSM) GetFaultStartEpoch(fault int32) int64 {
	return fsm.state.GetFaultStartEpoch(fault)
}

// InitAlertInfos update the alertInfos map
func (fsm *FSM) InitAlertInfos(fault int32) error {
	return fsm.state.InitAlertInfos(fault)
}

// GetFaultInStorage checks if faultName already associated to an index
func (fsm *FSM) GetFaultInStorage(faultName string) int32 {
	return fsm.state.GetFaultInStorage(faultName)
}

// DeleteFaultInStorage delete Fault in storage
func (fsm *FSM) DeleteFaultInStorage(faultName string) error {
	return fsm.state.DeleteFaultInStorage(faultName)
}

// Apply applies a Raft log to this FSM
func (fsm *FSM) Apply(logEntry *raft.Log) interface{} {
	var cmd StateCmd
	if err := json.Unmarshal(logEntry.Data, &cmd); err != nil {
		log.Errorf("Cannot decode state command: %s. Reason: %s", string(logEntry.Data), err.Error())
		return err
	}
	if fsm.debug {
		log.Infof("Apply Log - [%s]", cmd.String())
	}
	res, err := fsm.processCmd(cmd)
	if err != nil {
		log.Errorf("Cannot appy command: %s", err.Error())
		return err
	}
	return res
}

// Snapshot creates and return a new snapshot of the surrent state
func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	log.Info("Snapshotting state")
	return fsm.state.Snapshot(), nil
}

// Restore deserialize and applies a snapshot to this FSM, discarding previous state
func (fsm *FSM) Restore(input io.ReadCloser) error {
	log.Info("Restoring snapshot")
	defer func() {
		if err := input.Close(); err != nil {
			log.Errorf(err.Error())
		}
	}()
	snapshot := AgentStateSnapshot{}
	if err := json.NewDecoder(input).Decode(&snapshot); err != nil {
		return err
	}
	fsm.state.Restore(&snapshot)
	return nil
}

func (fsm *FSM) processCmd(cmd StateCmd) (interface{}, error) {
	switch cmd.Type {
	case IncrementMeasIdx:
		return fsm.state.NextMeasurementIndex()
	case IncrementHeartbeatIdx:
		return fsm.state.NextHeartbeatIndex()
	case UpdateScheduler:
		return nil, fsm.handleSchedulerUpdate(cmd.UpdateScheduler)
	case IncrementFaultIdx:
		return fsm.state.NextFaultIndex()
	case UpdateFault:
		return nil, fsm.handleFaultUpdate(cmd.UpdateFault)
	case DeleteFault:
		return nil, fsm.state.DeleteFaultInStorage(cmd.DeleteFault.FaultName)
	default:
		return nil, fmt.Errorf("Unknown command type: %d", cmd.Type)
	}
}

func (fsm *FSM) handleSchedulerUpdate(fields *UpdateSchedulerFields) error {
	if fields == nil {
		return errors.New("UpdateScheduler field is absent")
	}
	if fields.Interval != nil {
		if err := fsm.state.UpdateInterval(fields.Name, *fields.Interval); err != nil {
			return err
		}
	}
	if fields.Next != nil {
		if err := fsm.state.UpdateNextRun(fields.Name, time.Unix(*fields.Next, 0)); err != nil {
			return err
		}
	}
	return nil
}

func (fsm *FSM) handleFaultUpdate(fields *UpdateFaultFields) error {
	//var debugMsg string
	if fields == nil {
		return errors.New("UpdateFaultFields field is absent")
	}
	if fields.FaultID != nil {
		//debugMsg = debugMsg + " faultId " + *fields.FaultID
		if fields.FaultName != "" {
			//debugMsg = debugMsg + "faultName " + fields.FaultName
			if err := fsm.state.StoreFaultInStorage(fields.FaultName, *fields.FaultID); err != nil {
				return err
			}
		}
		if fields.SequenceNumber != nil {
			//debugMsg = debugMsg + "sn " + *fields.SequenceNumber
			if err := fsm.state.IncrementFaultSn(*fields.FaultID); err != nil {
				return err
			}
		}
		if fields.StartEpoch != nil {
			//debugMsg = debugMsg + "epoch " + *fields.StartEpoch
			if err := fsm.state.SetFaultStartEpoch(*fields.FaultID, *fields.StartEpoch); err != nil {
				return err
			}
		}
	}
	return nil
}
