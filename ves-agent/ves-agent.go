package main

import (
	"os"
	"os/signal"
	"runtime"
	"github.com/nokia/onap-vespa/ves-agent/agent"
	"github.com/nokia/onap-vespa/ves-agent/config"
	"github.com/nokia/onap-vespa/govel"

	log "github.com/sirupsen/logrus"
)

// List of build information
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// initLogging initilize the logger.
// If `debug` is true, then debug traces
// are activated
func initLogging(debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug traces activated")
		// log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// launchVES launch the routine for
// - metric collection
// - heartbeat events
// - alert received events
func launchVES(ves govel.VESCollectorIf, conf *config.VESAgentConfiguration) {
	log.Info("Starting VES routine")
	defer log.Fatal("VES routine exited")
	agent := agent.NewAgent(conf)
	agent.StartAgent(conf.AlertManager.Bind, ves)
}

func main() {

	var conf config.VESAgentConfiguration
	if err := config.InitConf(&conf); err != nil {
		log.Fatal("Cannot read config file: ", err.Error())
	}

	initLogging(conf.Debug)

	log.Infof("Starting VES Agent version %s", version)
	log.Infof("Version=%s, Commit=%s, Date=%s, Go version=%s", version, commit, date, runtime.Version())

	ves, err := govel.NewCluster(&conf.PrimaryCollector, &conf.BackupCollector, &conf.Event, conf.CaCert)
	if err != nil {
		log.Fatal("Cannot initialize VES connection: ", err.Error())
	}

	go launchVES(ves, &conf)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Infof("Stopping VES Agent version %s", version)
}
