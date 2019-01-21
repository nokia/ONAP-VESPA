package evel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
	"ves-agent/config"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type EvelTestSuite struct {
	suite.Suite
	conf1 *config.CollectorConfiguration
	conf2 *config.CollectorConfiguration
	event *config.EventConfiguration
}

func TestEvel(t *testing.T) {
	suite.Run(t, new(EvelTestSuite))
}

func (s *EvelTestSuite) SetupTest() {
	s.conf1 = &config.CollectorConfiguration{
		FQDN:     "localhost",
		Port:     1234,
		Topic:    "mytopic",
		User:     "myuser",
		Password: "mypassword",
	}
	s.conf2 = &config.CollectorConfiguration{
		FQDN:     "localhost",
		Port:     1234,
		Topic:    "mytopic",
		User:     "myuser",
		Password: "mypassword",
	}
	s.event = &config.EventConfiguration{
		MaxMissed:           1,
		RetryInterval:       time.Second,
		ReportingEntityName: "reportingEntityName",
		ReportingEntityID:   "reportingEntityID",
	}
}

func (s *EvelTestSuite) TestInitialization() {
	evel, err := NewEvel(s.conf1, s.event, "")
	s.NoError(err)
	if !s.NotNil(evel) {
		s.FailNow("Could not initialize evel")
	}
	fmt.Printf("URL: %+v\n", evel)
	s.Equal("http://myuser:mypassword@localhost:1234/api/eventListener/v5", evel.baseURL.String())
	s.Equal("/mytopic", evel.topic)
	s.Equal(0*time.Second, evel.GetMeasurementInterval())
	s.Equal(0*time.Second, evel.GetHeartbeatInterval())
}

func (s *EvelTestSuite) TestProcessCommandList() {
	evel, err := NewEvel(s.conf1, s.event, "")
	s.NoError(err)
	if !s.NotNil(evel) {
		s.FailNow("Could not initialize evel")
	}
	commandList := []Command{
		{CommandType: CommandHeartbeatIntervalChange, HeartbeatInterval: 120},
		{CommandType: CommandMeasurementIntervalChange, MeasurementInterval: 60},
		{CommandType: "foobar"},
	}
	evel.processCommands(commandList)
	s.Equal(60*time.Second, evel.GetMeasurementInterval())
	s.Equal(120*time.Second, evel.GetHeartbeatInterval())
}

func (s *EvelTestSuite) TestPostHeartbeat() {
	type request struct {
		Event HeartbeatEvent `json:"event"`
	}
	var event *request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		event = new(request)
		s.Equal("application/json", req.Header.Get("Content-Type"))
		err := json.NewDecoder(req.Body).Decode(event)
		s.NoError(err)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	s.conf2.FQDN = u.Hostname()
	s.conf2.Port, _ = strconv.Atoi(u.Port())

	evel, err := NewEvel(s.conf2, s.event, "")
	s.NoError(err)

	s.Nil(event)
	hb := NewHeartbeat("id", "name", "mysource", 5)
	err = evel.PostEvent(hb)
	s.NoError(err)
	s.NotNil(event)
	s.Equal(hb, &event.Event)
	s.Equal((&event.Event).ReportingEntityID, "reportingEntityID")
}

func (s *EvelTestSuite) TestPostFault() {
	type request struct {
		Event EventFault `json:"event"`
	}
	var event *request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		event = new(request)
		s.Equal("application/json", req.Header.Get("Content-Type"))
		err := json.NewDecoder(req.Body).Decode(event)
		s.NoError(err)

	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	s.conf2.FQDN = u.Hostname()
	s.conf2.Port, _ = strconv.Atoi(u.Port())

	evel, err := NewEvel(s.conf2, s.event, "")
	s.NoError(err)

	s.Nil(event)
	fault := NewFault("myfault", "myid", "mycondition", "myproblem", PriorityMedium, SeverityMajor, SourceHost, StatusIdle, "mysource")
	err = evel.PostEvent(fault)
	s.NoError(err)
	s.NotNil(event)
	s.Equal(fault, &event.Event)
	s.Equal((&event.Event).ReportingEntityID, "reportingEntityID")
}

func (s *EvelTestSuite) TestPostMeasurements() {
	type request struct {
		Event EventMeasurements `json:"event"`
	}
	var event *request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		event = new(request)
		s.Equal("application/json", req.Header.Get("Content-Type"))
		err := json.NewDecoder(req.Body).Decode(event)
		s.NoError(err)

	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	s.conf2.FQDN = u.Hostname()
	s.conf2.Port, _ = strconv.Atoi(u.Port())

	evel, err := NewEvel(s.conf2, s.event, "")
	s.NoError(err)

	s.Nil(event)
	interval := 10 * time.Second
	now := time.Now()
	meas := NewMeasurements("mymeas", "myid", "source", interval, now, now.Add(interval))
	err = evel.PostEvent(meas)
	s.NoError(err)
	s.NotNil(event)
	s.Equal(meas, &event.Event)
	s.Equal((&event.Event).ReportingEntityID, "reportingEntityID")
}

func (s *EvelTestSuite) TestPostError() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	s.conf2.FQDN = u.Hostname()
	s.conf2.Port, _ = strconv.Atoi(u.Port())

	evel, err := NewEvel(s.conf2, s.event, "")
	s.NoError(err)

	hb := NewHeartbeat("id", "name", "mysource", 5)
	err = evel.PostEvent(hb)
	s.Error(err)
}

