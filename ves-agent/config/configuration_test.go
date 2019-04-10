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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const LineBreak = "\n"

type ConfigurationTestSuite struct {
	suite.Suite
	file *os.File
}

func TestConfiguration(t *testing.T) {
	suite.Run(t, new(ConfigurationTestSuite))
}

// Make sure the configuration files exists and is empty before each test
func (s *ConfigurationTestSuite) SetupTest() {
	var err error
	s.file, err = os.Create(configName + ".yml")
	s.NoError(err)
	if !s.NotNil(s.file) {
		s.FailNow("Cannot create configuration file")
	}
	os.Args = []string{"test"}
}

// Make sure the configuration file is closed after each test
func (s *ConfigurationTestSuite) TearDownTest() {
	s.file.Close()
	os.Remove(configName + ".yml")
	os.Clearenv()
}

func (s *ConfigurationTestSuite) TestMissingParameters() {
	var conf VESAgentConfiguration
	err := InitConf(&conf)
	for _, v := range requiredConfs {
		s.Error(err)
		s.file.WriteString(v + ": " + "default" + LineBreak)
		err = InitConf(&conf)
	}
	s.NoError(err)
	s.file.WriteString("BackupCollector.FQDN" + ": " + "default" + LineBreak)
	err = InitConf(&conf)
	s.Error(err)
	s.file.WriteString("BackupCollector.User" + ": " + "default" + LineBreak)
	err = InitConf(&conf)
	s.Error(err)
	s.file.WriteString("BackupCollector.Password" + ": " + "default" + LineBreak)
	err = InitConf(&conf)
	s.NoError(err)
}

func (s *ConfigurationTestSuite) TestDefaultsParameters() {
	s.file.WriteString("primaryCollector: " + LineBreak)
	s.file.WriteString("  user: user" + LineBreak)
	s.file.WriteString("  password: pass" + LineBreak)

	var conf VESAgentConfiguration
	err := InitConf(&conf)
	s.NoError(err)
	s.Equal("", conf.PrimaryCollector.ServerRoot)
	s.Equal("localhost", conf.PrimaryCollector.FQDN)
	s.Equal(8443, conf.PrimaryCollector.Port)
	s.Equal("", conf.PrimaryCollector.Topic)
	s.Equal("", conf.PrimaryCollector.PassPhrase)
	s.Equal("", conf.BackupCollector.FQDN)
	s.Equal(60*time.Second, conf.Heartbeat.DefaultInterval)
	s.Equal("Measurement", conf.Measurement.DomainAbbreviation)
	s.Equal(300*time.Second, conf.Measurement.DefaultInterval)
	s.Equal(200, conf.Event.MaxSize)
	s.Equal(10*time.Second, conf.Event.RetryInterval)
	s.Equal(3, conf.Event.MaxMissed)
	s.Equal("localhost:9095", conf.AlertManager.Bind)
	s.Equal(false, conf.Debug)
}

