package govel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ClusterTestSuite struct {
	suite.Suite
	conf1     *CollectorConfiguration
	conf2     *CollectorConfiguration
	confEmpty *CollectorConfiguration
	event     *EventConfiguration
}

func TestCluster(t *testing.T) {
	suite.Run(t, new(ClusterTestSuite))
}

func (s *ClusterTestSuite) SetupTest() {
	s.conf1 = &CollectorConfiguration{
		FQDN:     "localhost", 
		Port:     1234,
		Topic:    "mytopic",
		User:     "myuser",
		Password: "mypassword",
	}
	s.conf2 = &CollectorConfiguration{
		FQDN:     "localhost2",
		Port:     5678,
		Topic:    "mytopic2",
		User:     "myuser2",
		Password: "mypassword2",
	}
	s.confEmpty = &CollectorConfiguration{
		FQDN:     "",
		Port:     0,
		Topic:    "",
		User:     "",
		Password: "",
	}
	s.event = &EventConfiguration{
		MaxMissed:     1,
		RetryInterval: time.Second,
	}
}
func (s *ClusterTestSuite) TestInitialization() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	s.Equal(cluster.primaryVES, cluster.activVES)
	s.Equal("http://myuser:mypassword@localhost:1234/api/eventListener/v5", cluster.primaryVES.baseURL.String())
	s.Equal("http://myuser2:mypassword2@localhost2:5678/api/eventListener/v5", cluster.backupVES.baseURL.String())
	s.Equal(1, cluster.maxMissed)
	s.Equal(time.Second, cluster.retryInterval)
}

func (s *ClusterTestSuite) TestInitializationHttps() {
	const cacert = `-----BEGIN CERTIFICATE-----
MIIC4jCCAcqgAwIBAgIQCunY77fnqyu57rOU31qeRzANBgkqhkiG9w0BAQsFADAQ
MQ4wDAYDVQQKEwVOb2tpYTAeFw0xODEwMDgxNTI1MTVaFw0yODEwMDgxNTI1MTVa
MBAxDjAMBgNVBAoTBU5va2lhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAu/U7NUjz3jhPpj3nEJ9DOLn/3cGO1bqKxLaCPQ5OliKiKbUlOdEPUlW5zjV1
X7nVMc51hLg1vsweWs5Uxc7NbasE+8SY1jN/jYZS8TZxQr7uKGRItERK9b2dbLiB
UYZdnaYxTdxiYGetDGn2Fm4PQ0Ei/lqM7TSJ1/L93QWrHwDJLg6zHplvr6rLh0jJ
pmF4OZk0VtZmI4w6ekjY3dyecfNzFAubiQBAZ84UK/9TgjaRtcnjciH+5fTKNGoq
hIDwuSJqU8GprhWx3M+ovdAgoTmIL8sddIT8FnvGHM0cecvm5GBTs0qO5fBmhbuq
slfzeWpKpgWgh+l6T/VNSig/VQIDAQABozgwNjAOBgNVHQ8BAf8EBAMCAqQwEwYD
VR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsF
AAOCAQEAPv/idweSFuimFc+GHlunylTsnX5pdSIDtlIRqU77pyi4icDqiIInJlAx
r31VBn5kepeFZELB/BKl5ulTWXVDCt9u0Tw+VgrIO0+sCUMG37pTHyANdVNVmTZ+
zZHHwxf1sPBVSx+pKQAPc676TrH8PByW6cO2juCwUYKLjMtfBfT60pCoNbUSi/3V
6XQMQm0JhsgXxH6kbPWpQ8wG0aFCa8uMF6tt9a2UloWdoCKvV2IGxV4hferiD6un
LAaDIw5mTQ6JrQc8OYazF7j0LZLro/8BC0A+24NErwB/KwxfQvMOGmHrT7gstNOu
unwD+TffUa2jWGiKGjohv2u18i+Gyw==
-----END CERTIFICATE-----`

	cluster, err := NewCluster(s.conf1, s.conf2, s.event, cacert)
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	s.Equal(cluster.primaryVES, cluster.activVES)
	s.Equal("https://myuser:mypassword@localhost:1234/api/eventListener/v5", cluster.primaryVES.baseURL.String())
	s.Equal("https://myuser2:mypassword2@localhost2:5678/api/eventListener/v5", cluster.backupVES.baseURL.String())
	s.Equal(1, cluster.maxMissed)
	s.Equal(time.Second, cluster.retryInterval)
}

func (s *ClusterTestSuite) TestInitializationEmptyBackup() {
	cluster, err := NewCluster(s.conf1, s.confEmpty, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	s.Nil(cluster.backupVES)
	s.Equal(cluster.primaryVES, cluster.activVES)
	s.Equal("http://myuser:mypassword@localhost:1234/api/eventListener/v5", cluster.primaryVES.baseURL.String())
	s.Equal(1, cluster.maxMissed)
	s.Equal(time.Second, cluster.retryInterval)
}

func (s *ClusterTestSuite) TestIsPrimaryActive() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	s.True(cluster.isPrimaryActive())
	cluster.activVES = cluster.backupVES
	s.False(cluster.isPrimaryActive())
}

func (s *ClusterTestSuite) TestGetMeasurementInterval() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	commandList := []Command{
		{CommandType: CommandMeasurementIntervalChange, MeasurementInterval: 1},
	}
	cluster.activVES.processCommands(commandList)
	s.Equal(time.Second, cluster.GetMeasurementInterval())
	cluster.switchCollector()
	s.Equal(0*time.Second, cluster.GetMeasurementInterval())
}

