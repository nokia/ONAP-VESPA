package convert

import (
	"github.com/nokia/onap-vespa/govel"
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/alertmanager/template"
	log "github.com/sirupsen/logrus"
)

const domain = "Fault"

var severityToPriority = map[govel.Severity]govel.EventPriority{
	govel.SeverityCritical: govel.PriorityHigh,
	govel.SeverityMajor:    govel.PriorityMedium,
	govel.SeverityMinor:    govel.PriorityLow,
	govel.SeverityWarning:  govel.PriorityLow,
	govel.SeverityNormal:   govel.PriorityNormal,
}

// CommitFunc is a function used to commit a fault conversion operation
type CommitFunc func() error

func mustNotCall() error {
	log.Panic("Commit function should not have been called")
	return nil
}

// AlertToFault convert Alert to VES fault and store it into a map;
// return status, id and a function used to finalize the operation.
// The function returned must be called after having successfully sent the alert to VES
// status could be: inError,alreadyExist, stored, cleared
func AlertToFault(alert template.Alert, fm *FaultManager, namingCodes map[string]string) (StatusResult, *govel.EventFault, CommitFunc) {
	var storeStatus StatusResult
	var alertName string
	var specificProblem string
	//var eventFault *govel.EventFault
	var id int32

	label := alert.Labels
	annotations := alert.Annotations
	log.Debugln("convert alert to VES event fault: ", label["alertname"])
	severity := govel.Severity(strings.ToUpper(label["severity"]))
	priority, exist := severityToPriority[severity]
	// check severity value consistence
	if !exist {
		log.Debugln("Error in severityToPriority convert for severity : " + label["severity"])
		return InError, nil, mustNotCall
	}
	// depending of the alarm type (generic or specific) sourceName can be system_name or VNFC
	sourceName := label["VNFC"]
	if sourceName == "" {
		sourceName = label["system_name"]
	}

	//faultName := label["id"] + "_" + annotations["service"] + "_" + sourceName
	faultName := buildFaultName(label["id"], annotations["service"], sourceName)
	nfNamingCode := fm.GetEventConf().NfNamingCode

	if alert.Status == "resolved" {
		storeStatus, id = fm.clearFault(faultName)
		alertName = annotations["clearAlertName"]
		specificProblem = annotations["clearDescription"]
		severity = "NORMAL"
	} else {
		storeStatus, id = fm.storeFault(faultName)
		alertName = label["alertname"]
		specificProblem = annotations["description"]
	}

	if storeStatus == InError || storeStatus == NotExist {
		return storeStatus, nil, mustNotCall
	}

	eventName := domain + "_" + nfNamingCode + "_" + alertName
	vesID := fmt.Sprintf("fault%010d", id)
	eventFault := govel.NewFault(
		eventName,
		vesID,
		alertName,
		specificProblem,
		priority,
		severity,
		govel.SourceVirtualMachine,
		govel.StatusActive,
		sourceName)

	eventFault.NfNamingCode = nfNamingCode
	eventFault.NfcNamingCode = namingCodes[sourceName]

	eventFault.Sequence = fm.state.GetFaultSn(id)

	// startEpoch is initialized at first event
	// same startEpoch is used in other case
	if storeStatus == Stored {
		if err := fm.GetFaultState().SetFaultStartEpoch(id, eventFault.StartEpochMicrosec); err != nil {
			log.Error(err.Error())
			return InError, nil, nil
		}
	} else {
		eventFault.StartEpochMicrosec = fm.GetFaultState().GetFaultStartEpoch(id)
	}

	if aaiMapping, exist := annotations["aaiMapping"]; exist {
		additionalInfos := buildAdditionalInfos(annotations["service"], aaiMapping)
		if len(additionalInfos) != 0 {
			eventFault.AlarmAdditionalInformation = additionalInfos
		}
	}

	log.Debugf("AlertToFault success for id %s sequence %d: \n", vesID, eventFault.Sequence)

	var commitFunc CommitFunc
	if storeStatus == Cleared {
		commitFunc = func() error {
			log.Debugf("Delete Fault %s with id %d in storage \n ", faultName, id)
			return fm.state.DeleteFaultInStorage(faultName)
		}
	} else {
		commitFunc = func() error {
			return fm.state.IncrementFaultSn(id)
		}
	}

	return storeStatus, eventFault, commitFunc
}

// buildFaultName built faultName
// faultName = <id>_<service>_<sourceName>
func buildFaultName(id string, service string, sourceName string) string {
	var serv string
	// check if service contains dynamic information
	// if yes dont used it to build faultName
	if strings.Contains(service, "=") {
		index := strings.Index(service, "=")
		serv = service[0:index]
	} else {
		serv = service
	}
	return id + "_" + serv + "_" + sourceName
}

// GetFaultID return the id of the stored fault
func GetFaultID(vesID string) int32 {
	idString := vesID[5:]
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Warning("problem to extract id from vesID")
		return 0
	}
	return int32(id)
}

// buildAdditionalInfos build the alarmAdditionalInformation Events
// aaiMapping contains the information label; format is <label1>*<label2>_...
// service contains the information value; format is <val1>_<val2>_...
func buildAdditionalInfos(service string, aaiMapping string) []govel.EventField {
	service = strings.Replace(service, "_", "*", 1)
	service = strings.Replace(service, "=", "*", 1)
	alaInfos := strings.Split(service, "*")
	eventNames := strings.Split(aaiMapping, "_")

	var events []govel.EventField

	if len(alaInfos) == 0 || len(alaInfos) != len(eventNames)+1 {
		log.Warning("incorrect information in alert service field")
	} else {
		for i := 1; i < len(alaInfos); i++ {
			events = append(events, govel.EventField{Name: eventNames[i-1], Value: alaInfos[i]})
		}
	}
	return events
}
