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

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"github.com/nokia/onap-vespa/govel"
	"github.com/nokia/onap-vespa/govel/schema"

	log "github.com/sirupsen/logrus"
)

// CommandList change request
type commandListRequest struct {
	CommandList []govel.Command `json:"commandList"`
}

type publishEventRequest struct {
	Event map[string]interface{} `json:"event"`
}

type publishBatchRequest struct {
	EventList []map[string]interface{} `json:"eventList,omitempty"`
}

func checkAuth(req *http.Request) error {
	if reqUser, reqPass, ok := req.BasicAuth(); !ok {
		return errors.New("Authentication needed")
	} else if reqUser != *user || reqPass != *pass {
		return errors.New("Authentication failed")
	}
	log.Info("Authentication success")
	return nil
}

func close(closer io.Closer) {
	if closer == nil {
		return
	}
	if err := closer.Close(); err != nil {
		log.Error(err.Error())
	}
}

func sendVESReply(w http.ResponseWriter) error {
	if len(commandList) > 0 {
		log.Debugf("Replying with commands %+v", commandList)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		resp := govel.VESResponse{
			CommandList: commandList,
		}
		// fmt.Printf("resp : %+v\n", resp)
		commandList = make([]govel.Command, 0)
		renc := json.NewEncoder(w)
		renc.SetIndent("", "\t")
		if err := renc.Encode(&resp); err != nil {
			return err
		}
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
	return nil
}

func decodeAndValidateJSON(req *http.Request, data interface{}) error {
	if req.Body == nil {
		return errors.New("Request has no body")
	}
	defer close(req.Body)
	buf := bytes.Buffer{}
	if _, err := buf.ReadFrom(req.Body); err != nil {
		return err
	}
	if *eventMaxSize > 0 && buf.Len() > *eventMaxSize {
		return fmt.Errorf("Received event is too big (has %d bytes, limit is %d bytes)", buf.Len(), *eventMaxSize)
	}
	// Copy received data into a new byte array
	bytes := append(([]byte)(nil), buf.Bytes()...)
	buf.Reset()
	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}
	// Pretty print JSON
	if err := json.Indent(&buf, bytes, "", "  "); err != nil {
		return err
	}
	log.Infof("\n%s", buf.String())
	// Check schema
	if err := schema.V2841().Validate(data); err != nil {
		return errors.New("Schema check failed: " + err.Error())
	}
	log.Info("Schema validation succeeded")
	return nil
}

func handlePostEvent(w http.ResponseWriter, req *http.Request) error {
	mutex.Lock()
	defer mutex.Unlock()
	log.Info("****************** Received event *******************")
	if err := checkAuth(req); err != nil {
		return err
	}
	event := publishEventRequest{}
	if err := decodeAndValidateJSON(req, &event); err != nil {
		return err
	}
	appendEvent(event.Event)
	return sendVESReply(w)
}

func handlePostBatch(w http.ResponseWriter, req *http.Request) error {
	mutex.Lock()
	defer mutex.Unlock()
	log.Info("****************** Received batch *******************")
	if err := checkAuth(req); err != nil {
		return err
	}
	batch := publishBatchRequest{}
	if err := decodeAndValidateJSON(req, &batch); err != nil {
		return err
	}
	for _, event := range batch.EventList {
		appendEvent(event)
	}
	stats.Batch++
	return sendVESReply(w)
}

func handleGetEvents(w http.ResponseWriter, req *http.Request) error {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	domain := req.URL.Query().Get("domain")
	entity := req.URL.Query().Get("entity")
	clear := req.URL.Query().Get("clear")
	filters := make(map[string]string)
	filters["domain"] = domain
	filters["reportingEntityName"] = entity
	filteredEvents := filterEvents(filters, false)
	if err := encoder.Encode(filteredEvents); err != nil {
		return err
	}
	if clear == "1" {
		clearEvents(filters)
	}
	return nil
}

func handleClearEvents(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	clearEvents(nil)
}

func handleSetCommandList(w http.ResponseWriter, req *http.Request) error {
	mutex.Lock()
	defer mutex.Unlock()

	if req.Body == nil {
		return errors.New("Request has no body")
	}
	defer close(req.Body)
	commands := commandListRequest{}
	if err := json.NewDecoder(req.Body).Decode(&commands); err != nil {
		return err
	}
	for _, cmd := range commands.CommandList {
		addCommand(cmd)
	}
	return nil
}

func handleGetStats(w http.ResponseWriter, req *http.Request) error {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(&stats)
}

func errorWrapper(hdl func(w http.ResponseWriter, req *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := hdl(w, req); err != nil {
			stats.Errors++
			log.Errorf("Invalid request: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}