func (s *ConfigurationTestSuite) TestYamlParameters() {
	s.file.WriteString("primaryCollector: " + LineBreak)
	s.file.WriteString("  serverRoot: api" + LineBreak)
	s.file.WriteString("  fqdn: 127.0.0.1" + LineBreak)
	s.file.WriteString("  port: 30000" + LineBreak)
	s.file.WriteString("  topic: mytopic" + LineBreak)
	s.file.WriteString("  user: user" + LineBreak)
	s.file.WriteString("  password: pass" + LineBreak)
	s.file.WriteString("  passphrase: mypassphrase" + LineBreak)
	s.file.WriteString("backupCollector: " + LineBreak)
	s.file.WriteString("  fqdn: 127.0.0.2" + LineBreak)
	s.file.WriteString("  port: 40000" + LineBreak)
	s.file.WriteString("  topic: mytopic2" + LineBreak)
	s.file.WriteString("  user: user2" + LineBreak)
	s.file.WriteString("  password: pass2" + LineBreak)
	s.file.WriteString("  passphrase: mypassphrase2" + LineBreak)
	s.file.WriteString("heartbeat:" + LineBreak)
	s.file.WriteString("  defaultInterval: 5s" + LineBreak)
	s.file.WriteString("measurement: " + LineBreak)
	s.file.WriteString("  domainAbbreviation: Mfvs" + LineBreak)
	s.file.WriteString("  defaultInterval: 100s" + LineBreak)
	s.file.WriteString("event: " + LineBreak)
	s.file.WriteString("  maxSize: 100" + LineBreak)
	s.file.WriteString("  retryInterval: 6s" + LineBreak)
	s.file.WriteString("  maxMissed: 7" + LineBreak)
	s.file.WriteString("alertManager: " + LineBreak)
	s.file.WriteString("  bind: localhost:8091" + LineBreak)
	s.file.WriteString("  path: /alerts" + LineBreak)
	s.file.WriteString("  user: user" + LineBreak)
	s.file.WriteString("  password: pass" + LineBreak)
	s.file.WriteString("debug: true" + LineBreak)
	checkAll(s, false)
}

func (s *ConfigurationTestSuite) TestEnvParameters() {
	os.Setenv("VES_PRIMARYCOLLECTOR_SERVERROOT", "api")
	os.Setenv("VES_PRIMARYCOLLECTOR_FQDN", "127.0.0.1")
	os.Setenv("VES_PRIMARYCOLLECTOR_PORT", "30000")
	os.Setenv("VES_PRIMARYCOLLECTOR_TOPIC", "mytopic")
	os.Setenv("VES_PRIMARYCOLLECTOR_USER", "user")
	os.Setenv("VES_PRIMARYCOLLECTOR_PASSWORD", "pass")
	os.Setenv("VES_PRIMARYCOLLECTOR_PASSPHRASE", "mypassphrase")
	os.Setenv("VES_BACKUPCOLLECTOR_FQDN", "127.0.0.2")
	os.Setenv("VES_BACKUPCOLLECTOR_PORT", "40000")
	os.Setenv("VES_BACKUPCOLLECTOR_TOPIC", "mytopic2")
	os.Setenv("VES_BACKUPCOLLECTOR_USER", "user2")
	os.Setenv("VES_BACKUPCOLLECTOR_PASSWORD", "pass2")
	os.Setenv("VES_BACKUPCOLLECTOR_PASSPHRASE", "mypassphrase2")
	os.Setenv("VES_HEARTBEAT_DEFAULTINTERVAL", "5s")
	os.Setenv("VES_MEASUREMENT_DOMAINABBREVIATION", "Mfvs")
	os.Setenv("VES_MEASUREMENT_DEFAULTINTERVAL", "100s")
	os.Setenv("VES_EVENT_MAXSIZE", "100")
	os.Setenv("VES_EVENT_RETRYINTERVAL", "6s")
	os.Setenv("VES_EVENT_MAXMISSED", "7")
	os.Setenv("VES_DEBUG", "true")
	os.Setenv("VES_ALERTMANAGER_BIND", "localhost:8091")
	checkAll(s, false)
}

func (s *ConfigurationTestSuite) TestCLIParameters() {
	os.Args = append(os.Args, "--PrimaryCollector.ServerRoot=api", "--PrimaryCollector.FQDN=127.0.0.1", "--PrimaryCollector.Port=30000", "--PrimaryCollector.Topic=mytopic")
	os.Args = append(os.Args, "--PrimaryCollector.User=user", "--PrimaryCollector.Password=pass")
	os.Args = append(os.Args, "--BackupCollector.FQDN=127.0.0.2", "--BackupCollector.Port=40000", "--BackupCollector.Topic=mytopic2")
	os.Args = append(os.Args, "--BackupCollector.User=user2", "--BackupCollector.Password=pass2")
	os.Args = append(os.Args, "--Heartbeat.DefaultInterval=5s", "--Measurement.DomainAbbreviation=Mfvs", "--Measurement.DefaultInterval=100s")
	os.Args = append(os.Args, "--Event.MaxSize=100", "--Event.RetryInterval=6s", "--Event.MaxMissed=7")
	os.Args = append(os.Args, "--Debug")
	checkAll(s, true)
}

