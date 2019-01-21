package rest

import (
	"net/http/httptest"
	"strings"
	"testing"
	"ves-agent/config"
	"ves-agent/convert"

	"github.com/stretchr/testify/suite"
)

var alertRoute = Route{
	"AlertReceiver",
	"POST",
	"/alerts",
	nil,
}

var postdata1 = []byte(`
{
	"receiver": "calm-webhook",
	"status": "firing",
	"alerts": [
		{
			"status": "firing",
			"labels": {
				"id": "201",
				"system_name": "TEST",
				"alertname": "AlertNodeFailure1",
				"severity": "critical",
				"VNFC": "hspx5001vm001",
				"calm": "true",
				"event_type": "x2"
			},
			"annotations": {
				"service": "NodeSupervision",
				"summary": "Node pilot-0 down",
				"description": "Node pilot-0 is seen disconnected from cluster"
			}
		}
	]		
}
`)

var postdata2 = []byte(`
{
	"receiver": "calm-webhook",
	"status": "firing",
	"alerts": [
		{
			"status": "firing",
			"labels": {
				"id": "201",
				"system_name": "TEST1",
				"alertname": "AlertNodeFailure21",
				"severity": "critical",
				"VNFC": "hspx5001vm002",
				"calm": "true",
				"event_type": "x2"
			},
			"annotations": {
				"service": "NodeSupervision",
				"summary": "Node pilot-0 down",
				"description": "Node pilot-0 is seen disconnected from cluster"
			}
		},
        {
            "status": "firing",
            "labels": {
                "id": "201",
                "system_name": "TEST2",
                "alertname": "AlertNodeFailure22",
                "severity": "critical",
                "VNFC": "hspx5001vm003",
                "calm": "true",
                "event_type": "x2"
            },
            "annotations": {
				"service": "NodeSupervision",
                "summary": "Node pilot-1 down",
                "description": "Node pilot-1 is seen disconnected from cluster"
            }
        }
	]		
}
`)

var postinvaliddata = []byte(`
{
	"receiver": "calm-webhook",
	"status": "firing",
	"alerts": [
		{
			"status": "firing",
			"labels": {
				"id": "201",
				"system_name": "TEST",
				"alertname": "AlertNodeFailure3",
				"severity": "critical",
				"VNFC": "hspx5001vm001",
				"calm": "true",
				"event_type": "x2"
			},
			"annotations": {
				"service": "NodeSupervision",
				"summary": "Node pilot-0 down",
				"description": "Node pilot-0 is seen disconnected from cluster"
		}
	]		
}
`)

type HandlerTestSuite struct {
	suite.Suite
	confEvent config.EventConfiguration
	fm        *convert.FaultManager
	faultID   int32
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupSuite() {
	suite.confEvent = config.EventConfiguration{
		MaxSize:      200,
		NfNamingCode: "hspx",
	}
	suite.fm = convert.NewFaultManager(&suite.confEvent)
	suite.faultID = 0
}

func (suite *HandlerTestSuite) TestHandlerData1Ok() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/alerts", strings.NewReader(string(postdata1)))
	req.Header.Set("content-Type", "application/json")

	alertCh := make(chan MessageFault)

	alertRoute.HandlerFunc = AlertReceiver(alertCh)
	// create an unstarted new server to receive http POST from prometheus
	alertHandler := NewServer([]Route{alertRoute})
	go alertHandler.ServeHTTP(resp, req)
	messageFault := <-alertCh
	//suite.Equal("Fault_hspx_AlertNodeFailure", messageFault.Fault.EventName)
	suite.Equal("AlertNodeFailure1", messageFault.Alert.Labels["alertname"])
	suite.Equal(200, resp.Code, "Bad HTTP response status code")

	close(messageFault.Response)
}

func (suite *HandlerTestSuite) TestHandlerData2Ok() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/alerts", strings.NewReader(string(postdata2)))
	req.Header.Set("content-Type", "application/json")

	alertCh := make(chan MessageFault)

	alertRoute.HandlerFunc = AlertReceiver(alertCh)
	alertHandler := NewServer([]Route{alertRoute})
	go alertHandler.ServeHTTP(resp, req)
	messageFault := <-alertCh
	//suite.Equal("Fault_hspx_AlertNodeFailure", messageFault.Fault.EventName)
	suite.Equal("AlertNodeFailure21", messageFault.Alert.Labels["alertname"])
	suite.Equal(200, resp.Code, "Bad HTTP response status code")
	close(messageFault.Response)
	messageFault = <-alertCh
	//suite.Equal("Fault_hspx_AlertNodeFailure", messageFault.Fault.EventName)
	suite.Equal("AlertNodeFailure22", messageFault.Alert.Labels["alertname"])
	suite.Equal(200, resp.Code, "Bad HTTP response status code")
	close(messageFault.Response)

	// Also check that blocked channel will receive nothing and won't cause deadlock
	select {
	case <-alertCh:
		suite.Fail("Channel should be empty")
	default:
	}
}

func (suite *HandlerTestSuite) TestHandlerData1InvalidContent() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/alerts", strings.NewReader(string(postdata1)))
	//req.Header.Set("content-Type", "application/json")

	alertCh := make(chan MessageFault)

	alertRoute.HandlerFunc = AlertReceiver(alertCh)
	alertHandler := NewServer([]Route{alertRoute})
	alertHandler.ServeHTTP(resp, req)

	suite.Equal(500, resp.Code, "Bad HTTP response status code")
	//check invalide fault is not sent to channel
	select {
	case <-alertCh:
		suite.Fail("Channel should be empty")
	default:
	}
}

func (suite *HandlerTestSuite) TestHandlerData1InvalidData() {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/alerts", strings.NewReader(string(postinvaliddata)))
	req.Header.Set("content-Type", "application/json")

	alertCh := make(chan MessageFault)

	alertRoute.HandlerFunc = AlertReceiver(alertCh)
	alertHandler := NewServer([]Route{alertRoute})
	alertHandler.ServeHTTP(resp, req)

	suite.Equal(400, resp.Code, "Bad HTTP response status code")
	//check invalide fault is not sent to channel
	select {
	case <-alertCh:
		suite.Fail("Channel should be empty")
	default:
	}
}

func (suite *HandlerTestSuite) TestHandlerDataFaultAlredyExist() {
	resp := httptest.NewRecorder()
	req1 := httptest.NewRequest("POST", "/alerts", strings.NewReader(string(postdata1)))
	req1.Header.Set("content-Type", "application/json")

	alertCh := make(chan MessageFault)

	alertRoute.HandlerFunc = AlertReceiver(alertCh)
	alertHandler := NewServer([]Route{alertRoute})
	go alertHandler.ServeHTTP(resp, req1)
	messageFault := <-alertCh
	close(messageFault.Response)

	suite.Equal(200, resp.Code, "Bad HTTP response status code")
	select {
	case messageFault := <-alertCh:
		//suite.Equal("Fault_hspx_AlertNodeFailure", messageFault.Fault.EventName)
		suite.Equal("AlertNodeFailure1", messageFault.Alert.Labels["alertname"])
		suite.Equal(200, resp.Code, "Bad HTTP response status code")
	default:
	}
}
