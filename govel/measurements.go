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

import "time"

// VNICPerformance describes the performance and errors of an identified virtual network interface card
type VNICPerformance struct {
	ReceivedBroadcastPacketsAccumulated    *float64         `json:"receivedBroadcastPacketsAccumulated,omitempty"`    // Cumulative count of broadcast packets received as read at the end of the measurement; interval
	ReceivedBroadcastPacketsDelta          *float64         `json:"receivedBroadcastPacketsDelta,omitempty"`          // Count of broadcast packets received within the measurement interval
	ReceivedDiscardedPacketsAccumulated    *float64         `json:"receivedDiscardedPacketsAccumulated,omitempty"`    // Cumulative count of discarded packets received as read at the end of the measurement; interval
	ReceivedDiscardedPacketsDelta          *float64         `json:"receivedDiscardedPacketsDelta,omitempty"`          // Count of discarded packets received within the measurement interval
	ReceivedErrorPacketsAccumulated        *float64         `json:"receivedErrorPacketsAccumulated,omitempty"`        // Cumulative count of error packets received as read at the end of the measurement interval
	ReceivedErrorPacketsDelta              *float64         `json:"receivedErrorPacketsDelta,omitempty"`              // Count of error packets received within the measurement interval
	ReceivedMulticastPacketsAccumulated    *float64         `json:"receivedMulticastPacketsAccumulated,omitempty"`    // Cumulative count of multicast packets received as read at the end of the measurement; interval
	ReceivedMulticastPacketsDelta          *float64         `json:"receivedMulticastPacketsDelta,omitempty"`          // Count of multicast packets received within the measurement interval
	ReceivedOctetsAccumulated              *float64         `json:"receivedOctetsAccumulated,omitempty"`              // Cumulative count of octets received as read at the end of the measurement interval
	ReceivedOctetsDelta                    *float64         `json:"receivedOctetsDelta,omitempty"`                    // Count of octets received within the measurement interval
	ReceivedTotalPacketsAccumulated        *float64         `json:"receivedTotalPacketsAccumulated,omitempty"`        // Cumulative count of all packets received as read at the end of the measurement interval
	ReceivedTotalPacketsDelta              *float64         `json:"receivedTotalPacketsDelta,omitempty"`              // Count of all packets received within the measurement interval
	ReceivedUnicastPacketsAccumulated      *float64         `json:"receivedUnicastPacketsAccumulated,omitempty"`      // Cumulative count of unicast packets received as read at the end of the measurement; interval
	ReceivedUnicastPacketsDelta            *float64         `json:"receivedUnicastPacketsDelta,omitempty"`            // Count of unicast packets received within the measurement interval
	TransmittedBroadcastPacketsAccumulated *float64         `json:"transmittedBroadcastPacketsAccumulated,omitempty"` // Cumulative count of broadcast packets transmitted as read at the end of the measurement; interval
	TransmittedBroadcastPacketsDelta       *float64         `json:"transmittedBroadcastPacketsDelta,omitempty"`       // Count of broadcast packets transmitted within the measurement interval
	TransmittedDiscardedPacketsAccumulated *float64         `json:"transmittedDiscardedPacketsAccumulated,omitempty"` // Cumulative count of discarded packets transmitted as read at the end of the measurement; interval
	TransmittedDiscardedPacketsDelta       *float64         `json:"transmittedDiscardedPacketsDelta,omitempty"`       // Count of discarded packets transmitted within the measurement interval
	TransmittedErrorPacketsAccumulated     *float64         `json:"transmittedErrorPacketsAccumulated,omitempty"`     // Cumulative count of error packets transmitted as read at the end of the measurement; interval
	TransmittedErrorPacketsDelta           *float64         `json:"transmittedErrorPacketsDelta,omitempty"`           // Count of error packets transmitted within the measurement interval
	TransmittedMulticastPacketsAccumulated *float64         `json:"transmittedMulticastPacketsAccumulated,omitempty"` // Cumulative count of multicast packets transmitted as read at the end of the measurement; interval
	TransmittedMulticastPacketsDelta       *float64         `json:"transmittedMulticastPacketsDelta,omitempty"`       // Count of multicast packets transmitted within the measurement interval
	TransmittedOctetsAccumulated           *float64         `json:"transmittedOctetsAccumulated,omitempty"`           // Cumulative count of octets transmitted as read at the end of the measurement interval
	TransmittedOctetsDelta                 *float64         `json:"transmittedOctetsDelta,omitempty"`                 // Count of octets transmitted within the measurement interval
	TransmittedTotalPacketsAccumulated     *float64         `json:"transmittedTotalPacketsAccumulated,omitempty"`     // Cumulative count of all packets transmitted as read at the end of the measurement interval
	TransmittedTotalPacketsDelta           *float64         `json:"transmittedTotalPacketsDelta,omitempty"`           // Count of all packets transmitted within the measurement interval
	TransmittedUnicastPacketsAccumulated   *float64         `json:"transmittedUnicastPacketsAccumulated,omitempty"`   // Cumulative count of unicast packets transmitted as read at the end of the measurement; interval
	TransmittedUnicastPacketsDelta         *float64         `json:"transmittedUnicastPacketsDelta,omitempty"`         // Count of unicast packets transmitted within the measurement interval
	ValuesAreSuspect                       ValuesAreSuspect `json:"valuesAreSuspect"`                                 // Indicates whether vNicPerformance values are likely inaccurate due to counter overflow or; other condtions
	VNICIdentifier                         string           `json:"vNicIdentifier"`                                   // vNic identification
}

