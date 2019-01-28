package main

// Simulator internal state
import (
	"fmt"
	"mime"
	"sync"
	"github.com/nokia/onap-vespa/govel"

	log "github.com/sirupsen/logrus"

	"github.com/gobuffalo/packr"
)

//go:generate packr -z

var (
	events      = make([]map[string]interface{}, 0)
	commandList = make([]govel.Command, 0)
	mutex       = sync.RWMutex{}
	assets      = packr.NewBox("assets")
	stats       = struct {
		Batch uint64 `json:"batch"`
		// Heartbeat uint64
		// Faults    uint64
		// Metrics   uint64
		Errors     uint64 `json:"errors"`
		LastSender string `json:"last_sender"`
	}{}
)

func init() {
	for _, ext := range []string{".yml", ".yaml"} {
		if err := mime.AddExtensionType(ext, "text/plain"); err != nil {
			log.Panic(err)
		}
	}
}

func appendEvent(event map[string]interface{}) {
	events = append(events, event)
	if *maxEventsKeep > 0 && len(events) > *maxEventsKeep {
		log.Warn("Max event buffer size reached. Dismissing oldest events")
		// Slide down array's elements
		copy(events, events[len(events)-*maxEventsKeep:])
	}
	stats.LastSender = event["commonEventHeader"].(map[string]interface{})["reportingEntityName"].(string)
}

// Add a command to send to next reply
func addCommand(cmd govel.Command) {
	log.Debugf("Adding command %+v", cmd)
	for i, command := range commandList {
		if command.CommandType == cmd.CommandType {
			commandList[i] = cmd
			return
		}
	}
	commandList = append(commandList, cmd)
}

// Find received events for specifics filters
func filterEvents(filters map[string]string, exclude bool) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, evt := range events {
		check := true
		for filter := range filters {
			if filters[filter] != "" {
				value := fmt.Sprint(evt["commonEventHeader"].(map[string]interface{})[filter])
				check = check && value == filters[filter]
			}
		}
		if check != exclude {
			res = append(res, evt)
		}
	}
	log.Debugf("Found %d events", len(res))
	return res
}

func clearEvents(filters map[string]string) {
	log.Debugf("Clearing events")
	events = filterEvents(filters, true)
}
