package heartbeat

import (
	"fmt"
	"time"
	"github.com/nokia/onap-vespa/govel"
)

// MonitorState handles the monitor internal state
type MonitorState interface {
	// NextHeartbeatIndex return the next event index and increments it
	NextHeartbeatIndex() (int64, error)
}

type inMemState struct {
	index int64
}

func (mem *inMemState) NextHeartbeatIndex() (int64, error) {
	i := mem.index
	mem.index++
	return i, nil
}

// Monitor is an utility to create heartbeats
type Monitor struct {
	sourceName   string //sourcename for building the heartbeat
	nfNamingCode string
	state        MonitorState      // Monitor internal state
	namingCodes  map[string]string // Cache for VnfcNamingCode from VnfcName
}

// NewMonitorWithState creates a new Heartbeat Monitor from provided configuration
// and provided state handler
func NewMonitorWithState(conf *govel.EventConfiguration, namingCodes map[string]string, state MonitorState) (*Monitor, error) {
	return &Monitor{sourceName: conf.VNFName, nfNamingCode: conf.NfNamingCode, state: state, namingCodes: namingCodes}, nil
}

// NewMonitor creates a new Heartbeat Monitor from provided configuration
// that use an in memory state
func NewMonitor(conf *govel.EventConfiguration, namingCodes map[string]string) (*Monitor, error) {
	return NewMonitorWithState(conf, namingCodes, &inMemState{index: 0})
}

// Run creates a new Heartbeat
func (mon *Monitor) Run(from, to time.Time, interval time.Duration) (interface{}, error) {
	idx, err := mon.state.NextHeartbeatIndex()
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("heartbeat%.10d", idx)
	eventName := "heartbeat_" + mon.nfNamingCode
	hb := govel.NewHeartbeat(id, eventName, mon.sourceName, int(interval.Seconds()))
	hb.NfNamingCode = mon.nfNamingCode
	hb.NfcNamingCode = mon.namingCodes[mon.sourceName]
	return hb, nil
}