// ValuesAreSuspect Indicates whether vNicPerformance values are likely inaccurate due to counter overflow or
// other condtions
type ValuesAreSuspect string

// Values for ValuesAreSuspect
const (
	False ValuesAreSuspect = "false"
	True  ValuesAreSuspect = "true"
)

// MemoryUsage memory usage of an identified virtual machine
type MemoryUsage struct {
	MemoryBuffered   *float64 `json:"memoryBuffered,omitempty"`   // kibibytes of temporary storage for raw disk blocks
	MemoryCached     *float64 `json:"memoryCached,omitempty"`     // kibibytes of memory used for cache
	MemoryConfigured *float64 `json:"memoryConfigured,omitempty"` // kibibytes of memory configured in the virtual machine on which the VNFC reporting the; event is running
	MemoryFree       float64  `json:"memoryFree"`                 // kibibytes of physical RAM left unused by the system
	MemorySlabRecl   *float64 `json:"memorySlabRecl,omitempty"`   // the part of the slab that can be reclaimed such as caches measured in kibibytes
	MemorySlabUnrecl *float64 `json:"memorySlabUnrecl,omitempty"` // the part of the slab that cannot be reclaimed even when lacking memory measured in; kibibytes
	MemoryUsed       float64  `json:"memoryUsed"`                 // total memory minus the sum of free, buffered, cached and slab memory measured in kibibytes
	VMIdentifier     string   `json:"vmIdentifier"`               // virtual machine identifier associated with the memory metrics
}

// LatencyBucketMeasure number of counts falling within a defined latency bucket
type LatencyBucketMeasure struct {
	CountsInTheBucket      float64  `json:"countsInTheBucket"`
	HighEndOfLatencyBucket *float64 `json:"highEndOfLatencyBucket,omitempty"`
	LowEndOfLatencyBucket  *float64 `json:"lowEndOfLatencyBucket,omitempty"`
}

// FilesystemUsage disk usage of an identified virtual machine in gigabytes and/or gigabytes per second
type FilesystemUsage struct {
	BlockConfigured     float64 `json:"blockConfigured"`
	BlockIops           float64 `json:"blockIops"`
	BlockUsed           float64 `json:"blockUsed"`
	EphemeralConfigured float64 `json:"ephemeralConfigured"`
	EphemeralIops       float64 `json:"ephemeralIops"`
	EphemeralUsed       float64 `json:"ephemeralUsed"`
	FilesystemName      string  `json:"filesystemName"`
}

