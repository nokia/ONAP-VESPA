package ha

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type StateTestSuite struct {
	suite.Suite
	state    AgentState
	setup    func(s *StateTestSuite)
	teardown func(s *StateTestSuite)
}

func TestInMemState(t *testing.T) {
	suite.Run(t, &StateTestSuite{
		setup: func(s *StateTestSuite) {
			s.state = NewInMemState()
		},
	})
}

func (s *StateTestSuite) SetupTest() {
	if s.setup != nil {
		s.setup(s)
	}
}

func (s *StateTestSuite) TearDownTest() {
	if s.teardown != nil {
		s.teardown(s)
	}
}

func (s *StateTestSuite) TestMeasurementIndex() {
	meas, err := s.state.NextMeasurementIndex()
	s.NoError(err)
	s.EqualValues(0, meas)
	meas, err = s.state.NextMeasurementIndex()
	s.NoError(err)
	s.EqualValues(1, meas)
	meas, err = s.state.NextMeasurementIndex()
	s.NoError(err)
	s.EqualValues(2, meas)
}

func (s *StateTestSuite) TestHeartbeatIndex() {
	meas, err := s.state.NextHeartbeatIndex()
	s.NoError(err)
	s.EqualValues(0, meas)
	meas, err = s.state.NextHeartbeatIndex()
	s.NoError(err)
	s.EqualValues(1, meas)
	meas, err = s.state.NextHeartbeatIndex()
	s.NoError(err)
	s.EqualValues(2, meas)
}

func (s *StateTestSuite) TestScheduler() {
	schn := "test"
	now := time.Now()
	s.EqualValues(0, s.state.Interval(schn))
	s.Equal(time.Time{}, s.state.NextRun(schn))

	s.NoError(s.state.UpdateInterval(schn, 10*time.Second))
	s.Equal(10*time.Second, s.state.Interval(schn))
	s.Equal(time.Time{}, s.state.NextRun(schn))

	s.NoError(s.state.UpdateNextRun(schn, now))
	s.Equal(10*time.Second, s.state.Interval(schn))
	s.Equal(now.Unix(), s.state.NextRun(schn).Unix())

	now = now.AddDate(1, 0, 0)
	s.NoError(s.state.UpdateScheduler(schn, 20*time.Second, now))
	s.Equal(20*time.Second, s.state.Interval(schn))
	s.Equal(now.Unix(), s.state.NextRun(schn).Unix())
}

func (s *StateTestSuite) TestNextFaultIndex() {
	idx, _ := s.state.NextFaultIndex()
	s.Equal(int32(1), idx)
	idx, _ = s.state.NextFaultIndex()
	s.Equal(int32(2), idx)
	idx, _ = s.state.NextFaultIndex()
	s.Equal(int32(3), idx)
}

func (s *StateTestSuite) TestFaultSN() {
	s.state.InitAlertInfos(12)
	s.Equal(int64(1), s.state.GetFaultSn(12))
	s.NoError(s.state.IncrementFaultSn(12))
	s.Equal(int64(2), s.state.GetFaultSn(12))

	s.Equal(int64(0), s.state.GetFaultSn(42))
	s.Error(s.state.IncrementFaultSn(42))
}

func (s *StateTestSuite) TestFaultStore() {
	faultName := "mysuperdummyfault"
	s.Equal(int32(0), s.state.GetFaultInStorage(faultName))
	s.state.StoreFaultInStorage(faultName, 12)
	s.Equal(int32(12), s.state.GetFaultInStorage(faultName))
	s.state.DeleteFaultInStorage(faultName)
	s.Equal(int32(0), s.state.GetFaultInStorage(faultName))
}

func (s *StateTestSuite) TestFaultEpoch() {
	s.state.InitAlertInfos(12)
	s.Equal(int64(0), s.state.GetFaultStartEpoch(12))
	s.state.SetFaultStartEpoch(12, 12345)
	s.Equal(int64(12345), s.state.GetFaultStartEpoch(12))

	s.Equal(int64(0), s.state.GetFaultStartEpoch(42))
	s.state.SetFaultStartEpoch(42, 54321)
	s.Equal(int64(54321), s.state.GetFaultStartEpoch(42))
}
