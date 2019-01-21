package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/alertmanager/template"
	log "github.com/sirupsen/logrus"
)

// MessageFault contains
// - the fault received by the server, to be sent to the collector
// - a channel to get the error (if any) after posting the fault.
type MessageFault struct {
	Alert    template.Alert
	Response chan error
}

// decodeJSON function used to extract Alerts from http datas
func decodeJSON(resp http.ResponseWriter, req *http.Request) template.Alerts {
	//fmt.Println("decode json request")

	data := template.Data{}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Errorf("Bad request from %s: %s\n", req.RemoteAddr, err.Error())
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}
	alertsmsg := data.Alerts

	//fmt.Printf("data received: %+v\n", alertsmsg)
	return alertsmsg
}

// AlertReceiver is an handler to manage  http POST alert
func AlertReceiver(alertCh chan MessageFault) http.Handler {
	hd1 := func(resp http.ResponseWriter, req *http.Request) error {
		contentType := req.Header.Get("Content-Type")
		if contentType != "application/json" {
			log.Errorf("content-type %s not managed \n", contentType)
			//resp.WriteHeader(http.StatusInternalServerError)
			return errors.New("content-type %s not managed")
		}
		alertsmsg := decodeJSON(resp, req)
		for i := range alertsmsg {
			errorCh := make(chan error)
			message := MessageFault{Alert: alertsmsg[i], Response: errorCh}
			// Non blocking write, to avoid a dead lock situation
			select {
			case alertCh <- message:
				//wait for PostEvent result
				err := <-errorCh
				if err != nil {
					log.Errorf("Cannot process alert: %s", err.Error())
					return err
				}
			default:
				err := fmt.Errorf("Alert %s could not be sent to a channel", alertsmsg[i].Labels["alertname"])
				log.Warn(err.Error())
				return err
			}
		}
		return nil
	}
	return errorWrapper(hd1)
}