// FeaturesInUse number of times an identified feature was used over the measurementInterval
type FeaturesInUse struct {
	FeatureIdentifier  string `json:"featureIdentifier"`
	FeatureUtilization int64  `json:"featureUtilization"`
}

// DiskUsage usage of an identified disk
type DiskUsage struct {
	DiskIdentifier            string   `json:"diskIdentifier"`                      // disk identifier
	DiskIoTimeAvg             *float64 `json:"diskIoTimeAvg,omitempty"`             // milliseconds spent doing input/output operations over 1 sec; treat this metric as a; device load percentage where 1000ms  matches 100% load; provide the average over the; measurement interval
	DiskIoTimeLast            *float64 `json:"diskIoTimeLast,omitempty"`            // milliseconds spent doing input/output operations over 1 sec; treat this metric as a; device load percentage where 1000ms  matches 100% load; provide the last value; measurement within the measurement interval
	DiskIoTimeMax             *float64 `json:"diskIoTimeMax,omitempty"`             // milliseconds spent doing input/output operations over 1 sec; treat this metric as a; device load percentage where 1000ms  matches 100% load; provide the maximum value; measurement within the measurement interval
	DiskIoTimeMin             *float64 `json:"diskIoTimeMin,omitempty"`             // milliseconds spent doing input/output operations over 1 sec; treat this metric as a; device load percentage where 1000ms  matches 100% load; provide the minimum value; measurement within the measurement interval
	DiskMergedReadAvg         *float64 `json:"diskMergedReadAvg,omitempty"`         // number of logical read operations that were merged into physical read operations, e.g.,; two logical reads were served by one physical disk access; provide the average; measurement within the measurement interval
	DiskMergedReadLast        *float64 `json:"diskMergedReadLast,omitempty"`        // number of logical read operations that were merged into physical read operations, e.g.,; two logical reads were served by one physical disk access; provide the last value; measurement within the measurement interval
	DiskMergedReadMax         *float64 `json:"diskMergedReadMax,omitempty"`         // number of logical read operations that were merged into physical read operations, e.g.,; two logical reads were served by one physical disk access; provide the maximum value; measurement within the measurement interval
	DiskMergedReadMin         *float64 `json:"diskMergedReadMin,omitempty"`         // number of logical read operations that were merged into physical read operations, e.g.,; two logical reads were served by one physical disk access; provide the minimum value; measurement within the measurement interval
	DiskMergedWriteAvg        *float64 `json:"diskMergedWriteAvg,omitempty"`        // number of logical write operations that were merged into physical write operations, e.g.,; two logical writes were served by one physical disk access; provide the average; measurement within the measurement interval
	DiskMergedWriteLast       *float64 `json:"diskMergedWriteLast,omitempty"`       // number of logical write operations that were merged into physical write operations, e.g.,; two logical writes were served by one physical disk access; provide the last value; measurement within the measurement interval
	DiskMergedWriteMax        *float64 `json:"diskMergedWriteMax,omitempty"`        // number of logical write operations that were merged into physical write operations, e.g.,; two logical writes were served by one physical disk access; provide the maximum value; measurement within the measurement interval
	DiskMergedWriteMin        *float64 `json:"diskMergedWriteMin,omitempty"`        // number of logical write operations that were merged into physical write operations, e.g.,; two logical writes were served by one physical disk access; provide the minimum value; measurement within the measurement interval
	DiskOctetsReadAvg         *float64 `json:"diskOctetsReadAvg,omitempty"`         // number of octets per second read from a disk or partition; provide the average; measurement within the measurement interval
	DiskOctetsReadLast        *float64 `json:"diskOctetsReadLast,omitempty"`        // number of octets per second read from a disk or partition; provide the last measurement; within the measurement interval
	DiskOctetsReadMax         *float64 `json:"diskOctetsReadMax,omitempty"`         // number of octets per second read from a disk or partition; provide the maximum; measurement within the measurement interval
	DiskOctetsReadMin         *float64 `json:"diskOctetsReadMin,omitempty"`         // number of octets per second read from a disk or partition; provide the minimum; measurement within the measurement interval
	DiskOctetsWriteAvg        *float64 `json:"diskOctetsWriteAvg,omitempty"`        // number of octets per second written to a disk or partition; provide the average; measurement within the measurement interval
	DiskOctetsWriteLast       *float64 `json:"diskOctetsWriteLast,omitempty"`       // number of octets per second written to a disk or partition; provide the last measurement; within the measurement interval
	DiskOctetsWriteMax        *float64 `json:"diskOctetsWriteMax,omitempty"`        // number of octets per second written to a disk or partition; provide the maximum; measurement within the measurement interval
	DiskOctetsWriteMin        *float64 `json:"diskOctetsWriteMin,omitempty"`        // number of octets per second written to a disk or partition; provide the minimum; measurement within the measurement interval
	DiskOpsReadAvg            *float64 `json:"diskOpsReadAvg,omitempty"`            // number of read operations per second issued to the disk; provide the average measurement; within the measurement interval
	DiskOpsReadLast           *float64 `json:"diskOpsReadLast,omitempty"`           // number of read operations per second issued to the disk; provide the last measurement; within the measurement interval
	DiskOpsReadMax            *float64 `json:"diskOpsReadMax,omitempty"`            // number of read operations per second issued to the disk; provide the maximum measurement; within the measurement interval
	DiskOpsReadMin            *float64 `json:"diskOpsReadMin,omitempty"`            // number of read operations per second issued to the disk; provide the minimum measurement; within the measurement interval
	DiskOpsWriteAvg           *float64 `json:"diskOpsWriteAvg,omitempty"`           // number of write operations per second issued to the disk; provide the average measurement; within the measurement interval
	DiskOpsWriteLast          *float64 `json:"diskOpsWriteLast,omitempty"`          // number of write operations per second issued to the disk; provide the last measurement; within the measurement interval
	DiskOpsWriteMax           *float64 `json:"diskOpsWriteMax,omitempty"`           // number of write operations per second issued to the disk; provide the maximum measurement; within the measurement interval
	DiskOpsWriteMin           *float64 `json:"diskOpsWriteMin,omitempty"`           // number of write operations per second issued to the disk; provide the minimum measurement; within the measurement interval
	DiskPendingOperationsAvg  *float64 `json:"diskPendingOperationsAvg,omitempty"`  // queue size of pending I/O operations per second; provide the average measurement within; the measurement interval
	DiskPendingOperationsLast *float64 `json:"diskPendingOperationsLast,omitempty"` // queue size of pending I/O operations per second; provide the last measurement within the; measurement interval
	DiskPendingOperationsMax  *float64 `json:"diskPendingOperationsMax,omitempty"`  // queue size of pending I/O operations per second; provide the maximum measurement within; the measurement interval
	DiskPendingOperationsMin  *float64 `json:"diskPendingOperationsMin,omitempty"`  // queue size of pending I/O operations per second; provide the minimum measurement within; the measurement interval
	DiskTimeReadAvg           *float64 `json:"diskTimeReadAvg,omitempty"`           // milliseconds a read operation took to complete; provide the average measurement within; the measurement interval
	DiskTimeReadLast          *float64 `json:"diskTimeReadLast,omitempty"`          // milliseconds a read operation took to complete; provide the last measurement within the; measurement interval
	DiskTimeReadMax           *float64 `json:"diskTimeReadMax,omitempty"`           // milliseconds a read operation took to complete; provide the maximum measurement within; the measurement interval
	DiskTimeReadMin           *float64 `json:"diskTimeReadMin,omitempty"`           // milliseconds a read operation took to complete; provide the minimum measurement within; the measurement interval
	DiskTimeWriteAvg          *float64 `json:"diskTimeWriteAvg,omitempty"`          // milliseconds a write operation took to complete; provide the average measurement within; the measurement interval
	DiskTimeWriteLast         *float64 `json:"diskTimeWriteLast,omitempty"`         // milliseconds a write operation took to complete; provide the last measurement within the; measurement interval
	DiskTimeWriteMax          *float64 `json:"diskTimeWriteMax,omitempty"`          // milliseconds a write operation took to complete; provide the maximum measurement within; the measurement interval
	DiskTimeWriteMin          *float64 `json:"diskTimeWriteMin,omitempty"`          // milliseconds a write operation took to complete; provide the minimum measurement within; the measurement interval
}

