package ha

import (
	"bytes"
	"os"
	"testing"
	"time"
	"github.com/nokia/onap-vespa/ves-agent/config"

	"github.com/stretchr/testify/suite"
)

func fillState(state AgentState, now time.Time) error {
	for i := 0; i < 10; i++ {
		if _, err := state.NextHeartbeatIndex(); err != nil {
			return err
		}
	}
	for i := 0; i < 15; i++ {
		if _, err := state.NextMeasurementIndex(); err != nil {
			return err
		}
	}
	if err := state.UpdateScheduler("foobar", 170*time.Minute, now); err != nil {
		return err
	}
	faultIdx, err := state.NextFaultIndex()
	if err != nil {
		return err
	}
	if err = state.StoreFaultInStorage("MyFault", faultIdx); err != nil {
		return err
	}
	if err = state.InitAlertInfos(faultIdx); err != nil {
		return err
	}
	if err = state.IncrementFaultSn(faultIdx); err != nil {
		return err
	}
	return state.SetFaultStartEpoch(faultIdx, 123456)
}

type SnapshotTestSuite struct {
	suite.Suite
}

func TestSnaphotAndRestore(t *testing.T) {
	suite.Run(t, new(SnapshotTestSuite))
}

func (s *SnapshotTestSuite) TestSnapshotAndRestore() {
	state := NewInMemState()
	now := time.Now().Truncate(0)
	s.NoError(fillState(state, now))

	snap := state.Snapshot()
	s.EqualValues(10, snap.HbIdx)
	s.EqualValues(15, snap.MeasIdx)
	s.Len(snap.Schedulers, 1)
	sched, ok := snap.Schedulers["foobar"]
	s.True(ok)
	s.Equal(now.UTC(), sched.Next)
	s.Equal(170*time.Minute, sched.Interval)

	newState := NewInMemState()
	newState.Restore(snap)
	s.EqualValues(state, newState)
}

func (s *SnapshotTestSuite) TestFsmSnapshotAndRetsore() {
	cluster, _ := NewCluster("./test_datadir", &config.ClusterConfiguration{Debug: true, DisplayLogs: true}, NewInMemState())
	defer func() {
		cluster.Shutdown()
		os.RemoveAll("./test_datadir")
	}()
	for !<-cluster.LeaderCh() {
	}
	now := time.Now().Truncate(0)
	s.NoError(fillState(cluster, now))

	snap, err := cluster.fsm.Snapshot()
	s.NoError(err)
	sink := memorySnapshotSink{}

	s.NoError(snap.Persist(&sink))

	newFSM := NewFSM(NewInMemState(), true)
	s.NoError(newFSM.Restore(&sink))
	s.EqualValues(cluster.fsm.state.Snapshot(), newFSM.state.Snapshot())
	s.EqualValues(cluster.fsm.state, newFSM.state)
}

type memorySnapshotSink struct {
	buf bytes.Buffer
}

func (m *memorySnapshotSink) ID() string { return "" }

func (m *memorySnapshotSink) Cancel() error { return nil }

func (m *memorySnapshotSink) Close() error { return nil }

func (m *memorySnapshotSink) Write(p []byte) (n int, err error) {
	return m.buf.Write(p)
}

func (m *memorySnapshotSink) Read(p []byte) (n int, err error) {
	return m.buf.Read(p)
}
