package govel

// CommandType specifies the kind of command sent by server
type CommandType string

// Possible values for CommandType
const (
	CommandHeartbeatIntervalChange   CommandType = "heartbeatIntervalChange"
	CommandMeasurementIntervalChange CommandType = "measurementIntervalChange"
	CommandProvideThrottlingState    CommandType = "provideThrottlingState"
	CommandThrottllingSpecification  CommandType = "throttllingSpecification"
)

// SuppressedNvPairs datatype is a list of specific NvPairsNames to suppress within a given Name-Value Field (for event throttling);
type SuppressedNvPairs struct {
	NvPairFieldName       string   `json:"nvPairFieldName"`       // Name of the field within which are the nvpair names to suppress
	SuppressedNvPairNames []string `json:"suppressedNvPairNames"` // Array of nvpair names to suppress (within the nvpairFieldName)
}

// EventDomainThrottleSpecification datatype specifies what fields to suppress within an event domain
type EventDomainThrottleSpecification struct {
	EventDomain           EventDomain         `json:"eventDomain"`                     // Event domain enum from the commonEventHeader domain field
	SuppressedFieldNames  []string            `json:"suppressedFieldNames,omitempty"`  // List of optional field names in the event block that should not be sent to the Event Listener
	SuppressedNvPairsList []SuppressedNvPairs `json:"suppressedNvPairsList,omitempty"` // Optional list of specific NvPairsNames to suppress within a given Name-Value Field
}

// Command describe a command sent by server in replies
type Command struct {
	CommandType                      CommandType                       `json:"commandType"`
	EventDomainThrottleSpecification *EventDomainThrottleSpecification `json:"eventDomainThrottleSpecification,omitempty"` // If commandType is ‘throttlingSpecification’, the fields to suppress within an event domain
	HeartbeatInterval                int                               `json:"heartbeatInterval,omitempty"`                // If commandType is ‘heartbeatIntervalChange’, the heartbeatInterval duration to use in seconds
	MeasurementInterval              int                               `json:"measurementInterval,omitempty"`              // If commandType is ‘measurementIntervalChange’, the measurementInterval duration to use in seconds
}