func (s *ClusterTestSuite) TestGetHeartbeatInterval() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	commandList := []Command{
		{CommandType: CommandHeartbeatIntervalChange, HeartbeatInterval: 1},
	}
	cluster.activVES.processCommands(commandList)
	s.Equal(time.Second, cluster.GetHeartbeatInterval())
	cluster.switchCollector()
	s.Equal(0*time.Second, cluster.GetHeartbeatInterval())
}

func (s *ClusterTestSuite) TestNotifyMeasurementIntervalChanged() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	//Check that subscribed channels receive the notification when meas interval changes
	c1 := cluster.NotifyMeasurementIntervalChanged(make(chan time.Duration, 1))
	c2 := cluster.NotifyMeasurementIntervalChanged(make(chan time.Duration, 1))
	cluster.primaryVES.processCommands([]Command{Command{CommandType: CommandMeasurementIntervalChange, MeasurementInterval: 12}})
	cluster.backupVES.processCommands([]Command{Command{CommandType: CommandMeasurementIntervalChange, MeasurementInterval: 12}})
	for _, c := range [](<-chan time.Duration){c1, c2} {
		select {
		case v := <-c:
			s.Equal(12*time.Second, v)
		default:
			s.Fail("Interval changed not sent to all channels")
		}
	}

	// Also check that blocked channel will receive nothing and won't cause deadlock
	c3 := cluster.NotifyMeasurementIntervalChanged(make(chan time.Duration))
	cluster.activVES.processCommands([]Command{Command{CommandType: CommandMeasurementIntervalChange, MeasurementInterval: 14}})
	select {
	case <-c3:
		s.Fail("Channel should be empty")
	default:
	}
}

func (s *ClusterTestSuite) TestNotifyHeartbeatIntervalChanged() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	//Check that subscribed channels receive the notification when meas interval changes
	c1 := cluster.NotifyHeartbeatIntervalChanged(make(chan time.Duration, 1))
	c2 := cluster.NotifyHeartbeatIntervalChanged(make(chan time.Duration, 1))
	cluster.primaryVES.processCommands([]Command{Command{CommandType: CommandHeartbeatIntervalChange, HeartbeatInterval: 12}})
	cluster.backupVES.processCommands([]Command{Command{CommandType: CommandHeartbeatIntervalChange, HeartbeatInterval: 12}})
	for _, c := range [](<-chan time.Duration){c1, c2} {
		select {
		case v := <-c:
			s.Equal(12*time.Second, v)
		default:
			s.Fail("Interval changed not sent to all channels")
		}
	}

	// Also check that blocked channel will receive nothing and won't cause deadlock
	c3 := cluster.NotifyHeartbeatIntervalChanged(make(chan time.Duration))
	cluster.activVES.processCommands([]Command{Command{CommandType: CommandHeartbeatIntervalChange, HeartbeatInterval: 14}})
	select {
	case <-c3:
		s.Fail("Channel should be empty")
	default:
	}
}

func (s *ClusterTestSuite) TestPostEvent() {
	var event *HeartbeatEvent
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		event = new(HeartbeatEvent)
		s.Equal("application/json", req.Header.Get("Content-Type"))
		err := json.NewDecoder(req.Body).Decode(event)
		s.NoError(err)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	port, _ := strconv.Atoi(u.Port())

	s.conf1.FQDN = u.Hostname()
	s.conf1.Port = port
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	s.Nil(event)
	hb := NewHeartbeat("id", "name", "mysource", 5)
	err = cluster.PostEvent(hb)
	s.NoError(err)
	s.Equal(cluster.activVES, cluster.primaryVES)
}
func (s *ClusterTestSuite) TestPostEventSwitch() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	hb := NewHeartbeat("id", "name", "mysource", 5)
	err = cluster.PostEvent(hb)
	s.Error(err)
	s.Equal(cluster.activVES, cluster.backupVES)
}

func (s *ClusterTestSuite) TestSwitch() {
	cluster, err := NewCluster(s.conf1, s.conf2, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	cluster.switchCollector()
	s.Equal(cluster.activVES, cluster.backupVES)
	s.Equal("http://myuser:mypassword@localhost:1234/api/eventListener/v5", cluster.primaryVES.baseURL.String())
	s.Equal("http://myuser2:mypassword2@localhost2:5678/api/eventListener/v5", cluster.backupVES.baseURL.String())
	s.Equal(1, cluster.maxMissed)
	s.Equal(time.Second, cluster.retryInterval)
	cluster.switchCollector()
	s.Equal(cluster.primaryVES, cluster.activVES)
	s.Equal("http://myuser:mypassword@localhost:1234/api/eventListener/v5", cluster.primaryVES.baseURL.String())
	s.Equal("http://myuser2:mypassword2@localhost2:5678/api/eventListener/v5", cluster.backupVES.baseURL.String())
}

func (s *ClusterTestSuite) TestSwitchEmptyBackup() {
	cluster, err := NewCluster(s.conf1, s.confEmpty, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	cluster.switchCollector()
	s.Equal(cluster.primaryVES, cluster.activVES)
	s.Equal("http://myuser:mypassword@localhost:1234/api/eventListener/v5", cluster.primaryVES.baseURL.String())
	s.Equal(1, cluster.maxMissed)
	s.Equal(time.Second, cluster.retryInterval)
}

func (s *ClusterTestSuite) TestSwitchEmptyPrimary() {
	cluster, err := NewCluster(s.confEmpty, s.conf1, s.event, "")
	s.NoError(err)
	if !s.NotNil(cluster) {
		s.FailNow("Could not initialize evel")
	}
	cluster.primaryVES = nil
	cluster.switchCollector()
	s.Equal(cluster.backupVES, cluster.activVES)
	s.Equal("http://myuser:mypassword@localhost:1234/api/eventListener/v5", cluster.backupVES.baseURL.String())
	s.Equal(1, cluster.maxMissed)
	s.Equal(time.Second, cluster.retryInterval)
}