// CPUUsage usage of an identified CPU
type CPUUsage struct {
	CPUIdentifier     string   `json:"cpuIdentifier"`               // cpu identifer
	CPUIdle           *float64 `json:"cpuIdle,omitempty"`           // percentage of CPU time spent in the idle task
	CPUUsageInterrupt *float64 `json:"cpuUsageInterrupt,omitempty"` // percentage of time spent servicing interrupts
	CPUUsageNice      *float64 `json:"cpuUsageNice,omitempty"`      // percentage of time spent running user space processes that have been niced
	CPUUsageSoftIRQ   *float64 `json:"cpuUsageSoftIrq,omitempty"`   // percentage of time spent handling soft irq interrupts
	CPUUsageSteal     *float64 `json:"cpuUsageSteal,omitempty"`     // percentage of time spent in involuntary wait which is neither user, system or idle time; and is effectively time that went missing
	CPUUsageSystem    *float64 `json:"cpuUsageSystem,omitempty"`    // percentage of time spent on system tasks running the kernel
	CPUUsageUser      *float64 `json:"cpuUsageUser,omitempty"`      // percentage of time spent running un-niced user space processes
	CPUWait           *float64 `json:"cpuWait,omitempty"`           // percentage of CPU time spent waiting for I/O operations to complete
	PercentUsage      float64  `json:"percentUsage"`                // aggregate cpu usage of the virtual machine on which the VNFC reporting the event is; running
}

