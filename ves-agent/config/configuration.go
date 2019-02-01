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

package config

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// List of required configuration parameters
var (
	requiredConfs = []string{"PrimaryCollector.User", "PrimaryCollector.Password"}
	configName    = "ves-agent"
	configPaths   = []string{"/etc/ves-agent/", "."}
)

func setFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP("PrimaryCollector.FQDN", "f", "localhost", "VES Collector FQDN")
	flagSet.IntP("PrimaryCollector.Port", "p", 8443, "VES Collector Port")
	flagSet.StringP("PrimaryCollector.Topic", "t", "", "VES Collector Topic")
	flagSet.StringP("PrimaryCollector.User", "u", "", "VES Username")
	flagSet.StringP("PrimaryCollector.Password", "k", "", "VES Password")
	flagSet.String("PrimaryCollector.PassPhrase", "", "VES PassPhrase")
	flagSet.String("BackupCollector.FQDN", "", "VES Collector FQDN")
	flagSet.Int("BackupCollector.Port", 0, "VES Collector Port")
	flagSet.String("BackupCollector.Topic", "", "VES Collector Topic")
	flagSet.String("BackupCollector.User", "", "VES Username")
	flagSet.String("BackupCollector.Password", "", "VES Password")
	flagSet.String("BackupCollector.PassPhrase", "", "VES PassPhrase")
	flagSet.DurationP("Heartbeat.DefaultInterval", "i", 60*time.Second, "VES heartbeat interval")
	flagSet.StringP("Measurement.DomainAbbreviation", "d", "Measurement", "Domain Abbreviation")
	flagSet.DurationP("Measurement.DefaultInterval", "m", 300*time.Second, "Measurement interval")
	flagSet.String("Measurement.Prometheus.Address", "http://localhost:9090", "Base url to of Prometheus server's API")
	flagSet.Duration("Measurement.MaxBufferingDuration", time.Hour, "Maximum timeframe size of buffering")
	flagSet.IntP("Event.MaxSize", "s", 200, "Max Event Size")
	retrieveReportingEntityName(flagSet)
	flagSet.DurationP("Event.RetryInterval", "r", 10*time.Second, "VES heartbeat retry interval")
	flagSet.IntP("Event.MaxMissed", "a", 3, "Missed heartbeats until switching collector")
	flagSet.String("AlertManager.Bind", "localhost:9095", "Alert Manager Bind address")
	flagSet.String("AlertManager.Path", "/alerts", "Alert Manager Path")
	flagSet.String("AlertManager.User", "", "Alert Manager Username")
	flagSet.String("AlertManager.Password", "", "Alert Manager Password")
	flagSet.String("Cluster.ID", "", "Override the cluster's node ID")
	flagSet.StringP("DataDir", "D", "/var/lib/ves-agent/data", "Path to directory where to store data")
	flagSet.Bool("Debug", false, "Activate debug traces")
}

// InitConf initilize the config store from config file, env and cli variables.
func InitConf(conf *VESAgentConfiguration) error {

	//bind env variable
	viper.SetEnvPrefix("ves")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// Viper will check for an environment variable any time a viper.Get request is made
	viper.AutomaticEnv()

	//read config file
	viper.SetConfigName(configName)
	for _, v := range configPaths {
		viper.AddConfigPath(v)
	}

	if err := viper.ReadInConfig(); err != nil {
		if reflect.TypeOf(err).String() == "viper.ConfigFileNotFoundError" {
			log.Warnf("Config file not found. Go on...")
		} else {
			return err
		}
	}

	//bind arguments variable
	flagSet := pflag.NewFlagSet("conf", pflag.ExitOnError)
	setFlags(flagSet)
	if err := flagSet.Parse(os.Args); err != nil {
		log.Panic(err)
	}
	if err := viper.BindPFlags(flagSet); err != nil {
		log.Panic(err)
	}

	//check required values
	for _, v := range requiredConfs {
		if !viper.IsSet(v) || viper.GetString(v) == "" {
			return errors.New("Missing required configuration parameter: " + v)
		}
	}

	// Viper will check in the following order: override, flag, env, config file, key/value store, default
	return viper.Unmarshal(conf)
}

func retrieveReportingEntityName(flagSet *pflag.FlagSet) {
	var out []byte
	var err error
	if runtime.GOOS == "windows" {
		out, err = exec.Command("hostname").Output()
	} else {
		out, err = exec.Command("hostname", "-s").Output()
	}
	if err != nil {
		log.Warnf("Cannot retrieve hostname: %s", err.Error())
	} else {
		hostname := strings.TrimSuffix(string(out), "\n")
		flagSet.StringP("Event.ReportingEntityName", "H", hostname, "Reporting entity name to add to event's header")
	}
}
