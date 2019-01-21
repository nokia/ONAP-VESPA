package evel

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
