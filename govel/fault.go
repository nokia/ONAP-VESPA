package govel

import (
	"time"
)

// Severity for faults
type Severity string

// Possible values for Severity
const (
	SeverityCritical Severity = "CRITICAL"
	SeverityMajor    Severity = "MAJOR"
	SeverityMinor    Severity = "MINOR"
	SeverityWarning  Severity = "WARNING"
	SeverityNormal   Severity = "NORMAL"
)

// VfStatus is the virtual function status
type VfStatus string

// Possible values for VfStatus
const (
	StatusActive           VfStatus = "Active"
	StatusIdle             VfStatus = "Idle"
	StatusPrepTerminate    VfStatus = "Preparing to terminate"
	StatusReadyTerminate   VfStatus = "Ready to terminate"
	StatusRequestTerminate VfStatus = "Requesting termination"
)

// SourceType of the fault originator
type SourceType string

// Possible values for SourceType
const (
	SourceOther                  SourceType = "other"
	SourceRouter                 SourceType = "router"
	SourceSwitch                 SourceType = "switch"
	SourceHost                   SourceType = "host"
	SourceCard                   SourceType = "card"
	SourcePort                   SourceType = "port"
	SourceSlotThreshold          SourceType = "slotThreshold"
	SourcePortThreshold          SourceType = "portThreshold"
	SourceVirtualMachine         SourceType = "virtualMachine"
	SourceVirtualNetworkFunction SourceType = "virtualNetworkFunction"
)

type faultFields struct {
	AlarmAdditionalInformation []EventField `json:"alarmAdditionalInformation,omitempty"`
	AlarmCondition             string       `json:"alarmCondition"`
	AlarmInterfaceA            string       `json:"alarmInterfaceA,omitempty"`
	EventCategory              string       `json:"eventCategory,omitempty"`
	EventSeverity              Severity     `json:"eventSeverity"`
	EventSourceType            SourceType   `json:"eventSourceType"`
	FaultFieldsVersion         float32      `json:"faultFieldsVersion"`
	SpecificProblem            string       `json:"specificProblem"`
	VfStatus                   VfStatus     `json:"vfStatus"`
}

//EventFault is a fault event
type EventFault struct {
	EventHeader `json:"commonEventHeader"`
	faultFields `json:"faultFields"`
}

// NewFault creates a new fault event
func NewFault(name, id, condition, specificProblem string, priority EventPriority, severity Severity, sourceType SourceType, status VfStatus, sourceName string) *EventFault {
	fault := new(EventFault)

	fault.AlarmCondition = condition
	fault.SpecificProblem = specificProblem
	fault.EventSeverity = severity
	fault.EventSourceType = sourceType
	fault.VfStatus = status
	fault.FaultFieldsVersion = 2.0

	fault.Domain = DomainFault
	fault.SourceName = sourceName
	fault.EventName = name
	fault.EventID = id
	fault.Version = 3.0
	fault.Priority = priority

	fault.StartEpochMicrosec = time.Now().UnixNano() / 1000
	fault.LastEpochMicrosec = fault.StartEpochMicrosec

	return fault
}
