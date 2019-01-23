package agent

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"github.com/nokia/onap-vespa/ves-agent/config"
	"github.com/nokia/onap-vespa/ves-agent/convert"
	"github.com/nokia/onap-vespa/govel"
	"github.com/nokia/onap-vespa/ves-agent/ha"
	"github.com/nokia/onap-vespa/ves-agent/heartbeat"
	"github.com/nokia/onap-vespa/ves-agent/metrics"
	"github.com/nokia/onap-vespa/ves-agent/rest"
	"github.com/nokia/onap-vespa/ves-agent/scheduler"

	log "github.com/sirupsen/logrus"
)

// Agent sends heartbeat, measurement and fault events.
// It initializes schedulers to trigger heartbeat and metric events.
// The schedulers are updated on new heartbeat and measurement interval.
// It initializes the AlertReceiver server to receive and handle Alert event from prometheus.
type Agent struct {
	measSched, hbSched           *scheduler.Scheduler
	measTimer, hbTimer           *time.Timer
	measIntervalCh, hbIntervalCh <-chan time.Duration
	alertCh                      chan rest.MessageFault
	fm                           *convert.FaultManager
	alertRoute                   rest.Route
	state                        *ha.Cluster
	namingCodes                  map[string]string
}

// NewAgent initializes schedulers to trigger heartbeat and metric events.
// It initializes the AlertReceiver server to receive Alert event from prometheus.
func NewAgent(conf *config.VESAgentConfiguration) *Agent {
	state, err := ha.NewCluster(conf.DataDir, conf.Cluster, ha.NewInMemState())
	if err != nil {
		log.Panic(err)
	}
	namingCodes := initNfcNamingCode(conf.Event.NfcNamingCodes)

	log.Info("Create measurement scheduler")
	// Create a new Scheduler used to trigger the measurements collector
	measSched := initMeasScheduler(conf, namingCodes, state)

	log.Info("Create heartbeat scheduler")
	// Create a new Scheduler used to trigger the heartbeat events
	hbSched := initHbScheduler(&conf.Event, conf.Heartbeat.DefaultInterval, namingCodes, state)

	// create a FaultManager
	fm := convert.NewFaultManagerWithState(&conf.Event, state)
	// declare the AlertReceiver route where the server will be listening
	alertRoute := rest.Route{
		Name:        "AlertReceiver",
		Method:      "POST",
		Pattern:     conf.AlertManager.Path,
		HandlerFunc: nil,
	}

	return &Agent{
		measSched:   measSched,
		hbSched:     hbSched,
		fm:          fm,
		alertRoute:  alertRoute,
		state:       state,
		namingCodes: namingCodes,
	}
}

// Stats exposes some internal stats. This MUST be used ONLY
// for debugging purpose, and MUST NOT be considered stable
func (agent *Agent) Stats() map[string]interface{} {
	return map[string]interface{}{
		"raft": agent.state.Stats(),
	}
}

func initMeasScheduler(conf *config.VESAgentConfiguration, namingCodes map[string]string, state ha.AgentState) *scheduler.Scheduler {
	// Creates a new measurements collector
	prom, err := metrics.NewCollectorWithState(&conf.Measurement, &conf.Event, namingCodes, state)
	if err != nil {
		log.Panic(err)
	}
	measSched := scheduler.NewSchedulerWithState("measurements", prom, conf.Measurement.DefaultInterval, state)
	return measSched
}

func initHbScheduler(conf *govel.EventConfiguration, defaultInterval time.Duration, namingCodes map[string]string, state ha.AgentState) *scheduler.Scheduler {
	// Creates a new heartbeat monitor
	hbMonitor, err := heartbeat.NewMonitorWithState(conf, namingCodes, state)
	if err != nil {
		log.Panic(err)
	}
	hbSched := scheduler.NewSchedulerWithState("heartbeats", hbMonitor, defaultInterval, state)
	return hbSched
}

// initNfcNamingCode extract the vnfcNamingCode from vnfcName
func initNfcNamingCode(nfcNamingCodes []govel.NfcNamingCode) map[string]string {
	namingCodes := make(map[string]string)
	for _, nfcCode := range nfcNamingCodes {
		for _, vnfc := range nfcCode.Vnfcs {
			namingCodes[vnfc] = nfcCode.Type
		}
	}
	return namingCodes
}

// StartAgent registers to heartbeat and measurement interval changed events, and triggers the events.
// It initializes the AlertReceiver server to receive and handle Alert event from prometheus.
func (agent *Agent) StartAgent(bind string, ves govel.VESCollectorIf) {
	agent.listen(bind, ves)
	agent.serve(ves)
}

func (agent *Agent) listen(bind string, ves govel.VESCollectorIf) {
	log.Info("Setup measurement collection routine")
	// Subscribe to measurement interval changed events
	agent.measIntervalCh = ves.NotifyMeasurementIntervalChanged(make(chan time.Duration, 1024))

	log.Info("Setup heartbeat routine")
	// Subscribe to heartbeat interval changed events
	agent.hbIntervalCh = ves.NotifyHeartbeatIntervalChanged(make(chan time.Duration, 1024))

	log.Info("Setup alert receiver server")
	// Setup the AlertReceiver and subscribe to alert events
	agent.notifyAlertEventReceived(bind)
}

