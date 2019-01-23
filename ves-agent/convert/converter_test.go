package convert

import (
	"encoding/json"
	"testing"
	"github.com/nokia/onap-vespa/govel"

	"github.com/stretchr/testify/suite"

	"github.com/prometheus/alertmanager/template"
)

var alertData1 = []byte(`
	{
			"status": "firing",
			"labels": {
				"id": "201",
				"system_name": "TEST1",
				"alertname": "NodeFailure",
				"severity": "critical",
				"VNFC": "dpa2bhsxp5001vm001oam001",
				"calm": "false",
				"event_type": "x2"
			},
			"annotations": {
				"service": "NodeSupervision",
				"summary": "Node pilot-0 down",
				"description": "VM node is seen disconnected from Cluster",
				"clearAlertName": "NodeEndOfFailure",
				"clearDescription": "VM node is seen  again connected from Cluster"
			}
	}`)

var alertData2 = []byte(`
	{
			"status": "firing",
			"labels": {
				"id": "201",
				"system_name": "TEST2",
				"alertname": "NodeFailure",
				"severity": "critical",
				"VNFC": "dpa2bhsxp5001vm001oam002",
				"calm": "false",
				"event_type": "x2"
			},
			"annotations": {
				"service": "NodeSupervision",
				"summary": "Node pilot-1 down",
				"description": "Node pilot-1 is seen disconnected from cluster",
				"clearAlertName": "NodeEndOfFailure"
			}
	}`)

var alertData3 = []byte(`
	{
			"status": "firing",
			"labels": {
				"id": "100",
				"system_name": "TEST3",
				"alertname": "VudrBegCommunicationAuditFailure",
				"severity": "minor",
				"calm": "false",
				"event_type": "x2"
			},
			"annotations": {
				"service": "COMMAudit_VUDR_1*BEG_1",
				"summary": "Node pilot-0 down",
				"description": "On demand or periodic audit between SDME and PGW End Point SPML has failed",
				"clearAlertName": "VudrBegCommunicationAuditEndOfFailure",
				"clearDescription": "On demand or periodic audit between SDME and PGW End Point now ok",
				"aaiMapping": "VudrName_BegName"
			}
	}`)

var alertData4 = []byte(`
	{
			"status": "firing",
			"labels": {
				"id": "302",
				"instance": "mjves-ope-1:9100",
				"monitor": "sdmexpert",
				"system_name": "TEST4",
				"alertname": "FileSystemOccupancyCrossedLowThreshold",
				"severity": "Minor",
				"VNFC": "mjves-ope-1",
				"calm": "false",
				"event_type": "x2",
				"job": "3gpp",
				"probable_cause": "351"
			},
			"annotations": {
				"service": "FileSystemSupervision_/var/log/audit=85.1234",
				"summary": "Filesystem /var/log/audit usage on instance mjves-ope-1",
				"description": "File system almost full",
				"clearAlertName": "FileSystemOccupancyReturnedUnderLowThreshold",
				"clearDescription": "File system occupancy returned under low threshold",
				"aaiMapping": "FileSystemName_ObservedThreshold"
			}
	}`)

var alertData5 = []byte(`
	{
			"status": "firing",
			"labels": {
				"id": "306",
				"instance": "mjves-ope-0:9100",
				"monitor": "sdmexpert",
				"system_name": "TEST5",
				"alertname": "MemoryOccupancyCrossedLowThreshold",
				"severity": "Minor",
				"VNFC": "mjves-ope-1",
				"calm": "false",
				"event_type": "x2",
				"job": "3gpp",
				"probable_cause": "351"
			},
			"annotations": {
				"service": "MemorySupervision=82.4321",
				"summary": "High memory usage on Instance mjves-ope-1",
				"description": "Memory high occupancy",
				"clearAlertName": "MemoryOccupancyReturnedUnderLowThreshold",
				"clearDescription": "Memory occupancy returned under low threshold",
				"aaiMapping": "ObservedThreshold"
			}
	}`)

type ConvertTestSuite struct {
	suite.Suite
	alert1      template.Alert
	alert2      template.Alert
	alert3      template.Alert
	alert4      template.Alert
	alert5      template.Alert
	confEvent   govel.EventConfiguration
	namingCodes map[string]string
}

