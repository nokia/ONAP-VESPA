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

package schema

const v2841 = `{
	"$schema": "http://json-schema.org/draft-04/schema#",
	"title": "VES Event Listener",
	"type": "object",
	"properties": {
		"event": {
			"$ref": "#/definitions/event"
		},
		"eventList": {
			"$ref": "#/definitions/eventList"
		}
	},
	"definitions": {
		"schemaHeaderBlock": {
			"description": "schema date, version, author and associated API",
			"type": "object",
			"properties": {
				"associatedApi": {
					"description": "VES Event Listener",
					"type": "string"
				},
				"lastUpdatedBy": {
					"description": "re2947",
					"type": "string"
				},
				"schemaDate": {
					"description": "September 19, 2017",
					"type": "string"
				},
				"schemaVersion": {
					"description": "28.4.1",
					"type": "number"
				}
			}
		},
		"schemaLicenseAndCopyrightNotice": {
			"description": "Copyright (c) 2017, AT&T Intellectual Property.  All rights reserved",
			"type": "object",
			"properties": {
				"apacheLicense2.0": {
					"description": "Licensed under the Apache License, Version 2.0 (the 'License'); you may not use this file except in compliance with the License. You may obtain a copy of the License at:",
					"type": "string"
				},
				"licenseUrl": {
					"description": "http://www.apache.org/licenses/LICENSE-2.0",
					"type": "string"
				},
				"asIsClause": {
					"description": "Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an 'AS IS' BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.",
					"type": "string"
				},
				"permissionsAndLimitations": {
					"description": "See the License for the specific language governing permissions and limitations under the License.",
					"type": "string"
				}
			}
		},
		"codecsInUse": {
			"description": "number of times an identified codec was used over the measurementInterval",
			"type": "object",
			"properties": {
				"codecIdentifier": {
					"type": "string"
				},
				"numberInUse": {
					"type": "integer"
				}
			},
			"required": [
				"codecIdentifier",
				"numberInUse"
			]
		},
		"command": {
			"description": "command from an event collector toward an event source",
			"type": "object",
			"properties": {
				"commandType": {
					"type": "string",
					"enum": [
						"heartbeatIntervalChange",
						"measurementIntervalChange",
						"provideThrottlingState",
						"throttlingSpecification"
					]
				},
				"eventDomainThrottleSpecification": {
					"$ref": "#/definitions/eventDomainThrottleSpecification"
				},
				"heartbeatInterval": {
					"type": "integer"
				},
				"measurementInterval": {
					"type": "integer"
				}
			},
			"required": [
				"commandType"
			]
		},
		"commandList": {
			"description": "array of commands from an event collector toward an event source",
			"type": "array",
			"items": {
				"$ref": "#/definitions/command"
			},
			"minItems": 0
		},
		"commonEventHeader": {
			"description": "fields common to all events",
			"type": "object",
			"properties": {
				"domain": {
					"description": "the eventing domain associated with the event",
					"type": "string",
					"enum": [
						"fault",
						"heartbeat",
						"measurementsForVfScaling",
						"mobileFlow",
						"other",
						"sipSignaling",
						"stateChange",
						"syslog",
						"thresholdCrossingAlert",
						"voiceQuality"
					]
				},
				"eventId": {
					"description": "event key that is unique to the event source",
					"type": "string"
				},
				"eventName": {
					"description": "unique event name",
					"type": "string"
				},
				"eventType": {
					"description": "for example - applicationVnf, guestOS, hostOS, platform",
					"type": "string"
				},
				"internalHeaderFields": {
					"$ref": "#/definitions/internalHeaderFields"
				},
				"lastEpochMicrosec": {
					"description": "the latest unix time aka epoch time associated with the event from any component--as microseconds elapsed since 1 Jan 1970 not including leap seconds",
					"type": "number"
				},
				"nfcNamingCode": {
					"description": "3 character network function component type, aligned with vfc naming standards",
					"type": "string"
				},
				"nfNamingCode": {
					"description": "4 character network function type, aligned with vnf naming standards",
					"type": "string"
				},
				"priority": {
					"description": "processing priority",
					"type": "string",
					"enum": [
						"High",
						"Medium",
						"Normal",
						"Low"
					]
				},
				"reportingEntityId": {
					"description": "UUID identifying the entity reporting the event, for example an OAM VM; must be populated by the ATT enrichment process",
					"type": "string"
				},
				"reportingEntityName": {
					"description": "name of the entity reporting the event, for example, an EMS name; may be the same as sourceName",
					"type": "string"
				},
				"sequence": {
					"description": "ordering of events communicated by an event source instance or 0 if not needed",
					"type": "integer"
				},
				"sourceId": {
					"description": "UUID identifying the entity experiencing the event issue; must be populated by the ATT enrichment process",
					"type": "string"
				},
				"sourceName": {
					"description": "name of the entity experiencing the event issue",
					"type": "string"
				},
				"startEpochMicrosec": {
					"description": "the earliest unix time aka epoch time associated with the event from any component--as microseconds elapsed since 1 Jan 1970 not including leap seconds",
					"type": "number"
				},
				"version": {
					"description": "version of the event header",
					"type": "number"
				}
			},
			"required": [
				"domain",
				"eventId",
				"eventName",
				"lastEpochMicrosec",
				"priority",
				"reportingEntityName",
				"sequence",
				"sourceName",
				"startEpochMicrosec",
				"version"
			]
		},
		"counter": {
			"description": "performance counter",
			"type": "object",
			"properties": {
				"criticality": {
					"type": "string",
					"enum": [
						"CRIT",
						"MAJ"
					]
				},
				"name": {
					"type": "string"
				},
				"thresholdCrossed": {
					"type": "string"
				},
				"value": {
					"type": "string"
				}
			},
			"required": [
				"criticality",
				"name",
				"thresholdCrossed",
				"value"
			]
		},
		"cpuUsage": {
			"description": "usage of an identified CPU",
			"type": "object",
			"properties": {
				"cpuIdentifier": {
					"description": "cpu identifer",
					"type": "string"
				},
				"cpuIdle": {
					"description": "percentage of CPU time spent in the idle task",
					"type": "number"
				},
				"cpuUsageInterrupt": {
					"description": "percentage of time spent servicing interrupts",
					"type": "number"
				},
				"cpuUsageNice": {
					"description": "percentage of time spent running user space processes that have been niced",
					"type": "number"
				},
				"cpuUsageSoftIrq": {
					"description": "percentage of time spent handling soft irq interrupts",
					"type": "number"
				},
				"cpuUsageSteal": {
					"description": "percentage of time spent in involuntary wait which is neither user, system or idle time and is effectively time that went missing",
					"type": "number"
				},
				"cpuUsageSystem": {
					"description": "percentage of time spent on system tasks running the kernel",
					"type": "number"
				},
				"cpuUsageUser": {
					"description": "percentage of time spent running un-niced user space processes",
					"type": "number"
				},
				"cpuWait": {
					"description": "percentage of CPU time spent waiting for I/O operations to complete",
					"type": "number"
				},
				"percentUsage": {
					"description": "aggregate cpu usage of the virtual machine on which the VNFC reporting the event is running",
					"type": "number"
				}
			},
			"required": [
				"cpuIdentifier",
				"percentUsage"
			]
		},
		"diskUsage": {
			"description": "usage of an identified disk",
			"type": "object",
			"properties": {
				"diskIdentifier": {
					"description": "disk identifier",
					"type": "string"
				},
				"diskIoTimeAvg": {
					"description": "milliseconds spent doing input/output operations over 1 sec; treat this metric as a device load percentage where 1000ms  matches 100% load; provide the average over the measurement interval",
					"type": "number"
				},
				"diskIoTimeLast": {
					"description": "milliseconds spent doing input/output operations over 1 sec; treat this metric as a device load percentage where 1000ms  matches 100% load; provide the last value measurement within the measurement interval",
					"type": "number"
				},
				"diskIoTimeMax": {
					"description": "milliseconds spent doing input/output operations over 1 sec; treat this metric as a device load percentage where 1000ms  matches 100% load; provide the maximum value measurement within the measurement interval",
					"type": "number"
				},
				"diskIoTimeMin": {
					"description": "milliseconds spent doing input/output operations over 1 sec; treat this metric as a device load percentage where 1000ms  matches 100% load; provide the minimum value measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedReadAvg": {
					"description": "number of logical read operations that were merged into physical read operations, e.g., two logical reads were served by one physical disk access; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedReadLast": {
					"description": "number of logical read operations that were merged into physical read operations, e.g., two logical reads were served by one physical disk access; provide the last value measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedReadMax": {
					"description": "number of logical read operations that were merged into physical read operations, e.g., two logical reads were served by one physical disk access; provide the maximum value measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedReadMin": {
					"description": "number of logical read operations that were merged into physical read operations, e.g., two logical reads were served by one physical disk access; provide the minimum value measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedWriteAvg": {
					"description": "number of logical write operations that were merged into physical write operations, e.g., two logical writes were served by one physical disk access; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedWriteLast": {
					"description": "number of logical write operations that were merged into physical write operations, e.g., two logical writes were served by one physical disk access; provide the last value measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedWriteMax": {
					"description": "number of logical write operations that were merged into physical write operations, e.g., two logical writes were served by one physical disk access; provide the maximum value measurement within the measurement interval",
					"type": "number"
				},
				"diskMergedWriteMin": {
					"description": "number of logical write operations that were merged into physical write operations, e.g., two logical writes were served by one physical disk access; provide the minimum value measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsReadAvg": {
					"description": "number of octets per second read from a disk or partition; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsReadLast": {
					"description": "number of octets per second read from a disk or partition; provide the last measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsReadMax": {
					"description": "number of octets per second read from a disk or partition; provide the maximum measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsReadMin": {
					"description": "number of octets per second read from a disk or partition; provide the minimum measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsWriteAvg": {
					"description": "number of octets per second written to a disk or partition; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsWriteLast": {
					"description": "number of octets per second written to a disk or partition; provide the last measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsWriteMax": {
					"description": "number of octets per second written to a disk or partition; provide the maximum measurement within the measurement interval",
					"type": "number"
				},
				"diskOctetsWriteMin": {
					"description": "number of octets per second written to a disk or partition; provide the minimum measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsReadAvg": {
					"description": "number of read operations per second issued to the disk; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsReadLast": {
					"description": "number of read operations per second issued to the disk; provide the last measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsReadMax": {
					"description": "number of read operations per second issued to the disk; provide the maximum measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsReadMin": {
					"description": "number of read operations per second issued to the disk; provide the minimum measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsWriteAvg": {
					"description": "number of write operations per second issued to the disk; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsWriteLast": {
					"description": "number of write operations per second issued to the disk; provide the last measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsWriteMax": {
					"description": "number of write operations per second issued to the disk; provide the maximum measurement within the measurement interval",
					"type": "number"
				},
				"diskOpsWriteMin": {
					"description": "number of write operations per second issued to the disk; provide the minimum measurement within the measurement interval",
					"type": "number"
				},
				"diskPendingOperationsAvg": {
					"description": "queue size of pending I/O operations per second; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskPendingOperationsLast": {
					"description": "queue size of pending I/O operations per second; provide the last measurement within the measurement interval",
					"type": "number"
				},
				"diskPendingOperationsMax": {
					"description": "queue size of pending I/O operations per second; provide the maximum measurement within the measurement interval",
					"type": "number"
				},
				"diskPendingOperationsMin": {
					"description": "queue size of pending I/O operations per second; provide the minimum measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeReadAvg": {
					"description": "milliseconds a read operation took to complete; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeReadLast": {
					"description": "milliseconds a read operation took to complete; provide the last measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeReadMax": {
					"description": "milliseconds a read operation took to complete; provide the maximum measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeReadMin": {
					"description": "milliseconds a read operation took to complete; provide the minimum measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeWriteAvg": {
					"description": "milliseconds a write operation took to complete; provide the average measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeWriteLast": {
					"description": "milliseconds a write operation took to complete; provide the last measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeWriteMax": {
					"description": "milliseconds a write operation took to complete; provide the maximum measurement within the measurement interval",
					"type": "number"
				},
				"diskTimeWriteMin": {
					"description": "milliseconds a write operation took to complete; provide the minimum measurement within the measurement interval",
					"type": "number"
				}
			},
			"required": [
				"diskIdentifier"
			]
		},
		"endOfCallVqmSummaries": {
			"description": "provides end of call voice quality metrics",
			"type": "object",
			"properties": {
				"adjacencyName": {
					"description": " adjacency name",
					"type": "string"
				},
				"endpointDescription": {
					"description": "Either Caller or Callee",
					"type": "string",
					"enum": [
						"Caller",
						"Callee"
					]
				},
				"endpointJitter": {
					"description": "",
					"type": "number"
				},
				"endpointRtpOctetsDiscarded": {
					"description": "",
					"type": "number"
				},
				"endpointRtpOctetsReceived": {
					"description": "",
					"type": "number"
				},
				"endpointRtpOctetsSent": {
					"description": "",
					"type": "number"
				},
				"endpointRtpPacketsDiscarded": {
					"description": "",
					"type": "number"
				},
				"endpointRtpPacketsReceived": {
					"description": "",
					"type": "number"
				},
				"endpointRtpPacketsSent": {
					"description": "",
					"type": "number"
				},
				"localJitter": {
					"description": "",
					"type": "number"
				},
				"localRtpOctetsDiscarded": {
					"description": "",
					"type": "number"
				},
				"localRtpOctetsReceived": {
					"description": "",
					"type": "number"
				},
				"localRtpOctetsSent": {
					"description": "",
					"type": "number"
				},
				"localRtpPacketsDiscarded": {
					"description": "",
					"type": "number"
				},
				"localRtpPacketsReceived": {
					"description": "",
					"type": "number"
				},
				"localRtpPacketsSent": {
					"description": "",
					"type": "number"
				},
				"mosCqe": {
					"description": "1-5 1dp",
					"type": "number"
				},
				"packetsLost": {
					"description": "",
					"type": "number"
				},
				"packetLossPercent": {
					"description": "Calculated percentage packet loss based on Endpoint RTP packets lost (as reported in RTCP) and Local RTP packets sent. Direction is based on Endpoint description (Caller, Callee). Decimal (2 dp)",
					"type": "number"
				},
				"rFactor": {
					"description": "0-100",
					"type": "number"
				},
				"roundTripDelay": {
					"description": "millisecs",
					"type": "number"
				}
			},
			"required": [
				"adjacencyName",
				"endpointDescription"
			]
		},
		"event": {
			"description": "the root level of the common event format",
			"type": "object",
			"properties": {
				"commonEventHeader": {
					"$ref": "#/definitions/commonEventHeader"
				},
				"faultFields": {
					"$ref": "#/definitions/faultFields"
				},
				"heartbeatFields": {
					"$ref": "#/definitions/heartbeatFields"
				},
				"measurementsForVfScalingFields": {
					"$ref": "#/definitions/measurementsForVfScalingFields"
				},
				"mobileFlowFields": {
					"$ref": "#/definitions/mobileFlowFields"
				},
				"otherFields": {
					"$ref": "#/definitions/otherFields"
				},
				"sipSignalingFields": {
					"$ref": "#/definitions/sipSignalingFields"
				},
				"stateChangeFields": {
					"$ref": "#/definitions/stateChangeFields"
				},
				"syslogFields": {
					"$ref": "#/definitions/syslogFields"
				},
				"thresholdCrossingAlertFields": {
					"$ref": "#/definitions/thresholdCrossingAlertFields"
				},
				"voiceQualityFields": {
					"$ref": "#/definitions/voiceQualityFields"
				}
			},
			"required": [
				"commonEventHeader"
			]
		},
		"eventDomainThrottleSpecification": {
			"description": "specification of what information to suppress within an event domain",
			"type": "object",
			"properties": {
				"eventDomain": {
					"description": "Event domain enum from the commonEventHeader domain field",
					"type": "string"
				},
				"suppressedFieldNames": {
					"description": "List of optional field names in the event block that should not be sent to the Event Listener",
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"suppressedNvPairsList": {
					"description": "Optional list of specific NvPairsNames to suppress within a given Name-Value Field",
					"type": "array",
					"items": {
						"$ref": "#/definitions/suppressedNvPairs"
					}
				}
			},
			"required": [
				"eventDomain"
			]
		},
		"eventDomainThrottleSpecificationList": {
			"description": "array of eventDomainThrottleSpecifications",
			"type": "array",
			"items": {
				"$ref": "#/definitions/eventDomainThrottleSpecification"
			},
			"minItems": 0
		},
		"eventList": {
			"description": "array of events",
			"type": "array",
			"items": {
				"$ref": "#/definitions/event"
			}
		},
		"eventThrottlingState": {
			"description": "reports the throttling in force at the event source",
			"type": "object",
			"properties": {
				"eventThrottlingMode": {
					"description": "Mode the event manager is in",
					"type": "string",
					"enum": [
						"normal",
						"throttled"
					]
				},
				"eventDomainThrottleSpecificationList": {
					"$ref": "#/definitions/eventDomainThrottleSpecificationList"
				}
			},
			"required": [
				"eventThrottlingMode"
			]
		},
		"faultFields": {
			"description": "fields specific to fault events",
			"type": "object",
			"properties": {
				"alarmAdditionalInformation": {
					"description": "additional alarm information",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"alarmCondition": {
					"description": "alarm condition reported by the device",
					"type": "string"
				},
				"alarmInterfaceA": {
					"description": "card, port, channel or interface name of the device generating the alarm",
					"type": "string"
				},
				"eventCategory": {
					"description": "Event category, for example: license, link, routing, security, signaling",
					"type": "string"
				},
				"eventSeverity": {
					"description": "event severity",
					"type": "string",
					"enum": [
						"CRITICAL",
						"MAJOR",
						"MINOR",
						"WARNING",
						"NORMAL"
					]
				},
				"eventSourceType": {
					"description": "type of event source; examples: card, host, other, port, portThreshold, router, slotThreshold, switch, virtualMachine, virtualNetworkFunction",
					"type": "string"
				},
				"faultFieldsVersion": {
					"description": "version of the faultFields block",
					"type": "number"
				},
				"specificProblem": {
					"description": "short description of the alarm or problem",
					"type": "string"
				},
				"vfStatus": {
					"description": "virtual function status enumeration",
					"type": "string",
					"enum": [
						"Active",
						"Idle",
						"Preparing to terminate",
						"Ready to terminate",
						"Requesting termination"
					]
				}
			},
			"required": [
				"alarmCondition",
				"eventSeverity",
				"eventSourceType",
				"faultFieldsVersion",
				"specificProblem",
				"vfStatus"
			]
		},
		"featuresInUse": {
			"description": "number of times an identified feature was used over the measurementInterval",
			"type": "object",
			"properties": {
				"featureIdentifier": {
					"type": "string"
				},
				"featureUtilization": {
					"type": "integer"
				}
			},
			"required": [
				"featureIdentifier",
				"featureUtilization"
			]
		},
		"field": {
			"description": "name value pair",
			"type": "object",
			"properties": {
				"name": {
					"type": "string"
				},
				"value": {
					"type": "string"
				}
			},
			"required": [
				"name",
				"value"
			]
		},
		"filesystemUsage": {
			"description": "disk usage of an identified virtual machine in gigabytes and/or gigabytes per second",
			"type": "object",
			"properties": {
				"blockConfigured": {
					"type": "number"
				},
				"blockIops": {
					"type": "number"
				},
				"blockUsed": {
					"type": "number"
				},
				"ephemeralConfigured": {
					"type": "number"
				},
				"ephemeralIops": {
					"type": "number"
				},
				"ephemeralUsed": {
					"type": "number"
				},
				"filesystemName": {
					"type": "string"
				}
			},
			"required": [
				"blockConfigured",
				"blockIops",
				"blockUsed",
				"ephemeralConfigured",
				"ephemeralIops",
				"ephemeralUsed",
				"filesystemName"
			]
		},
		"gtpPerFlowMetrics": {
			"description": "Mobility GTP Protocol per flow metrics",
			"type": "object",
			"properties": {
				"avgBitErrorRate": {
					"description": "average bit error rate",
					"type": "number"
				},
				"avgPacketDelayVariation": {
					"description": "Average packet delay variation or jitter in milliseconds for received packets: Average difference between the packet timestamp and time received for all pairs of consecutive packets",
					"type": "number"
				},
				"avgPacketLatency": {
					"description": "average delivery latency",
					"type": "number"
				},
				"avgReceiveThroughput": {
					"description": "average receive throughput",
					"type": "number"
				},
				"avgTransmitThroughput": {
					"description": "average transmit throughput",
					"type": "number"
				},
				"durConnectionFailedStatus": {
					"description": "duration of failed state in milliseconds, computed as the cumulative time between a failed echo request and the next following successful error request, over this reporting interval",
					"type": "number"
				},
				"durTunnelFailedStatus": {
					"description": "Duration of errored state, computed as the cumulative time between a tunnel error indicator and the next following non-errored indicator, over this reporting interval",
					"type": "number"
				},
				"flowActivatedBy": {
					"description": "Endpoint activating the flow",
					"type": "string"
				},
				"flowActivationEpoch": {
					"description": "Time the connection is activated in the flow (connection) being reported on, or transmission time of the first packet if activation time is not available",
					"type": "number"
				},
				"flowActivationMicrosec": {
					"description": "Integer microseconds for the start of the flow connection",
					"type": "number"
				},
				"flowActivationTime": {
					"description": "time the connection is activated in the flow being reported on, or transmission time of the first packet if activation time is not available; with RFC 2822 compliant format: Sat, 13 Mar 2010 11:29:05 -0800",
					"type": "string"
				},
				"flowDeactivatedBy": {
					"description": "Endpoint deactivating the flow",
					"type": "string"
				},
				"flowDeactivationEpoch": {
					"description": "Time for the start of the flow connection, in integer UTC epoch time aka UNIX time",
					"type": "number"
				},
				"flowDeactivationMicrosec": {
					"description": "Integer microseconds for the start of the flow connection",
					"type": "number"
				},
				"flowDeactivationTime": {
					"description": "Transmission time of the first packet in the flow connection being reported on; with RFC 2822 compliant format: Sat, 13 Mar 2010 11:29:05 -0800",
					"type": "string"
				},
				"flowStatus": {
					"description": "connection status at reporting time as a working / inactive / failed indicator value",
					"type": "string"
				},
				"gtpConnectionStatus": {
					"description": "Current connection state at reporting time",
					"type": "string"
				},
				"gtpTunnelStatus": {
					"description": "Current tunnel state  at reporting time",
					"type": "string"
				},
				"ipTosCountList": {
					"description": "array of key: value pairs where the keys are drawn from the IP Type-of-Service identifiers which range from '0' to '255', and the values are the count of packets that had those ToS identifiers in the flow",
					"type": "array",
					"items": {
						"type": "array",
						"items": [
							{
								"type": "string"
							},
							{
								"type": "number"
							}
						]
					}
				},
				"ipTosList": {
					"description": "Array of unique IP Type-of-Service values observed in the flow where values range from '0' to '255'",
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"largePacketRtt": {
					"description": "large packet round trip time",
					"type": "number"
				},
				"largePacketThreshold": {
					"description": "large packet threshold being applied",
					"type": "number"
				},
				"maxPacketDelayVariation": {
					"description": "Maximum packet delay variation or jitter in milliseconds for received packets: Maximum of the difference between the packet timestamp and time received for all pairs of consecutive packets",
					"type": "number"
				},
				"maxReceiveBitRate": {
					"description": "maximum receive bit rate",
					"type": "number"
				},
				"maxTransmitBitRate": {
					"description": "maximum transmit bit rate",
					"type": "number"
				},
				"mobileQciCosCountList": {
					"description": "array of key: value pairs where the keys are drawn from LTE QCI or UMTS class of service strings, and the values are the count of packets that had those strings in the flow",
					"type": "array",
					"items": {
						"type": "array",
						"items": [
							{
								"type": "string"
							},
							{
								"type": "number"
							}
						]
					}
				},
				"mobileQciCosList": {
					"description": "Array of unique LTE QCI or UMTS class-of-service values observed in the flow",
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"numActivationFailures": {
					"description": "Number of failed activation requests, as observed by the reporting node",
					"type": "number"
				},
				"numBitErrors": {
					"description": "number of errored bits",
					"type": "number"
				},
				"numBytesReceived": {
					"description": "number of bytes received, including retransmissions",
					"type": "number"
				},
				"numBytesTransmitted": {
					"description": "number of bytes transmitted, including retransmissions",
					"type": "number"
				},
				"numDroppedPackets": {
					"description": "number of received packets dropped due to errors per virtual interface",
					"type": "number"
				},
				"numGtpEchoFailures": {
					"description": "Number of Echo request path failures where failed paths are defined in 3GPP TS 29.281 sec 7.2.1 and 3GPP TS 29.060 sec. 11.2",
					"type": "number"
				},
				"numGtpTunnelErrors": {
					"description": "Number of tunnel error indications where errors are defined in 3GPP TS 29.281 sec 7.3.1 and 3GPP TS 29.060 sec. 11.1",
					"type": "number"
				},
				"numHttpErrors": {
					"description": "Http error count",
					"type": "number"
				},
				"numL7BytesReceived": {
					"description": "number of tunneled layer 7 bytes received, including retransmissions",
					"type": "number"
				},
				"numL7BytesTransmitted": {
					"description": "number of tunneled layer 7 bytes transmitted, excluding retransmissions",
					"type": "number"
				},
				"numLostPackets": {
					"description": "number of lost packets",
					"type": "number"
				},
				"numOutOfOrderPackets": {
					"description": "number of out-of-order packets",
					"type": "number"
				},
				"numPacketErrors": {
					"description": "number of errored packets",
					"type": "number"
				},
				"numPacketsReceivedExclRetrans": {
					"description": "number of packets received, excluding retransmission",
					"type": "number"
				},
				"numPacketsReceivedInclRetrans": {
					"description": "number of packets received, including retransmission",
					"type": "number"
				},
				"numPacketsTransmittedInclRetrans": {
					"description": "number of packets transmitted, including retransmissions",
					"type": "number"
				},
				"numRetries": {
					"description": "number of packet retries",
					"type": "number"
				},
				"numTimeouts": {
					"description": "number of packet timeouts",
					"type": "number"
				},
				"numTunneledL7BytesReceived": {
					"description": "number of tunneled layer 7 bytes received, excluding retransmissions",
					"type": "number"
				},
				"roundTripTime": {
					"description": "round trip time",
					"type": "number"
				},
				"tcpFlagCountList": {
					"description": "array of key: value pairs where the keys are drawn from TCP Flags and the values are the count of packets that had that TCP Flag in the flow",
					"type": "array",
					"items": {
						"type": "array",
						"items": [
							{
								"type": "string"
							},
							{
								"type": "number"
							}
						]
					}
				},
				"tcpFlagList": {
					"description": "Array of unique TCP Flags observed in the flow",
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"timeToFirstByte": {
					"description": "Time in milliseconds between the connection activation and first byte received",
					"type": "number"
				}
			},
			"required": [
				"avgBitErrorRate",
				"avgPacketDelayVariation",
				"avgPacketLatency",
				"avgReceiveThroughput",
				"avgTransmitThroughput",
				"flowActivationEpoch",
				"flowActivationMicrosec",
				"flowDeactivationEpoch",
				"flowDeactivationMicrosec",
				"flowDeactivationTime",
				"flowStatus",
				"maxPacketDelayVariation",
				"numActivationFailures",
				"numBitErrors",
				"numBytesReceived",
				"numBytesTransmitted",
				"numDroppedPackets",
				"numL7BytesReceived",
				"numL7BytesTransmitted",
				"numLostPackets",
				"numOutOfOrderPackets",
				"numPacketErrors",
				"numPacketsReceivedExclRetrans",
				"numPacketsReceivedInclRetrans",
				"numPacketsTransmittedInclRetrans",
				"numRetries",
				"numTimeouts",
				"numTunneledL7BytesReceived",
				"roundTripTime",
				"timeToFirstByte"
			]
		},
		"heartbeatFields": {
			"description": "optional field block for fields specific to heartbeat events",
			"type": "object",
			"properties": {
				"additionalFields": {
					"description": "additional heartbeat fields if needed",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"heartbeatFieldsVersion": {
					"description": "version of the heartbeatFields block",
					"type": "number"
				},
				"heartbeatInterval": {
					"description": "current heartbeat interval in seconds",
					"type": "integer"
				}
			},
			"required": [
				"heartbeatFieldsVersion",
				"heartbeatInterval"
			]
		},
		"internalHeaderFields": {
			"description": "enrichment fields for internal VES Event Listener service use only, not supplied by event sources",
			"type": "object"
		},
		"jsonObject": {
			"description": "json object schema, name and other meta-information along with one or more object instances",
			"type": "object",
			"properties": {
				"objectInstances": {
					"description": "one or more instances of the jsonObject",
					"type": "array",
					"items": {
						"$ref": "#/definitions/jsonObjectInstance"
					}
				},
				"objectName": {
					"description": "name of the JSON Object",
					"type": "string"
				},
				"objectSchema": {
					"description": "json schema for the object",
					"type": "string"
				},
				"objectSchemaUrl": {
					"description": "Url to the json schema for the object",
					"type": "string"
				},
				"nfSubscribedObjectName": {
					"description": "name of the object associated with the nfSubscriptonId",
					"type": "string"
				},
				"nfSubscriptionId": {
					"description": "identifies an openConfig telemetry subscription on a network function, which configures the network function to send complex object data associated with the jsonObject",
					"type": "string"
				}
			},
			"required": [
				"objectInstances",
				"objectName"
			]
		},
		"jsonObjectInstance": {
			"description": "meta-information about an instance of a jsonObject along with the actual object instance",
			"type": "object",
			"properties": {
				"objectInstance": {
					"description": "an instance conforming to the jsonObject schema",
					"type": "object"
				},
				"objectInstanceEpochMicrosec": {
					"description": "the unix time aka epoch time associated with this objectInstance--as microseconds elapsed since 1 Jan 1970 not including leap seconds",
					"type": "number"
				},
				"objectKeys": {
					"description": "an ordered set of keys that identifies this particular instance of jsonObject",
					"type": "array",
					"items": {
						"$ref": "#/definitions/key"
					}
				}
			},
			"required": [
				"objectInstance"
			]
		},
		"key": {
			"description": "tuple which provides the name of a key along with its value and relative order",
			"type": "object",
			"properties": {
				"keyName": {
					"description": "name of the key",
					"type": "string"
				},
				"keyOrder": {
					"description": "relative sequence or order of the key with respect to other keys",
					"type": "integer"
				},
				"keyValue": {
					"description": "value of the key",
					"type": "string"
				}
			},
			"required": [
				"keyName"
			]
		},
		"latencyBucketMeasure": {
			"description": "number of counts falling within a defined latency bucket",
			"type": "object",
			"properties": {
				"countsInTheBucket": {
					"type": "number"
				},
				"highEndOfLatencyBucket": {
					"type": "number"
				},
				"lowEndOfLatencyBucket": {
					"type": "number"
				}
			},
			"required": [
				"countsInTheBucket"
			]
		},
		"measurementsForVfScalingFields": {
			"description": "measurementsForVfScaling fields",
			"type": "object",
			"properties": {
				"additionalFields": {
					"description": "additional name-value-pair fields",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"additionalMeasurements": {
					"description": "array of named name-value-pair arrays",
					"type": "array",
					"items": {
						"$ref": "#/definitions/namedArrayOfFields"
					}
				},
				"additionalObjects": {
					"description": "array of JSON objects described by name, schema and other meta-information",
					"type": "array",
					"items": {
						"$ref": "#/definitions/jsonObject"
					}
				},
				"codecUsageArray": {
					"description": "array of codecs in use",
					"type": "array",
					"items": {
						"$ref": "#/definitions/codecsInUse"
					}
				},
				"concurrentSessions": {
					"description": "peak concurrent sessions for the VM or VNF over the measurementInterval",
					"type": "integer"
				},
				"configuredEntities": {
					"description": "over the measurementInterval, peak total number of: users, subscribers, devices, adjacencies, etc., for the VM, or subscribers, devices, etc., for the VNF",
					"type": "integer"
				},
				"cpuUsageArray": {
					"description": "usage of an array of CPUs",
					"type": "array",
					"items": {
						"$ref": "#/definitions/cpuUsage"
					}
				},
				"diskUsageArray": {
					"description": "usage of an array of disks",
					"type": "array",
					"items": {
						"$ref": "#/definitions/diskUsage"
					}
				},
				"featureUsageArray": {
					"description": "array of features in use",
					"type": "array",
					"items": {
						"$ref": "#/definitions/featuresInUse"
					}
				},
				"filesystemUsageArray": {
					"description": "filesystem usage of the VM on which the VNFC reporting the event is running",
					"type": "array",
					"items": {
						"$ref": "#/definitions/filesystemUsage"
					}
				},
				"latencyDistribution": {
					"description": "array of integers representing counts of requests whose latency in milliseconds falls within per-VNF configured ranges",
					"type": "array",
					"items": {
						"$ref": "#/definitions/latencyBucketMeasure"
					}
				},
				"meanRequestLatency": {
					"description": "mean seconds required to respond to each request for the VM on which the VNFC reporting the event is running",
					"type": "number"
				},
				"measurementInterval": {
					"description": "interval over which measurements are being reported in seconds",
					"type": "number"
				},
				"measurementsForVfScalingVersion": {
					"description": "version of the measurementsForVfScaling block",
					"type": "number"
				},
				"memoryUsageArray": {
					"description": "memory usage of an array of VMs",
					"type": "array",
					"items": {
						"$ref": "#/definitions/memoryUsage"
					}
				},
				"numberOfMediaPortsInUse": {
					"description": "number of media ports in use",
					"type": "integer"
				},
				"requestRate": {
					"description": "peak rate of service requests per second to the VNF over the measurementInterval",
					"type": "number"
				},
				"vnfcScalingMetric": {
					"description": "represents busy-ness of the VNF from 0 to 100 as reported by the VNFC",
					"type": "integer"
				},
				"vNicPerformanceArray": {
					"description": "usage of an array of virtual network interface cards",
					"type": "array",
					"items": {
						"$ref": "#/definitions/vNicPerformance"
					}
				}
			},
			"required": [
				"measurementInterval",
				"measurementsForVfScalingVersion"
			]
		},
		"memoryUsage": {
			"description": "memory usage of an identified virtual machine",
			"type": "object",
			"properties": {
				"memoryBuffered": {
					"description": "kibibytes of temporary storage for raw disk blocks",
					"type": "number"
				},
				"memoryCached": {
					"description": "kibibytes of memory used for cache",
					"type": "number"
				},
				"memoryConfigured": {
					"description": "kibibytes of memory configured in the virtual machine on which the VNFC reporting the event is running",
					"type": "number"
				},
				"memoryFree": {
					"description": "kibibytes of physical RAM left unused by the system",
					"type": "number"
				},
				"memorySlabRecl": {
					"description": "the part of the slab that can be reclaimed such as caches measured in kibibytes",
					"type": "number"
				},
				"memorySlabUnrecl": {
					"description": "the part of the slab that cannot be reclaimed even when lacking memory measured in kibibytes",
					"type": "number"
				},
				"memoryUsed": {
					"description": "total memory minus the sum of free, buffered, cached and slab memory measured in kibibytes",
					"type": "number"
				},
				"vmIdentifier": {
					"description": "virtual machine identifier associated with the memory metrics",
					"type": "string"
				}
			},
			"required": [
				"memoryFree",
				"memoryUsed",
				"vmIdentifier"
			]
		},
		"mobileFlowFields": {
			"description": "mobileFlow fields",
			"type": "object",
			"properties": {
				"additionalFields": {
					"description": "additional mobileFlow fields if needed",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"applicationType": {
					"description": "Application type inferred",
					"type": "string"
				},
				"appProtocolType": {
					"description": "application protocol",
					"type": "string"
				},
				"appProtocolVersion": {
					"description": "application protocol version",
					"type": "string"
				},
				"cid": {
					"description": "cell id",
					"type": "string"
				},
				"connectionType": {
					"description": "Abbreviation referencing a 3GPP reference point e.g., S1-U, S11, etc",
					"type": "string"
				},
				"ecgi": {
					"description": "Evolved Cell Global Id",
					"type": "string"
				},
				"flowDirection": {
					"description": "Flow direction, indicating if the reporting node is the source of the flow or destination for the flow",
					"type": "string"
				},
				"gtpPerFlowMetrics": {
					"$ref": "#/definitions/gtpPerFlowMetrics"
				},
				"gtpProtocolType": {
					"description": "GTP protocol",
					"type": "string"
				},
				"gtpVersion": {
					"description": "GTP protocol version",
					"type": "string"
				},
				"httpHeader": {
					"description": "HTTP request header, if the flow connects to a node referenced by HTTP",
					"type": "string"
				},
				"imei": {
					"description": "IMEI for the subscriber UE used in this flow, if the flow connects to a mobile device",
					"type": "string"
				},
				"imsi": {
					"description": "IMSI for the subscriber UE used in this flow, if the flow connects to a mobile device",
					"type": "string"
				},
				"ipProtocolType": {
					"description": "IP protocol type e.g., TCP, UDP, RTP...",
					"type": "string"
				},
				"ipVersion": {
					"description": "IP protocol version e.g., IPv4, IPv6",
					"type": "string"
				},
				"lac": {
					"description": "location area code",
					"type": "string"
				},
				"mcc": {
					"description": "mobile country code",
					"type": "string"
				},
				"mnc": {
					"description": "mobile network code",
					"type": "string"
				},
				"mobileFlowFieldsVersion": {
					"description": "version of the mobileFlowFields block",
					"type": "number"
				},
				"msisdn": {
					"description": "MSISDN for the subscriber UE used in this flow, as an integer, if the flow connects to a mobile device",
					"type": "string"
				},
				"otherEndpointIpAddress": {
					"description": "IP address for the other endpoint, as used for the flow being reported on",
					"type": "string"
				},
				"otherEndpointPort": {
					"description": "IP Port for the reporting entity, as used for the flow being reported on",
					"type": "integer"
				},
				"otherFunctionalRole": {
					"description": "Functional role of the other endpoint for the flow being reported on e.g., MME, S-GW, P-GW, PCRF...",
					"type": "string"
				},
				"rac": {
					"description": "routing area code",
					"type": "string"
				},
				"radioAccessTechnology": {
					"description": "Radio Access Technology e.g., 2G, 3G, LTE",
					"type": "string"
				},
				"reportingEndpointIpAddr": {
					"description": "IP address for the reporting entity, as used for the flow being reported on",
					"type": "string"
				},
				"reportingEndpointPort": {
					"description": "IP port for the reporting entity, as used for the flow being reported on",
					"type": "integer"
				},
				"sac": {
					"description": "service area code",
					"type": "string"
				},
				"samplingAlgorithm": {
					"description": "Integer identifier for the sampling algorithm or rule being applied in calculating the flow metrics if metrics are calculated based on a sample of packets, or 0 if no sampling is applied",
					"type": "integer"
				},
				"tac": {
					"description": "transport area code",
					"type": "string"
				},
				"tunnelId": {
					"description": "tunnel identifier",
					"type": "string"
				},
				"vlanId": {
					"description": "VLAN identifier used by this flow",
					"type": "string"
				}
			},
			"required": [
				"flowDirection",
				"gtpPerFlowMetrics",
				"ipProtocolType",
				"ipVersion",
				"mobileFlowFieldsVersion",
				"otherEndpointIpAddress",
				"otherEndpointPort",
				"reportingEndpointIpAddr",
				"reportingEndpointPort"
			]
		},
		"namedArrayOfFields": {
			"description": "an array of name value pairs along with a name for the array",
			"type": "object",
			"properties": {
				"name": {
					"type": "string"
				},
				"arrayOfFields": {
					"description": "array of name value pairs",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				}
			},
			"required": [
				"name",
				"arrayOfFields"
			]
		},
		"otherFields": {
			"description": "fields for events belonging to the 'other' domain of the commonEventHeader domain enumeration",
			"type": "object",
			"properties": {
				"hashOfNameValuePairArrays": {
					"description": "array of named name-value-pair arrays",
					"type": "array",
					"items": {
						"$ref": "#/definitions/namedArrayOfFields"
					}
				},
				"jsonObjects": {
					"description": "array of JSON objects described by name, schema and other meta-information",
					"type": "array",
					"items": {
						"$ref": "#/definitions/jsonObject"
					}
				},
				"nameValuePairs": {
					"description": "array of name-value pairs",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"otherFieldsVersion": {
					"description": "version of the otherFields block",
					"type": "number"
				}
			},
			"required": [
				"otherFieldsVersion"
			]
		},
		"requestError": {
			"description": "standard request error data structure",
			"type": "object",
			"properties": {
				"messageId": {
					"description": "Unique message identifier of the format ABCnnnn where ABC is either SVC for Service Exceptions or POL for Policy Exception",
					"type": "string"
				},
				"text": {
					"description": "Message text, with replacement variables marked with %n, where n is an index into the list of <variables> elements, starting at 1",
					"type": "string"
				},
				"url": {
					"description": "Hyperlink to a detailed error resource e.g., an HTML page for browser user agents",
					"type": "string"
				},
				"variables": {
					"description": "List of zero or more strings that represent the contents of the variables used by the message text",
					"type": "string"
				}
			},
			"required": [
				"messageId",
				"text"
			]
		},
		"sipSignalingFields": {
			"description": "sip signaling fields",
			"type": "object",
			"properties": {
				"additionalInformation": {
					"description": "additional sip signaling fields if needed",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"compressedSip": {
					"description": "the full SIP request/response including headers and bodies",
					"type": "string"
				},
				"correlator": {
					"description": "this is the same for all events on this call",
					"type": "string"
				},
				"localIpAddress": {
					"description": "IP address on VNF",
					"type": "string"
				},
				"localPort": {
					"description": "port on VNF",
					"type": "string"
				},
				"remoteIpAddress": {
					"description": "IP address of peer endpoint",
					"type": "string"
				},
				"remotePort": {
					"description": "port of peer endpoint",
					"type": "string"
				},
				"sipSignalingFieldsVersion": {
					"description": "version of the sipSignalingFields block",
					"type": "number"
				},
				"summarySip": {
					"description": "the SIP Method or Response (‘INVITE’, ‘200 OK’, ‘BYE’, etc)",
					"type": "string"
				},
				"vendorVnfNameFields": {
					"$ref": "#/definitions/vendorVnfNameFields"
				}
			},
			"required": [
				"correlator",
				"localIpAddress",
				"localPort",
				"remoteIpAddress",
				"remotePort",
				"sipSignalingFieldsVersion",
				"vendorVnfNameFields"
			]
		},
		"stateChangeFields": {
			"description": "stateChange fields",
			"type": "object",
			"properties": {
				"additionalFields": {
					"description": "additional stateChange fields if needed",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"newState": {
					"description": "new state of the entity",
					"type": "string",
					"enum": [
						"inService",
						"maintenance",
						"outOfService"
					]
				},
				"oldState": {
					"description": "previous state of the entity",
					"type": "string",
					"enum": [
						"inService",
						"maintenance",
						"outOfService"
					]
				},
				"stateChangeFieldsVersion": {
					"description": "version of the stateChangeFields block",
					"type": "number"
				},
				"stateInterface": {
					"description": "card or port name of the entity that changed state",
					"type": "string"
				}
			},
			"required": [
				"newState",
				"oldState",
				"stateChangeFieldsVersion",
				"stateInterface"
			]
		},
		"suppressedNvPairs": {
			"description": "List of specific NvPairsNames to suppress within a given Name-Value Field for event Throttling",
			"type": "object",
			"properties": {
				"nvPairFieldName": {
					"description": "Name of the field within which are the nvpair names to suppress",
					"type": "string"
				},
				"suppressedNvPairNames": {
					"description": "Array of nvpair names to suppress within the nvpairFieldName",
					"type": "array",
					"items": {
						"type": "string"
					}
				}
			},
			"required": [
				"nvPairFieldName",
				"suppressedNvPairNames"
			]
		},
		"syslogFields": {
			"description": "sysLog fields",
			"type": "object",
			"properties": {
				"additionalFields": {
					"description": "additional syslog fields if needed provided as name=value delimited by a pipe ‘|’ symbol, for example: 'name1=value1|name2=value2|…'",
					"type": "string"
				},
				"eventSourceHost": {
					"description": "hostname of the device",
					"type": "string"
				},
				"eventSourceType": {
					"description": "type of event source; examples: other, router, switch, host, card, port, slotThreshold, portThreshold, virtualMachine, virtualNetworkFunction",
					"type": "string"
				},
				"syslogFacility": {
					"description": "numeric code from 0 to 23 for facility--see table in documentation",
					"type": "integer"
				},
				"syslogFieldsVersion": {
					"description": "version of the syslogFields block",
					"type": "number"
				},
				"syslogMsg": {
					"description": "syslog message",
					"type": "string"
				},
				"syslogPri": {
					"description": "0-192 combined severity and facility",
					"type": "integer"
				},
				"syslogProc": {
					"description": "identifies the application that originated the message",
					"type": "string"
				},
				"syslogProcId": {
					"description": "a change in the value of this field indicates a discontinuity in syslog reporting",
					"type": "number"
				},
				"syslogSData": {
					"description": "syslog structured data consisting of a structured data Id followed by a set of key value pairs",
					"type": "string"
				},
				"syslogSdId": {
					"description": "0-32 char in format name@number for example ourSDID@32473",
					"type": "string"
				},
				"syslogSev": {
					"description": "numerical Code for  severity derived from syslogPri as remaider of syslogPri / 8",
					"type": "string",
					"enum": [
						"Alert",
						"Critical",
						"Debug",
						"Emergency",
						"Error",
						"Info",
						"Notice",
						"Warning"
					]
				},
				"syslogTag": {
					"description": "msgId indicating the type of message such as TCPOUT or TCPIN; NILVALUE should be used when no other value can be provided",
					"type": "string"
				},
				"syslogVer": {
					"description": "IANA assigned version of the syslog protocol specification - typically 1",
					"type": "number"
				}
			},
			"required": [
				"eventSourceType",
				"syslogFieldsVersion",
				"syslogMsg",
				"syslogTag"
			]
		},
		"thresholdCrossingAlertFields": {
			"description": "fields specific to threshold crossing alert events",
			"type": "object",
			"properties": {
				"additionalFields": {
					"description": "additional threshold crossing alert fields if needed",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"additionalParameters": {
					"description": "performance counters",
					"type": "array",
					"items": {
						"$ref": "#/definitions/counter"
					}
				},
				"alertAction": {
					"description": "Event action",
					"type": "string",
					"enum": [
						"CLEAR",
						"CONT",
						"SET"
					]
				},
				"alertDescription": {
					"description": "Unique short alert description such as IF-SHUB-ERRDROP",
					"type": "string"
				},
				"alertType": {
					"description": "Event type",
					"type": "string",
					"enum": [
						"CARD-ANOMALY",
						"ELEMENT-ANOMALY",
						"INTERFACE-ANOMALY",
						"SERVICE-ANOMALY"
					]
				},
				"alertValue": {
					"description": "Calculated API value (if applicable)",
					"type": "string"
				},
				"associatedAlertIdList": {
					"description": "List of eventIds associated with the event being reported",
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"collectionTimestamp": {
					"description": "Time when the performance collector picked up the data; with RFC 2822 compliant format: Sat, 13 Mar 2010 11:29:05 -0800",
					"type": "string"
				},
				"dataCollector": {
					"description": "Specific performance collector instance used",
					"type": "string"
				},
				"elementType": {
					"description": "type of network element - internal ATT field",
					"type": "string"
				},
				"eventSeverity": {
					"description": "event severity or priority",
					"type": "string",
					"enum": [
						"CRITICAL",
						"MAJOR",
						"MINOR",
						"WARNING",
						"NORMAL"
					]
				},
				"eventStartTimestamp": {
					"description": "Time closest to when the measurement was made; with RFC 2822 compliant format: Sat, 13 Mar 2010 11:29:05 -0800",
					"type": "string"
				},
				"interfaceName": {
					"description": "Physical or logical port or card (if applicable)",
					"type": "string"
				},
				"networkService": {
					"description": "network name - internal ATT field",
					"type": "string"
				},
				"possibleRootCause": {
					"description": "Reserved for future use",
					"type": "string"
				},
				"thresholdCrossingFieldsVersion": {
					"description": "version of the thresholdCrossingAlertFields block",
					"type": "number"
				}
			},
			"required": [
				"additionalParameters",
				"alertAction",
				"alertDescription",
				"alertType",
				"collectionTimestamp",
				"eventSeverity",
				"eventStartTimestamp",
				"thresholdCrossingFieldsVersion"
			]
		},
		"vendorVnfNameFields": {
			"description": "provides vendor, vnf and vfModule identifying information",
			"type": "object",
			"properties": {
				"vendorName": {
					"description": "VNF vendor name",
					"type": "string"
				},
				"vfModuleName": {
					"description": "ASDC vfModuleName for the vfModule generating the event",
					"type": "string"
				},
				"vnfName": {
					"description": "ASDC modelName for the VNF generating the event",
					"type": "string"
				}
			},
			"required": [
				"vendorName"
			]
		},
		"vNicPerformance": {
			"description": "describes the performance and errors of an identified virtual network interface card",
			"type": "object",
			"properties": {
				"receivedBroadcastPacketsAccumulated": {
					"description": "Cumulative count of broadcast packets received as read at the end of the measurement interval",
					"type": "number"
				},
				"receivedBroadcastPacketsDelta": {
					"description": "Count of broadcast packets received within the measurement interval",
					"type": "number"
				},
				"receivedDiscardedPacketsAccumulated": {
					"description": "Cumulative count of discarded packets received as read at the end of the measurement interval",
					"type": "number"
				},
				"receivedDiscardedPacketsDelta": {
					"description": "Count of discarded packets received within the measurement interval",
					"type": "number"
				},
				"receivedErrorPacketsAccumulated": {
					"description": "Cumulative count of error packets received as read at the end of the measurement interval",
					"type": "number"
				},
				"receivedErrorPacketsDelta": {
					"description": "Count of error packets received within the measurement interval",
					"type": "number"
				},
				"receivedMulticastPacketsAccumulated": {
					"description": "Cumulative count of multicast packets received as read at the end of the measurement interval",
					"type": "number"
				},
				"receivedMulticastPacketsDelta": {
					"description": "Count of multicast packets received within the measurement interval",
					"type": "number"
				},
				"receivedOctetsAccumulated": {
					"description": "Cumulative count of octets received as read at the end of the measurement interval",
					"type": "number"
				},
				"receivedOctetsDelta": {
					"description": "Count of octets received within the measurement interval",
					"type": "number"
				},
				"receivedTotalPacketsAccumulated": {
					"description": "Cumulative count of all packets received as read at the end of the measurement interval",
					"type": "number"
				},
				"receivedTotalPacketsDelta": {
					"description": "Count of all packets received within the measurement interval",
					"type": "number"
				},
				"receivedUnicastPacketsAccumulated": {
					"description": "Cumulative count of unicast packets received as read at the end of the measurement interval",
					"type": "number"
				},
				"receivedUnicastPacketsDelta": {
					"description": "Count of unicast packets received within the measurement interval",
					"type": "number"
				},
				"transmittedBroadcastPacketsAccumulated": {
					"description": "Cumulative count of broadcast packets transmitted as read at the end of the measurement interval",
					"type": "number"
				},
				"transmittedBroadcastPacketsDelta": {
					"description": "Count of broadcast packets transmitted within the measurement interval",
					"type": "number"
				},
				"transmittedDiscardedPacketsAccumulated": {
					"description": "Cumulative count of discarded packets transmitted as read at the end of the measurement interval",
					"type": "number"
				},
				"transmittedDiscardedPacketsDelta": {
					"description": "Count of discarded packets transmitted within the measurement interval",
					"type": "number"
				},
				"transmittedErrorPacketsAccumulated": {
					"description": "Cumulative count of error packets transmitted as read at the end of the measurement interval",
					"type": "number"
				},
				"transmittedErrorPacketsDelta": {
					"description": "Count of error packets transmitted within the measurement interval",
					"type": "number"
				},
				"transmittedMulticastPacketsAccumulated": {
					"description": "Cumulative count of multicast packets transmitted as read at the end of the measurement interval",
					"type": "number"
				},
				"transmittedMulticastPacketsDelta": {
					"description": "Count of multicast packets transmitted within the measurement interval",
					"type": "number"
				},
				"transmittedOctetsAccumulated": {
					"description": "Cumulative count of octets transmitted as read at the end of the measurement interval",
					"type": "number"
				},
				"transmittedOctetsDelta": {
					"description": "Count of octets transmitted within the measurement interval",
					"type": "number"
				},
				"transmittedTotalPacketsAccumulated": {
					"description": "Cumulative count of all packets transmitted as read at the end of the measurement interval",
					"type": "number"
				},
				"transmittedTotalPacketsDelta": {
					"description": "Count of all packets transmitted within the measurement interval",
					"type": "number"
				},
				"transmittedUnicastPacketsAccumulated": {
					"description": "Cumulative count of unicast packets transmitted as read at the end of the measurement interval",
					"type": "number"
				},
				"transmittedUnicastPacketsDelta": {
					"description": "Count of unicast packets transmitted within the measurement interval",
					"type": "number"
				},
				"valuesAreSuspect": {
					"description": "Indicates whether vNicPerformance values are likely inaccurate due to counter overflow or other condtions",
					"type": "string",
					"enum": [
						"true",
						"false"
					]
				},
				"vNicIdentifier": {
					"description": "vNic identification",
					"type": "string"
				}
			},
			"required": [
				"valuesAreSuspect",
				"vNicIdentifier"
			]
		},
		"voiceQualityFields": {
			"description": "provides statistics related to customer facing voice products",
			"type": "object",
			"properties": {
				"additionalInformation": {
					"description": "additional voice quality fields if needed",
					"type": "array",
					"items": {
						"$ref": "#/definitions/field"
					}
				},
				"calleeSideCodec": {
					"description": "callee codec for the call",
					"type": "string"
				},
				"callerSideCodec": {
					"description": "caller codec for the call",
					"type": "string"
				},
				"correlator": {
					"description": "this is the same for all events on this call",
					"type": "string"
				},
				"endOfCallVqmSummaries": {
					"$ref": "#/definitions/endOfCallVqmSummaries"
				},
				"phoneNumber": {
					"description": "phone number associated with the correlator",
					"type": "string"
				},
				"midCallRtcp": {
					"description": "Base64 encoding of the binary RTCP data excluding Eth/IP/UDP headers",
					"type": "string"
				},
				"vendorVnfNameFields": {
					"$ref": "#/definitions/vendorVnfNameFields"
				},
				"voiceQualityFieldsVersion": {
					"description": "version of the voiceQualityFields block",
					"type": "number"
				}
			},
			"required": [
				"calleeSideCodec",
				"callerSideCodec",
				"correlator",
				"midCallRtcp",
				"vendorVnfNameFields",
				"voiceQualityFieldsVersion"
			]
		}
	}
}`

var _v2841 *JSONSchema

func init() {
	sch, err := NewSchemaFromBytes([]byte(v2841))
	if err != nil {
		panic(err)
	}
	_v2841 = sch
}

// V2841 loads and returns VES Schema v28.4.1 (VES v5.4.1)
func V2841() *JSONSchema {
	return _v2841
}
