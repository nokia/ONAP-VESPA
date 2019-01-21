package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SchedulerTestSuite struct {
	suite.Suite
}

func TestScheduler(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}

func (s *SchedulerTestSuite) TestNullJob() {
	s.Panics(func() {
		NewScheduler("test", nil, 1*time.Second)
	})
}

func (s *SchedulerTestSuite) TestName() {
	defaultInterval := 1 * time.Second
	job := JobFunc(func(from, to time.Time, interval time.Duration) (interface{}, error) {
		return nil, nil
	})
	sched := NewScheduler("test", job, defaultInterval)
	s.Equal("test", sched.Name())
}

func (s *SchedulerTestSuite) TestFailedRun() {
	defaultInterval := 1 * time.Second
	job := JobFunc(func(from, to time.Time, interval time.Duration) (interface{}, error) {
		return nil, errors.New("blah blah")
	})
	sched := NewScheduler("test", job, defaultInterval)
	nextRun := sched.NextRun()
	_, err := sched.StepWait()
	s.Error(err)
	s.NoError(sched.Ack())
	s.Equal(nextRun, sched.NextRun())
}

func (s *SchedulerTestSuite) TestGetSetInterval() {
	defaultInterval := 5 * time.Second
	job := JobFunc(func(from, to time.Time, interval time.Duration) (interface{}, error) {
		return nil, nil
	})
	sched := NewScheduler("test", job, defaultInterval)

	s.Equal(defaultInterval, sched.GetInterval())
	s.NoError(sched.SetInterval(12 * time.Minute))
	s.Equal(12*time.Minute, sched.GetInterval())
	s.NoError(sched.SetInterval(12 * time.Minute))
	s.Equal(12*time.Minute, sched.GetInterval())
	s.NoError(sched.SetInterval(time.Duration(0)))
	s.Equal(defaultInterval, sched.GetInterval())
}

func (s *SchedulerTestSuite) TestStepWaitAck() {
	defaultInterval := 5 * time.Second
	var lastFrom, lastTo time.Time
	var lastInterval time.Duration
	job := JobFunc(func(from, to time.Time, interval time.Duration) (interface{}, error) {
		lastFrom = from
		lastTo = to
		lastInterval = interval
		return nil, nil
	})

	sched := NewScheduler("test", job, defaultInterval)
	sched.StepWait()
	nextRun := sched.NextRun()
	sched.StepWait()
	s.Equal(nextRun, lastFrom)
	s.Equal(defaultInterval, lastInterval)

	sched.StepWait()
	s.Equal(nextRun, lastFrom)
	s.Equal(defaultInterval, lastInterval)

	s.NoError(sched.Ack())
	s.Equal(nextRun.Add(defaultInterval), sched.NextRun())
	s.NoError(sched.Ack()) // Should not change anything
	s.Equal(nextRun.Add(defaultInterval), sched.NextRun())

	nextRun = sched.NextRun()
	sched.StepWait()
	s.Equal(nextRun, lastFrom)
	s.Equal(defaultInterval, lastInterval)

	s.NoError(sched.Ack())
	nextRun = sched.NextRun()
	time.Sleep(3 * defaultInterval)
	sched.StepWait()
	s.Equal(nextRun, lastFrom)
	s.True(lastTo.After(lastFrom.Add(2 * defaultInterval)))
}

func (s *SchedulerTestSuite) TestStepNotReady() {
	defaultInterval := 5 * time.Second
	job := JobFunc(func(from, to time.Time, interval time.Duration) (interface{}, error) {
		return nil, nil
	})

	sched := NewScheduler("test", job, defaultInterval)
	err := sched.state.UpdateNextRun("test", time.Now().Add(1*time.Minute))
	s.NoError(err)
	_, err = sched.Step()
	s.Error(err)
	s.Equal(ErrNotReady, err)
}

func (s *SchedulerTestSuite) TestWaitTimeout() {
	defaultInterval := 5 * time.Second
	job := JobFunc(func(from, to time.Time, interval time.Duration) (interface{}, error) {
		return nil, nil
	})

	sched := NewScheduler("test", job, defaultInterval)
	err := sched.state.UpdateNextRun("test", time.Now().Add(1*time.Minute))
	s.NoError(err)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	err = sched.Wait(ctx)
	cancel()
	s.Error(err)
}