func TestConvert(t *testing.T) {
	suite.Run(t, new(ConvertTestSuite))
}

func (suite *ConvertTestSuite) SetupSuite() {
	var err error
	suite.confEvent = govel.EventConfiguration{
		MaxSize:      200,
		NfNamingCode: "hspx",
	}

	suite.namingCodes = make(map[string]string)
	suite.namingCodes["dpa2bhsxp5001vm001oam001"] = "oam"
	suite.namingCodes["mjves-ope-1"] = "etl"

	err = json.Unmarshal(alertData1, &suite.alert1)
	if err != nil {
		suite.Fail("Error in unmarshall function for test alert1")
	}
	err = json.Unmarshal(alertData2, &suite.alert2)
	if err != nil {
		suite.Fail("Error in unmarshall function for test alert2")
	}
	err = json.Unmarshal(alertData3, &suite.alert3)
	if err != nil {
		suite.Fail("Error in unmarshall function for test alert3")
	}
	err = json.Unmarshal(alertData4, &suite.alert4)
	if err != nil {
		suite.Fail("Error in unmarshall function for test alert4")
	}
	err = json.Unmarshal(alertData5, &suite.alert5)
	if err != nil {
		suite.Fail("Error in unmarshall function for test alert5")
	}
}

func (suite *ConvertTestSuite) TestConvertAlertRaiseOK() {
	var eventFault *govel.EventFault
	var status StatusResult
	var fm = NewFaultManager(&suite.confEvent)
	var sequence int64

	// test alarm1, id=1, status=stored, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert1, fm, suite.namingCodes)
	if status == InError {
		suite.FailNow("Alert to fault conversion Error")
	}
	suite.Equal(Stored, status)
	suite.Equal("fault0000000001", eventFault.EventID)
	sequence = eventFault.Sequence
	suite.Equal(int64(1), sequence)

	faultName := eventFault.EventName
	suite.Equal("Fault_hspx_NodeFailure", faultName)

	// test alarm1, id=1, status=alreadyExist, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert1, fm, suite.namingCodes)
	suite.Equal(AlreadyExist, status)
	suite.Equal("fault0000000001", eventFault.EventID)
	sequence = eventFault.Sequence
	// sequence incrementation is done by handler if POST sending is OK
	suite.Equal(int64(1), sequence)

	// test alarm2, id=2, status=stored, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert2, fm, suite.namingCodes)
	if status == InError {
		suite.FailNow("Alert to fault conversion Error")
	}
	suite.Equal(Stored, status)
	suite.Equal("fault0000000002", eventFault.EventID)
	sequence = eventFault.Sequence
	suite.Equal(int64(1), sequence)

	faultName = eventFault.EventName
	suite.Equal("Fault_hspx_NodeFailure", faultName)

	// test alarm2, id=2, status=alreadyExist, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert2, fm, suite.namingCodes)
	suite.Equal(AlreadyExist, status)
	suite.Equal("fault0000000002", eventFault.EventID)
	sequence = eventFault.Sequence
	// sequence incrementation is done by handler if POST sending is OK
	suite.Equal(int64(1), sequence)

}

