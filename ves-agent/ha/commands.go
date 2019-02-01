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
	"fmt"
	"time"
)

const nullValue = "<null>"

// CmdType is the kind of command sent into commit logs
type CmdType int

// CmdType values
const (
	IncrementAlarmIdx CmdType = iota
	IncrementMeasIdx
	IncrementHeartbeatIdx
	UpdateScheduler
	IncrementFaultIdx
	UpdateFault
	DeleteFault
)

// StateCmd is a state change command sent through commit logs
type StateCmd struct {
	// Kind of command
	Type CmdType `json:"ty"`
	// Fields for command of kind UpdateScheduler
	UpdateScheduler *UpdateSchedulerFields `json:"sched,omitempty"`
	// Fields for command of kind UpdateFault
	UpdateFault *UpdateFaultFields `json:"updatefault,omitempty"`
	// Fields for command of kind DeleteFault
	DeleteFault *DeleteFaultFields `json:"deletefault,omitempty"`
}

func (cmd *StateCmd) String() string {
	switch cmd.Type {
	case IncrementAlarmIdx:
		return "IncrementAlarmIdx"
	case IncrementMeasIdx:
		return "IncrementMeasIdx"
	case IncrementHeartbeatIdx:
		return "IncrementHeartbeatIdx"
	case UpdateScheduler:
		return fmt.Sprintf("UpdateScheduler => %s", cmd.UpdateScheduler.String())
	case IncrementFaultIdx:
		return "IncrementFaultIdx"
	case UpdateFault:
		return fmt.Sprintf("UpdateFault => %s", cmd.UpdateFault.String())
	case DeleteFault:
		return fmt.Sprintf("DeleteFault => %s", cmd.DeleteFault.String())
	default:
		return fmt.Sprintf("Unknown command type: %d", cmd.Type)
	}
}

// UpdateSchedulerFields holds the fields for command of kind UpdateScheduler
type UpdateSchedulerFields struct {
	// Name of scheduler to update
	Name string `json:"name"`
	// New value of interval, if updated, or nil
	Interval *time.Duration `json:"intv,omitempty"`
	// New value of next run epoch time (in seconds), if updated, or nil
	Next *int64 `json:"nxt,omitempty"`
}

// UpdateFaultFields holds the fields for command of kind UpdateFault
type UpdateFaultFields struct {
	// FaultID of the fault to update
	FaultID *int32 `json:"faultId"`
	// FaultName of the fault to create
	FaultName string `json:"faultName,omitempty"`
	// New value of sequenceNumber
	SequenceNumber *int64 `json:"sn,omitempty"`
	// New value of startEpoch
	StartEpoch *int64 `json:"epoch,omitempty"`
}

// DeleteFaultFields holds the fields for command of kind DeleteFault
type DeleteFaultFields struct {
	// FaultName of the fault to create
	FaultName string `json:"faultName"`
}

func (fields *UpdateSchedulerFields) String() string {
	if fields == nil {
		return nullValue
	}
	var nxt *time.Time
	if fields.Next != nil {
		t := time.Unix(*fields.Next, 0)
		nxt = &t
	}
	return fmt.Sprintf("name: %s, interval: %s, next: %s", fields.Name, fields.Interval, nxt)
}

func (fields *UpdateFaultFields) String() string {
	if fields == nil {
		return nullValue
	}
	var sn int64
	if fields.SequenceNumber != nil {
		sn = *fields.SequenceNumber
	}
	var epoch int64
	if fields.StartEpoch != nil {
		epoch = *fields.StartEpoch
	}
	return fmt.Sprintf("faultId: %10d, faultName: %s, sn: %d, epoch: %d", *fields.FaultID, fields.FaultName, sn, epoch)
}

func (fields *DeleteFaultFields) String() string {
	if fields == nil {
		return nullValue
	}
	return fmt.Sprintf("faultName: %s", fields.FaultName)
}
