package agent

import (
	"encoding/json"
	"testing"
	"time"
	"github.com/nokia/onap-vespa/ves-agent/config"
	"github.com/nokia/onap-vespa/govel"
	"github.com/nokia/onap-vespa/ves-agent/ha"
	"github.com/nokia/onap-vespa/ves-agent/rest"

	"github.com/prometheus/alertmanager/template"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

//Cluster Mock

type ClusterMock struct {
	mock.Mock
}

func (ves *ClusterMock) PostEvent(evt govel.Event) error {
	args := ves.MethodCalled("PostEvent", evt)
	return args.Error(0)
}

func (ves *ClusterMock) PostBatch(batch govel.Batch) error {
	args := ves.MethodCalled("PostBatch", batch)
	return args.Error(0)
}

func (ves *ClusterMock) GetMeasurementInterval() time.Duration {
	args := ves.MethodCalled("GetMeasurementInterval")
	return args.Get(0).(time.Duration)
}
func (ves *ClusterMock) GetHeartbeatInterval() time.Duration {
	args := ves.MethodCalled("GetHeartbeatInterval")
	return args.Get(0).(time.Duration)
}
func (ves *ClusterMock) NotifyMeasurementIntervalChanged(ch chan time.Duration) <-chan time.Duration {
	args := ves.MethodCalled("NotifyMeasurementIntervalChanged", ch)
	return args.Get(0).(chan time.Duration)
}
func (ves *ClusterMock) NotifyHeartbeatIntervalChanged(ch chan time.Duration) <-chan time.Duration {
	args := ves.MethodCalled("NotifyHeartbeatIntervalChanged", ch)
	return args.Get(0).(chan time.Duration)
}

//Tests

type AgentTestSuite struct {
	suite.Suite
	hbConf      *config.HeartbeatConfiguration
	eventConf   *govel.EventConfiguration
	vesConf     *config.VESAgentConfiguration
	cluster     ClusterMock
	namingCodes map[string]string
}

func TestAgent(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}

func (suite *AgentTestSuite) SetupSuite() {
	hbConf := config.HeartbeatConfiguration{
		DefaultInterval: 1 * time.Second,
	}
	oam := govel.NfcNamingCode{
		Type:  "oam",
		Vnfcs: []string{"dpa2bhsxp5001vm001oam001", "ope-1", "ope-2"},
	}
	etl := govel.NfcNamingCode{
		Type:  "etl",
		Vnfcs: []string{"pro-0", "pro-1"},
	}
	eventConf := govel.EventConfiguration{
		MaxSize:        200,
		NfNamingCode:   "hspx",
		NfcNamingCodes: []govel.NfcNamingCode{oam, etl},
	}
	measConf := config.MeasurementConfiguration{
		DefaultInterval: 2 * time.Second,
	}
	collectorConf := govel.CollectorConfiguration{
		FQDN:     "localhost",
		Port:     1234,
		Topic:    "mytopic",
		User:     "myuser",
		Password: "mypassword",
	}
	backupConf := govel.CollectorConfiguration{
		FQDN:     "",
		Port:     0,
		Topic:    "",
		User:     "",
		Password: "",
	}
	suite.hbConf = &hbConf
	suite.eventConf = &eventConf
	suite.vesConf = &config.VESAgentConfiguration{
		Event:            eventConf,
		Heartbeat:        hbConf,
		Measurement:      measConf,
		PrimaryCollector: collectorConf,
		BackupCollector:  backupConf,
	}
	suite.cluster = ClusterMock{}

	suite.namingCodes = make(map[string]string)
	suite.namingCodes["dpa2bhsxp5001vm001oam001"] = "oam"
	suite.namingCodes["ope-1"] = "oam"
	suite.namingCodes["ope-2"] = "oam"
	suite.namingCodes["pro-0"] = "etl"
	suite.namingCodes["pro-1"] = "etl"
}

func (suite *AgentTestSuite) TestInitMeasScheduler() {
	state := ha.NewInMemState()
	//Without required interval
	measSched := initMeasScheduler(suite.vesConf, suite.namingCodes, state)
	suite.Equal(measSched.GetInterval(), 2*time.Second)
	//Withrequired interval
	state.UpdateInterval("measurements", 5*time.Second)
	measSched = initMeasScheduler(suite.vesConf, suite.namingCodes, state)
	suite.Equal(measSched.GetInterval(), 5*time.Second)
}

func (suite *AgentTestSuite) TestInitHbScheduler() {
	state := ha.NewInMemState()
	//Without required interval
	hbSched := initHbScheduler(suite.eventConf, 1*time.Second, suite.namingCodes, state)
	suite.Equal(hbSched.GetInterval(), 1*time.Second)
	//Withrequired interval
	state.UpdateInterval("heartbeats", 5*time.Second)
	hbSched = initHbScheduler(suite.eventConf, 1*time.Second, suite.namingCodes, state)
	suite.Equal(hbSched.GetInterval(), 5*time.Second)
}

func (suite *AgentTestSuite) TestNewAgent() {
	agent := NewAgent(suite.vesConf)
	suite.NotNil(agent)
	suite.Equal(agent.measSched.GetInterval(), 2*time.Second)
	suite.Equal(agent.hbSched.GetInterval(), 1*time.Second)
	suite.NotNil(agent.fm)
	suite.Equal(agent.alertRoute.Pattern, suite.vesConf.AlertManager.Path)
	suite.cluster.AssertExpectations(suite.T())
}

func (suite *AgentTestSuite) TestInitNfcNamingCode() {
	res := initNfcNamingCode(suite.eventConf.NfcNamingCodes)
	suite.Equal(res, suite.namingCodes)
}

func (suite *AgentTestSuite) TestAgent() {
	var alertData = []byte(`
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
	var alert template.Alert

	agent := NewAgent(suite.vesConf)
	suite.NotNil(agent)
	<-agent.state.LeaderCh()

	hbIntCh := make(chan time.Duration, 1024)
	measIntCh := make(chan time.Duration, 1024)
	suite.cluster.On("NotifyMeasurementIntervalChanged", mock.Anything).Once().Return(measIntCh)
	suite.cluster.On("NotifyHeartbeatIntervalChanged", mock.Anything).Once().Return(hbIntCh)
	agent.listen("localhost:0", &suite.cluster)
	suite.Nil(agent.measTimer)
	suite.Nil(agent.hbTimer)
	suite.NotNil(agent.measIntervalCh)
	suite.NotNil(agent.hbIntervalCh)
	suite.NotNil(agent.alertCh)

	suite.cluster.On("PostEvent", mock.AnythingOfType("*govel.EventFault")).Once().Return(nil)
	// PostEvent with heartbeat MUST be called at least once
	suite.cluster.On("PostEvent", mock.AnythingOfType("*govel.HeartbeatEvent")).Once().Return(nil)
	// But it sometimes may be called a second time, depending on timing.
	suite.cluster.On("PostEvent", mock.AnythingOfType("*govel.HeartbeatEvent")).Maybe().Return(nil)
	suite.cluster.On("PostBatch", mock.Anything).Once().Return(nil)
	err := json.Unmarshal(alertData, &alert)
	if err != nil {
		suite.Fail("Error in unmarshall function for alert")
	}
	agent.alertCh <- rest.MessageFault{Alert: alert, Response: make(chan error)}
	agent.measTimer = agent.measSched.WaitChan()
	agent.hbTimer = agent.hbSched.WaitChan()
	agent.leaderStep(&suite.cluster)
	agent.leaderStep(&suite.cluster)
	agent.leaderStep(&suite.cluster)
	suite.cluster.AssertExpectations(suite.T())
}

func (suite *AgentTestSuite) TestStats() {
	agent := NewAgent(suite.vesConf)
	suite.NotNil(agent)
	<-agent.state.LeaderCh()
	stats := agent.Stats()
	suite.Equal("Leader", stats["raft"].(map[string]string)["state"])
}