func (suite *ConvertTestSuite) TestConvertAlertClearOK() {
	//var id int32
	var eventFault *govel.EventFault
	var status StatusResult
	var commit CommitFunc
	var fm = NewFaultManager(&suite.confEvent)
	var sequence int64

	// test alarm1, id=1, status=stored, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert1, fm, suite.namingCodes)
	if status == InError {
		suite.FailNow("Alert to fault conversion Error")
	}
	suite.Equal(Stored, status)
	suite.Equal("fault0000000001", eventFault.EventID)

	// test alarm1, id=2, status=stored, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert2, fm, suite.namingCodes)
	if status == InError {
		suite.FailNow("Alert to fault conversion Error")
	}
	suite.Equal(Stored, status)
	suite.Equal("fault0000000002", eventFault.EventID)

	// test alarm1, id=1, status=cleared, sequence=1
	suite.alert1.Status = "resolved"
	suite.alert1.Labels["severity"] = "normal"

	status, eventFault, commit = AlertToFault(suite.alert1, fm, suite.namingCodes)
	suite.NoError(commit())
	suite.Equal(Cleared, status)
	suite.Equal("fault0000000001", eventFault.EventID)
	sequence = eventFault.Sequence
	// sequence incrementation is done by handler if POST sending is OK
	suite.Equal(int64(1), sequence)

	// set alarm1 on
	suite.alert1.Status = "firing"
	suite.alert1.Labels["severity"] = "critical"

	// test alarm2, id=2, status=cleared, sequence=1
	suite.alert2.Status = "resolved"
	suite.alert2.Labels["severity"] = "normal"

	status, eventFault, commit = AlertToFault(suite.alert2, fm, suite.namingCodes)
	suite.NoError(commit())
	suite.Equal(Cleared, status)
	suite.Equal("fault0000000002", eventFault.EventID)
	sequence = eventFault.Sequence
	// sequence incrementation is done by handler if POST sending is OK
	suite.Equal(int64(1), sequence)

	// set alarm2 on
	suite.alert2.Status = "firing"
	suite.alert2.Labels["severity"] = "critical"

	// test alarm1, id=3, status=stored, sequence=1

	status, eventFault, _ = AlertToFault(suite.alert1, fm, suite.namingCodes)
	suite.Equal(Stored, status)
	suite.Equal("fault0000000003", eventFault.EventID)
	sequence = eventFault.Sequence
	suite.Equal(int64(1), sequence)
}

func (suite *ConvertTestSuite) TestConvertAlert100Ok() {
	//var id int32
	var eventFault *govel.EventFault
	var status StatusResult
	var fm = NewFaultManager(&suite.confEvent)
	var startEpoch int64
	//var sequence int64

	// test raised alarm3, id=1, status=stored, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert3, fm, suite.namingCodes)
	if status == InError {
		suite.FailNow("Alert to fault conversion Error")
	}
	startEpoch = eventFault.StartEpochMicrosec
	suite.Equal(Stored, status)
	suite.Equal("fault0000000001", eventFault.EventID)
	suite.Equal("Fault_hspx_VudrBegCommunicationAuditFailure", eventFault.EventName)
	suite.Equal("VudrBegCommunicationAuditFailure", eventFault.AlarmCondition)
	suite.Equal("Low", string(eventFault.Priority))
	suite.Equal("MINOR", string(eventFault.EventSeverity))
	suite.Equal("Active", string(eventFault.VfStatus))
	suite.Equal("hspx", eventFault.NfNamingCode)
	suite.Equal("", eventFault.NfcNamingCode)
	suite.Equal(govel.SourceVirtualMachine, eventFault.EventSourceType)
	suite.Equal("TEST3", eventFault.SourceName)
	suite.Equal(int64(1), eventFault.Sequence)
	suite.Equal("VudrName", eventFault.AlarmAdditionalInformation[0].Name)
	suite.Equal("VUDR_1", eventFault.AlarmAdditionalInformation[0].Value)
	suite.Equal("BegName", eventFault.AlarmAdditionalInformation[1].Name)
	suite.Equal("BEG_1", eventFault.AlarmAdditionalInformation[1].Value)
	suite.Equal("On demand or periodic audit between SDME and PGW End Point SPML has failed", eventFault.SpecificProblem)

	//test cleared alarm3, status=cleared, sequence=1
	suite.alert3.Status = "resolved"
	suite.alert3.Labels["severity"] = "NORMAL"
	status, eventFault, _ = AlertToFault(suite.alert3, fm, suite.namingCodes)
	suite.Equal(Cleared, status)
	suite.Equal("fault0000000001", eventFault.EventID)
	suite.Equal("Fault_hspx_VudrBegCommunicationAuditEndOfFailure", eventFault.EventName)
	suite.Equal("VudrBegCommunicationAuditEndOfFailure", eventFault.AlarmCondition)
	suite.Equal("Normal", string(eventFault.Priority))
	suite.Equal("NORMAL", string(eventFault.EventSeverity))
	suite.Equal("Active", string(eventFault.VfStatus))
	suite.Equal("hspx", eventFault.NfNamingCode)
	suite.Equal("", eventFault.NfcNamingCode)
	suite.Equal(govel.SourceVirtualMachine, eventFault.EventSourceType)
	suite.Equal("TEST3", eventFault.SourceName)
	// sequence incrementation is done by handler if POST sending is OK
	suite.Equal(int64(1), eventFault.Sequence)
	suite.Equal("VudrName", eventFault.AlarmAdditionalInformation[0].Name)
	suite.Equal("VUDR_1", eventFault.AlarmAdditionalInformation[0].Value)
	suite.Equal("BegName", eventFault.AlarmAdditionalInformation[1].Name)
	suite.Equal("BEG_1", eventFault.AlarmAdditionalInformation[1].Value)
	suite.Equal("On demand or periodic audit between SDME and PGW End Point now ok", eventFault.SpecificProblem)
	suite.Equal(startEpoch, eventFault.StartEpochMicrosec)

}