// CodecsInUse number of times an identified codec was used over the measurementInterval
type CodecsInUse struct {
	CodecIdentifier string `json:"codecIdentifier"`
	NumberInUse     int64  `json:"numberInUse"`
}

// Key tuple which provides the name of a key along with its value and relative order
type Key struct {
	KeyName  string  `json:"keyName"`            // name of the key
	KeyOrder *int64  `json:"keyOrder,omitempty"` // relative sequence or order of the key with respect to other keys
	KeyValue *string `json:"keyValue,omitempty"` // value of the key
}

// JSONObjectInstance meta-information about an instance of a jsonObject along with the actual object instance
type JSONObjectInstance struct {
	ObjectInstance              map[string]interface{} `json:"objectInstance"`                        // an instance conforming to the jsonObject schema
	ObjectInstanceEpochMicrosec *float64               `json:"objectInstanceEpochMicrosec,omitempty"` // the unix time aka epoch time associated with this objectInstance--as microseconds elapsed; since 1 Jan 1970 not including leap seconds
	ObjectKeys                  []Key                  `json:"objectKeys"`                            // an ordered set of keys that identifies this particular instance of jsonObject
}

// JSONObject json object schema, name and other meta-information along with one or more object
// instances
type JSONObject struct {
	NfSubscribedObjectName *string              `json:"nfSubscribedObjectName,omitempty"` // name of the object associated with the nfSubscriptonId
	NfSubscriptionID       *string              `json:"nfSubscriptionId,omitempty"`       // identifies an openConfig telemetry subscription on a network function, which configures; the network function to send complex object data associated with the jsonObject
	ObjectInstances        []JSONObjectInstance `json:"objectInstances"`                  // one or more instances of the jsonObject
	ObjectName             string               `json:"objectName"`                       // name of the JSON Object
	ObjectSchema           *string              `json:"objectSchema,omitempty"`           // json schema for the object
	ObjectSchemaURL        *string              `json:"objectSchemaUrl,omitempty"`        // Url to the json schema for the object
}