func (agent *Agent) notifyAlertEventReceived(bind string) {
	// attach the AlertReceiver handler to the alert route managed by server
	agent.alertCh = make(chan rest.MessageFault, 1024)
	agent.alertRoute.HandlerFunc = rest.AlertReceiver(agent.alertCh)

	routes := []rest.Route{
		agent.alertRoute,
		{Name: "Stats", Method: "GET", Pattern: "/stats", HandlerFunc: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			if err := enc.Encode(agent.Stats()); err != nil {
				log.Errorf("HTTP Handler - Cannot write stats: %s", err.Error())
			}
		})},
	}

	// create an unstarted new server to receive http POST from prometheus
	alertHandler := rest.NewServer(routes)
	// start server
	go rest.StartServer(bind, alertHandler)
}

func (agent *Agent) serve(ves govel.VESCollectorIf) {
	for {
		// Wait to become cluster's leader
		log.Info("Waiting to obtain cluster leadership")
		for !agent.followerStep() {
		}
		log.Info("Gained cluster leadership")

		// Setup schedulers timers
		agent.measTimer = agent.measSched.WaitChan()
		agent.hbTimer = agent.hbSched.WaitChan()
		// Run leadership steps until we loose leader state
		for agent.leaderStep(ves) {
		}
		log.Info("Lost cluster leadership")
		agent.measTimer.Stop()
		agent.hbTimer.Stop()
	}
}

func (agent *Agent) followerStep() bool {
	select {
	case fault := <-agent.alertCh:
		fault.Response <- errors.New("Not the leader")
		close(fault.Response)
	case leader := <-agent.state.LeaderCh():
		return leader
	}
	return false
}

func (agent *Agent) leaderStep(ves govel.VESCollectorIf) bool {
	// Demultiplex events
	select {
	case measInterval := <-agent.measIntervalCh:
		// Measurements collection interval changed event
		agent.handleMeasurementIntervalChanged(measInterval)
	case hbInterval := <-agent.hbIntervalCh:
		// Heartbeat interval changed event
		agent.handleHeartbeatIntervalChanged(hbInterval)
	case fault := <-agent.alertCh:
		//alert received event
		agent.handleAlertReceived(ves, fault)
	case <-agent.measTimer.C:
		// It's time to collect and send some measurements
		agent.triggerMeasurementEvent(ves)
	case <-agent.hbTimer.C:
		// It's time to send the heartbeat
		agent.triggerHeatbeatEvent(ves)
	case leader := <-agent.state.LeaderCh():
		return leader
	}
	return true
}

func (agent *Agent) handleMeasurementIntervalChanged(interval time.Duration) {
	if err := agent.measSched.SetInterval(interval); err != nil {
		log.Errorf("Cannot update measurement interval : %s", err.Error())
		return
	}
	agent.measTimer.Stop()
	agent.measTimer = agent.measSched.WaitChan()
}

func (agent *Agent) handleHeartbeatIntervalChanged(interval time.Duration) {
	if err := agent.hbSched.SetInterval(interval); err != nil {
		log.Errorf("Cannot update heartbeat interval : %s", err.Error())
		return
	}
	agent.hbTimer.Stop()
	agent.hbTimer = agent.hbSched.WaitChan()
}

func (agent *Agent) handleAlertReceived(ves govel.VESCollectorIf, messageFault rest.MessageFault) {
	status, eventFault, commitFunc := convert.AlertToFault(messageFault.Alert, agent.fm, agent.namingCodes)
	if status == convert.InError || status == convert.NotExist {
		log.Warningln("!!!error in ConvertToFault process")
		if status == convert.InError {
			messageFault.Response <- errors.New("Cannot convert Fault to VES event")
		}
	} else {

		if err := ves.PostEvent(eventFault); err != nil {
			log.Error("Cannot post fault: ", err.Error())
			// Send result to fault handler.
			messageFault.Response <- err
		} else {
			// Commit the alert if successfully sent
			if err := commitFunc(); err != nil {
				messageFault.Response <- err
			}
		}
	}
	close(messageFault.Response)
}

func (agent *Agent) triggerMeasurementEvent(ves govel.VESCollectorIf) {
	triggerScheduler(agent.measSched, &agent.measTimer, func(res interface{}) error {
		return ves.PostBatch(res.(metrics.EventMeasurementSet).Batch())
	})
}

func (agent *Agent) triggerHeatbeatEvent(ves govel.VESCollectorIf) {
	triggerScheduler(agent.hbSched, &agent.hbTimer, func(res interface{}) error {
		return ves.PostEvent(res.(govel.Event))
	})
}

func triggerScheduler(sched *scheduler.Scheduler, timer **time.Timer, f func(interface{}) error) {
	res, err := sched.Step()
	if err != nil {
		log.Errorf("Cannot trigger scheduler %s: %s", sched.Name(), err.Error())
		// Setup a retry timer
		*timer = time.NewTimer(10 * time.Second)
		return
	}
	if err = f(res); err == nil {
		// Acknowledge the scheduler interval(s) if send is successful
		if err := sched.Ack(); err != nil {
			log.Errorf("Cannot acknowledge scheduler execution: %s", err.Error())
			return
		}
		// Set timer to the next interval
		*timer = sched.WaitChan()
	} else {
		// If Post to active ves collector failed: setup a retry timer before trying to second ves collector
		*timer = time.NewTimer(10 * time.Second)
	}
}