func (s *ConfigurationTestSuite) TestCLIShortParameters() {
	os.Args = append(os.Args, "--PrimaryCollector.ServerRoot=api", "-f=127.0.0.1", "-p=30000", "-t=mytopic")
	os.Args = append(os.Args, "-u=user", "-k=pass")
	os.Args = append(os.Args, "--BackupCollector.FQDN=127.0.0.2", "--BackupCollector.Port=40000", "--BackupCollector.Topic=mytopic2")
	os.Args = append(os.Args, "--BackupCollector.User=user2", "--BackupCollector.Password=pass2")
	os.Args = append(os.Args, "-i=5s", "-r=6s", "-a=7")
	os.Args = append(os.Args, "-d=Mfvs", "-m=100s")
	os.Args = append(os.Args, "-s=100")
	os.Args = append(os.Args, "--Debug")
	checkAll(s, true)
}

func (s *ConfigurationTestSuite) TestOverwriteParameters() {
	s.file.WriteString("primaryCollector: " + LineBreak)
	s.file.WriteString("  fqdn: 127.0.0.0" + LineBreak)
	s.file.WriteString("  port: 30000" + LineBreak)
	s.file.WriteString("  topic: mytopic" + LineBreak)
	s.file.WriteString("  user: user" + LineBreak)
	s.file.WriteString("  password: pass" + LineBreak)
	os.Setenv("VES_PRIMARYCOLLECTOR_FQDN", "127.0.0.1")
	os.Setenv("VES_PRIMARYCOLLECTOR_PORT", "30001")
	os.Args = append(os.Args, "-p=30002")

	var conf VESAgentConfiguration
	err := InitConf(&conf)
	s.NoError(err)
	s.Equal("127.0.0.1", conf.PrimaryCollector.FQDN)
	s.Equal(30002, conf.PrimaryCollector.Port)
}

func checkAll(s *ConfigurationTestSuite, cli bool) {
	var conf VESAgentConfiguration
	err := InitConf(&conf)
	s.NoError(err)
	s.Equal("api", conf.PrimaryCollector.ServerRoot)
	s.Equal("127.0.0.1", conf.PrimaryCollector.FQDN)
	s.Equal(30000, conf.PrimaryCollector.Port)
	s.Equal("mytopic", conf.PrimaryCollector.Topic)
	s.Equal("user", conf.PrimaryCollector.User)
	s.Equal("pass", conf.PrimaryCollector.Password)
	s.Equal("", conf.BackupCollector.ServerRoot)
	s.Equal("127.0.0.2", conf.BackupCollector.FQDN)
	s.Equal(40000, conf.BackupCollector.Port)
	s.Equal("mytopic2", conf.BackupCollector.Topic)
	s.Equal("user2", conf.BackupCollector.User)
	s.Equal("pass2", conf.BackupCollector.Password)
	s.Equal(5*time.Second, conf.Heartbeat.DefaultInterval)
	s.Equal("Mfvs", conf.Measurement.DomainAbbreviation)
	s.Equal(100*time.Second, conf.Measurement.DefaultInterval)
	s.Equal(100, conf.Event.MaxSize)
	s.Equal(6*time.Second, conf.Event.RetryInterval)
	s.Equal(7, conf.Event.MaxMissed)
	s.Equal(true, conf.Debug)
	if !cli {
		s.Equal("localhost:8091", conf.AlertManager.Bind)
	} else {
		s.Equal("localhost:9095", conf.AlertManager.Bind)
	}
}