// NamedArrayOfFields an array of name value pairs along with a name for the array
type NamedArrayOfFields struct {
	ArrayOfFields []Field `json:"arrayOfFields"` // array of name value pairs
	Name          string  `json:"name"`
}

// Field name value pair
type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type measurementsForVfScalingFields struct {
	AdditionalFields                []Field                `json:"additionalFields,omitempty"`        // additional name-value-pair fields
	AdditionalMeasurements          []NamedArrayOfFields   `json:"additionalMeasurements,omitempty"`  // array of named name-value-pair arrays
	AdditionalObjects               []JSONObject           `json:"additionalObjects,omitempty"`       // array of JSON objects described by name, schema and other meta-information
	CodecUsageArray                 []CodecsInUse          `json:"codecUsageArray,omitempty"`         // array of codecs in use
	ConcurrentSessions              *int64                 `json:"concurrentSessions,omitempty"`      // peak concurrent sessions for the VM or VNF over the measurementInterval
	ConfiguredEntities              *int64                 `json:"configuredEntities,omitempty"`      // over the measurementInterval, peak total number of: users, subscribers, devices,; adjacencies, etc., for the VM, or subscribers, devices, etc., for the VNF
	CPUUsageArray                   []CPUUsage             `json:"cpuUsageArray,omitempty"`           // usage of an array of CPUs
	DiskUsageArray                  []DiskUsage            `json:"diskUsageArray,omitempty"`          // usage of an array of disks
	FeatureUsageArray               []FeaturesInUse        `json:"featureUsageArray,omitempty"`       // array of features in use
	FilesystemUsageArray            []FilesystemUsage      `json:"filesystemUsageArray,omitempty"`    // filesystem usage of the VM on which the VNFC reporting the event is running
	LatencyDistribution             []LatencyBucketMeasure `json:"latencyDistribution,omitempty"`     // array of integers representing counts of requests whose latency in milliseconds falls; within per-VNF configured ranges
	MeanRequestLatency              *float64               `json:"meanRequestLatency,omitempty"`      // mean seconds required to respond to each request for the VM on which the VNFC reporting; the event is running
	MeasurementInterval             float64                `json:"measurementInterval"`               // interval over which measurements are being reported in seconds
	MeasurementsForVfScalingVersion float64                `json:"measurementsForVfScalingVersion"`   // version of the measurementsForVfScaling block
	MemoryUsageArray                []MemoryUsage          `json:"memoryUsageArray,omitempty"`        // memory usage of an array of VMs
	NumberOfMediaPortsInUse         *int64                 `json:"numberOfMediaPortsInUse,omitempty"` // number of media ports in use
	RequestRate                     *float64               `json:"requestRate,omitempty"`             // peak rate of service requests per second to the VNF over the measurementInterval
	VnfcScalingMetric               *int64                 `json:"vnfcScalingMetric,omitempty"`       // represents busy-ness of the VNF from 0 to 100 as reported by the VNFC
	VNICPerformanceArray            []VNICPerformance      `json:"vNicPerformanceArray,omitempty"`    // usage of an array of virtual network interface cards
}

// EventMeasurements is a metric event
type EventMeasurements struct {
	EventHeader                    `json:"commonEventHeader"`
	measurementsForVfScalingFields `json:"measurementsForVfScalingFields"`
}

// NewMeasurements creates a new measurements event
func NewMeasurements(name, id string, sourceName string, interval time.Duration, start, end time.Time) *EventMeasurements {
	meas := new(EventMeasurements)

	meas.MeasurementsForVfScalingVersion = 2.0
	meas.MeasurementInterval = interval.Seconds()

	meas.Domain = DomainMeasurementsForVfScaling
	meas.SourceName = sourceName
	meas.EventName = name
	meas.EventID = id
	meas.Version = 3.0
	meas.Priority = PriorityNormal

	meas.StartEpochMicrosec = start.UnixNano() / 1000
	meas.LastEpochMicrosec = end.UnixNano() / 1000
	meas.Sequence = 0

	return meas
}