func (s *EvelTestSuite) TestNotifyMeasurementIntervalChanged() {
	evel, err := NewEvel(s.conf1, s.event, "")
	s.NoError(err)
	//Check that subscribed channels receive the notification when meas interval changes
	c1 := evel.NotifyMeasurementIntervalChanged(make(chan time.Duration, 1))
	c2 := evel.NotifyMeasurementIntervalChanged(make(chan time.Duration, 1))
	evel.processCommands([]Command{{CommandType: CommandMeasurementIntervalChange, MeasurementInterval: 12}})
	for _, c := range [](<-chan time.Duration){c1, c2} {
		select {
		case v := <-c:
			s.Equal(12*time.Second, v)
		default:
			s.Fail("Interval changed not sent to all channels")
		}
	}

	// Also check that blocked channel will receive nothing and won't cause deadlock
	c3 := evel.NotifyMeasurementIntervalChanged(make(chan time.Duration))
	evel.processCommands([]Command{{CommandType: CommandMeasurementIntervalChange, MeasurementInterval: 14}})
	select {
	case <-c3:
		s.Fail("Channel should be empty")
	default:
	}
}

func (s *EvelTestSuite) TestNotifyHeartbeatIntervalChanged() {
	evel, err := NewEvel(s.conf1, s.event, "")
	s.NoError(err)
	//Check that subscribed channels receive the notification when meas interval changes
	c1 := evel.NotifyHeartbeatIntervalChanged(make(chan time.Duration, 1))
	c2 := evel.NotifyHeartbeatIntervalChanged(make(chan time.Duration, 1))
	evel.processCommands([]Command{{CommandType: CommandHeartbeatIntervalChange, HeartbeatInterval: 12}})
	for _, c := range [](<-chan time.Duration){c1, c2} {
		select {
		case v := <-c:
			s.Equal(12*time.Second, v)
		default:
			s.Fail("Interval changed not sent to all channels")
		}
	}

	// Also check that blocked channel will receive nothing and won't cause deadlock
	c3 := evel.NotifyHeartbeatIntervalChanged(make(chan time.Duration))
	evel.processCommands([]Command{{CommandType: CommandHeartbeatIntervalChange, HeartbeatInterval: 14}})
	select {
	case <-c3:
		s.Fail("Channel should be empty")
	default:
	}
}

func (s *EvelTestSuite) TestPostMeasurementsBatch() {
	type request struct {
		EventList []EventMeasurements `json:"eventList"`
	}
	var event *request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		event = new(request)
		s.Equal("application/json", req.Header.Get("Content-Type"))
		err := json.NewDecoder(req.Body).Decode(event)
		s.NoError(err)

	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	s.conf2.FQDN = u.Hostname()
	s.conf2.Port, _ = strconv.Atoi(u.Port())

	evel, err := NewEvel(s.conf2, s.event, "")
	s.NoError(err)

	s.Nil(event)
	interval := 10 * time.Second
	now := time.Now()
	batch := make(Batch, 2)

	batch[0] = NewMeasurements("mymeas", "myid", "source", interval, now, now.Add(interval))
	batch[1] = NewMeasurements("mymeas2", "myid2", "source", interval, now, now.Add(interval))
	err = evel.PostBatch(batch)
	s.NoError(err)
	s.NotNil(event)
	s.Len(event.EventList, 2)
	s.Equal(batch[0], &event.EventList[0])
	s.Equal(batch[1], &event.EventList[1])
}

func (s *EvelTestSuite) TestPostMeasurementsBatchTooLarge() {
	type request struct {
		EventList []EventMeasurements `json:"eventList"`
	}
	var events []request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s.Equal("application/json", req.Header.Get("Content-Type"))
		var event request
		err := json.NewDecoder(req.Body).Decode(&event)
		s.NoError(err)
		events = append(events, event)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	s.conf2.FQDN = u.Hostname()
	s.conf2.Port, _ = strconv.Atoi(u.Port())
	s.event.MaxSize = 500

	evel, err := NewEvel(s.conf2, s.event, "")
	s.NoError(err)

	s.Nil(events)
	interval := 10 * time.Second
	now := time.Now()
	batch := make(Batch, 2)

	batch[0] = NewMeasurements("mymeas", "myid", "source", interval, now, now.Add(interval))
	batch[1] = NewMeasurements("mymeas2", "myid2", "source", interval, now, now.Add(interval))
	err = evel.PostBatch(batch)
	s.NoError(err)
	s.NotNil(events)
	s.Len(events, 2)
	s.Len(events[0].EventList, 1)
	s.Equal(batch[0], &events[0].EventList[0])
	s.Equal(batch[1], &events[1].EventList[0])
	s.event.MaxSize = 100
}

func (s *EvelTestSuite) TestPostEmptyBatch() {
	type request struct {
		EventList []Event `json:"eventList"`
	}
	var events []request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s.Equal("application/json", req.Header.Get("Content-Type"))
		var event request
		err := json.NewDecoder(req.Body).Decode(&event)
		s.NoError(err)
		events = append(events, event)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	s.conf2.FQDN = u.Hostname()
	s.conf2.Port, _ = strconv.Atoi(u.Port())
	s.event.MaxSize = 500
	evel, err := NewEvel(s.conf2, s.event, "")
	s.NoError(err)
	s.Nil(events)

	err = evel.PostBatch(Batch{})
	s.NoError(err)
	s.Empty(events)
	s.event.MaxSize = 100
}

func (s *EvelTestSuite) TestCannotSplitBatch() {
	s.event.MaxSize = 12
	evel, err := NewEvel(s.conf1, s.event, "")
	s.NoError(err)
	s.Panics(assert.PanicTestFunc(func() {
		evel.PostBatch(Batch{NewHeartbeat("id", "name", "foo", 1234)})
	}))
	s.event.MaxSize = 100
}