func (suite *ConvertTestSuite) TestConvertAlert302Ok() {
	//var id int32
	var eventFault *govel.EventFault
	var status StatusResult
	var fm = NewFaultManager(&suite.confEvent)

	// test raised alarm3, id=1, status=stored, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert4, fm, suite.namingCodes)
	if status == InError {
		suite.FailNow("Alert to fault conversion Error")
	}

	suite.Equal(Stored, status)
	suite.Equal("fault0000000001", eventFault.EventID)
	suite.Equal("Fault_hspx_FileSystemOccupancyCrossedLowThreshold", eventFault.EventName)
	suite.Equal("FileSystemOccupancyCrossedLowThreshold", eventFault.AlarmCondition)
	suite.Equal("Low", string(eventFault.Priority))
	suite.Equal("MINOR", string(eventFault.EventSeverity))
	suite.Equal("Active", string(eventFault.VfStatus))
	suite.Equal("hspx", eventFault.NfNamingCode)
	suite.Equal("etl", eventFault.NfcNamingCode)
	suite.Equal(govel.SourceVirtualMachine, eventFault.EventSourceType)
	suite.Equal("mjves-ope-1", eventFault.SourceName)
	suite.Equal(int64(1), eventFault.Sequence)
	suite.Equal("FileSystemName", eventFault.AlarmAdditionalInformation[0].Name)
	suite.Equal("/var/log/audit", eventFault.AlarmAdditionalInformation[0].Value)
	suite.Equal("ObservedThreshold", eventFault.AlarmAdditionalInformation[1].Name)
	suite.Equal("85.1234", eventFault.AlarmAdditionalInformation[1].Value)
	suite.Equal("File system almost full", eventFault.SpecificProblem)

}

func (suite *ConvertTestSuite) TestConvertAlert305Ok() {
	//var id int32
	var eventFault *govel.EventFault
	var status StatusResult
	var fm = NewFaultManager(&suite.confEvent)

	// test raised alarm3, id=1, status=stored, sequence=1
	status, eventFault, _ = AlertToFault(suite.alert5, fm, suite.namingCodes)
	if status == InError {
		suite.FailNow("Alert to fault conversion Error")
	}

	suite.Equal(Stored, status)
	suite.Equal("fault0000000001", eventFault.EventID)
	suite.Equal("Fault_hspx_MemoryOccupancyCrossedLowThreshold", eventFault.EventName)
	suite.Equal("MemoryOccupancyCrossedLowThreshold", eventFault.AlarmCondition)
	suite.Equal("Low", string(eventFault.Priority))
	suite.Equal("MINOR", string(eventFault.EventSeverity))
	suite.Equal("Active", string(eventFault.VfStatus))
	suite.Equal("hspx", eventFault.NfNamingCode)
	suite.Equal("etl", eventFault.NfcNamingCode)
	suite.Equal(govel.SourceVirtualMachine, eventFault.EventSourceType)
	suite.Equal("mjves-ope-1", eventFault.SourceName)
	suite.Equal(int64(1), eventFault.Sequence)
	suite.Equal("ObservedThreshold", eventFault.AlarmAdditionalInformation[0].Name)
	suite.Equal("82.4321", eventFault.AlarmAdditionalInformation[0].Value)
	suite.Equal("Memory high occupancy", eventFault.SpecificProblem)

}
