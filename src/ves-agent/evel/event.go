package evel

// EventDomain is the kind of event
type EventDomain string

// List of possible values for EventDomain
const (
	DomainFault                    EventDomain = "fault"
	DomainHeartbeat                EventDomain = "heartbeat"
	DomainMeasurementsForVfScaling EventDomain = "measurementsForVfScaling"
	DomainMobileFlow               EventDomain = "mobileFlow"
	DomainOther                    EventDomain = "other"
	DomainSipSignaling             EventDomain = "sipSignaling"
	DomainStateChange              EventDomain = "stateChange"
	DomainSyslog                   EventDomain = "syslog"
	DomainThresholdCrossingAlert   EventDomain = "thresholdCrossingAlert"
	DomainVoiceQuality             EventDomain = "voiceQuality"
)

// EventPriority is the event's level of priority
type EventPriority string

// Possible values for EventPriority
const (
	PriorityHigh   EventPriority = "High"
	PriorityMedium EventPriority = "Medium"
	PriorityNormal EventPriority = "Normal"
	PriorityLow    EventPriority = "Low"
)

// Event has commont methods for VES Events  structures
type Event interface {
	// Header returns a reference to the event's commonEventHeader
	Header() *EventHeader
}

// Batch is a list of events
type Batch []Event

// Split extracts cut the batch into 2 batches of equal length ( +/- 1)
func (batch Batch) Split() (Batch, Batch) {
	i := len(batch) / 2
	return batch[:i], batch[i:]
}

// Len returns the number of events in the batch
func (batch Batch) Len() int {
	return len(batch)
}

// UpdateReportingEntityName will update `reportingEntityName` field
// on events of batch for which this field has no value. Other events
// which already have the field set will be left untouched
func (batch Batch) UpdateReportingEntityName(name string) {
	for _, evt := range batch {
		if evt.Header().ReportingEntityName == "" {
			evt.Header().ReportingEntityName = name
		}
	}
}

// UpdateReportingEntityID will update `reportingEntityID` field
// on events of batch for which this field has no value. Other events
// which already have the field set will be left untouched
func (batch Batch) UpdateReportingEntityID(id string) {
	for _, evt := range batch {
		if evt.Header().ReportingEntityID == "" {
			evt.Header().ReportingEntityID = id
		}
	}
}

// EventHeader is the common part of all kind of events
type EventHeader struct {
	Domain               EventDomain   `json:"domain"`
	EventID              string        `json:"eventId"`
	EventName            string        `json:"eventName"`
	EventType            string        `json:"eventType,omitempty"`
	InternalHeaderFields interface{}   `json:"internalHeaderFields,omitempty"`
	LastEpochMicrosec    int64         `json:"lastEpochMicrosec"`
	NfNamingCode         string        `json:"nfNamingCode,omitempty"`
	NfcNamingCode        string        `json:"nfcNamingCode,omitempty"`
	Priority             EventPriority `json:"priority"`
	ReportingEntityID    string        `json:"reportingEntityId,omitempty"`
	ReportingEntityName  string        `json:"reportingEntityName"`
	Sequence             int64         `json:"sequence"`
	SourceID             string        `json:"sourceId,omitempty"`
	SourceName           string        `json:"sourceName"`
	StartEpochMicrosec   int64         `json:"startEpochMicrosec"`
	Version              float32       `json:"version"`
}

// Header returns a reference self
func (hdr *EventHeader) Header() *EventHeader {
	return hdr
}

// EventField is used for additional events fields
type EventField struct {
	// Name of the field
	Name string `json:"name"`
	// Value of the field
	Value string `json:"value"`
}
