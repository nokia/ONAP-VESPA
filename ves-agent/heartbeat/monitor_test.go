package heartbeat

import (
	"testing"
	"time"
	"github.com/nokia/onap-vespa/govel"

	"github.com/stretchr/testify/suite"
)

type MonitorTestSuite struct {
	suite.Suite
	namingCodes map[string]string
}

func (suite *MonitorTestSuite) SetupSuite() {
	suite.namingCodes = make(map[string]string)
	suite.namingCodes["MyVNF"] = "oam"
	suite.namingCodes["ope-2"] = "pro"
}

func TestMonitor(t *testing.T) {
	suite.Run(t, new(MonitorTestSuite))
}

func (suite *MonitorTestSuite) TestNew() {
	mon, err := NewMonitor(&govel.EventConfiguration{VNFName: "MyVNF"}, suite.namingCodes)
	suite.NoError(err)
	suite.NotNil(mon)
	suite.Equal(mon.sourceName, "MyVNF")
}

func (suite *MonitorTestSuite) TestRun() {
	mon, err := NewMonitor(&govel.EventConfiguration{VNFName: "MyVNF", NfNamingCode: "hsxp"}, suite.namingCodes)
	suite.NoError(err)
	suite.NotNil(mon)
	res, err := mon.Run(time.Now(), time.Now(), 5*time.Second)
	hb := res.(*govel.HeartbeatEvent)
	suite.NoError(err)
	suite.Equal(hb.Domain, govel.DomainHeartbeat)
	suite.Equal(hb.Priority, govel.PriorityNormal)
	suite.Equal(hb.Version, float32(3.0))
	suite.Equal(hb.EventID, "heartbeat0000000000")
	suite.Equal(hb.EventName, "heartbeat_hsxp")
	suite.Equal(hb.SourceName, "MyVNF")
	suite.Equal(hb.NfcNamingCode, "oam")
	suite.Equal(hb.NfNamingCode, "hsxp")
	suite.Equal(hb.HeartbeatFieldsVersion, float32(1))
	suite.Equal(hb.HeartbeatInterval, 5)

	res, err = mon.Run(time.Now(), time.Now(), 5*time.Second)
	hb = res.(*govel.HeartbeatEvent)
	suite.NoError(err)
	suite.Equal(hb.EventID, "heartbeat0000000001")
}
