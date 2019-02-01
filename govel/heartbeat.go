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
	"time"
)

type heartbeatFields struct {
	AdditionalFields       []EventField `json:"additionalFields,omitempty"`
	HeartbeatFieldsVersion float32      `json:"heartbeatFieldsVersion"`
	HeartbeatInterval      int          `json:"heartbeatInterval"`
}

// HeartbeatEvent with optional field block for fields specific to heartbeat events
type HeartbeatEvent struct {
	EventHeader     `json:"commonEventHeader"`
	heartbeatFields `json:"heartbeatFields,omitempty"`
}

// NewHeartbeat creates a new heartbeat event
func NewHeartbeat(id, name, sourceName string, interval int) *HeartbeatEvent {
	hb := new(HeartbeatEvent)
	hb.Domain = DomainHeartbeat
	hb.Priority = PriorityNormal
	hb.Version = 3.0
	hb.EventID = id
	hb.EventName = name
	hb.SourceName = sourceName

	hb.StartEpochMicrosec = time.Now().UnixNano() / 1000
	hb.LastEpochMicrosec = hb.StartEpochMicrosec

	//hb fields
	hb.HeartbeatFieldsVersion = 1.0
	hb.HeartbeatInterval = interval

	return hb
}
