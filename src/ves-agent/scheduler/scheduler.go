package scheduler

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

// ErrNotReady is the error returned when the scheduler is not ready
// meaning it's not the time yet to run
var ErrNotReady = errors.New("Scheduler not ready")

// State handles schedulers internal state
type State interface {
	// NextRun returns the time at which next execution should occure
	NextRun(name string) time.Time
	// UpdateNextRun set the time of the next execution
	UpdateNextRun(name string, next time.Time) error
	// Interval returns the scheduler exceution interval
	Interval(name string) time.Duration
	// UpdateInterval set a new execution interval for the scheduler
	UpdateInterval(name string, interval time.Duration) error
	// UpdateScheduler set both interval and next execution time for the scheduler
	UpdateScheduler(name string, interval time.Duration, next time.Time) error
}

type inMemState struct {
	next     time.Time
	interval time.Duration
}

func (state *inMemState) NextRun(_ string) time.Time {
	return state.next
}
func (state *inMemState) UpdateNextRun(_ string, next time.Time) error {
	state.next = next
	return nil
}
func (state *inMemState) Interval(_ string) time.Duration {
	return state.interval
}
func (state *inMemState) UpdateInterval(_ string, interval time.Duration) error {
	state.interval = interval
	return nil
}
func (state *inMemState) UpdateScheduler(_ string, interval time.Duration, next time.Time) error {
	state.interval = interval
	state.next = next
	return nil
}

// Job is an interface to a schedulable task
type Job interface {
	// Run the job, providing it with timing information
	Run(from, to time.Time, interval time.Duration) (interface{}, error)
}

// JobFunc is a function implementing the Job interface
type JobFunc func(from, to time.Time, interval time.Duration) (interface{}, error)

// Run the job by calling the function
func (job JobFunc) Run(from, to time.Time, interval time.Duration) (interface{}, error) {
	return job(from, to, interval)
}

// Scheduler schedules a job on a periodic interval. It provides on each Job run
// some timing information about the interval being invoked.
// After a run, the scheduler must receive Akcnowledge to commit the interval
//
// WARNING: The Scheduler is not thread safe and must not be used as is
// (without additional care) from multiple goroutines
type Scheduler struct {
	name            string        // The scheduler name
	defaultInterval time.Duration // Default interval between each run
	job             Job           // Job to periodically execute
	lastTime        *time.Time    // Last time of successful, unacknowleged run
	state           State         // Scheduler internal state
}

// NewSchedulerWithState creates a new scheduler with the given Job, default interval, using the procided
// internal state handler
func NewSchedulerWithState(name string, job Job, defaultInterval time.Duration, state State) *Scheduler {
	if job == nil {
		log.Panic("job cannot be null")
	}
	log.Infof("Creating Scheduler '%s' (default interval: %s)", name, defaultInterval.String())
	return &Scheduler{name: name, defaultInterval: defaultInterval, job: job, state: state}
}

// NewScheduler creates a new scheduler with the given Job, and default interval.
// Scheduler internal state is handled in memory
func NewScheduler(name string, job Job, defaultInterval time.Duration) *Scheduler {
	return NewSchedulerWithState(name, job, defaultInterval, &inMemState{})
}

// Name returns the name of this scheduler
func (sched *Scheduler) Name() string {
	return sched.name
}

// GetInterval returns the scheduled interval set,
// or the default one
func (sched *Scheduler) GetInterval() time.Duration {
	interval := sched.state.Interval(sched.name)
	if interval > 0 {
		return interval
	}
	return sched.defaultInterval
}

// SetInterval modifies the scheduled interval
func (sched *Scheduler) SetInterval(interval time.Duration) error {
	if interval == sched.GetInterval() {
		// Nothing changed
		return nil
	}
	var newInterval time.Duration
	if interval <= 0 {
		newInterval = sched.defaultInterval
	} else {
		newInterval = interval
	}

	next := sched.NextRun()
	if next.After(time.Now()) {
		next = time.Now().Round(newInterval)
	} else {
		next = next.Round(newInterval)
	}
	if err := sched.state.UpdateScheduler(sched.name, newInterval, next); err != nil {
		return err
	}
	log.Infof("Scheduler %s: interval set to %s", sched.name, newInterval)
	log.Infof("Scheduler %s: next run at %s", sched.name, sched.NextRun().String())
	return nil
}

// NextRun returns the time of the next run
func (sched *Scheduler) NextRun() time.Time {
	nextTime := sched.state.NextRun(sched.name)
	if nextTime.IsZero() {
		// Lazy initialization
		nextTime = time.Now().Round(sched.defaultInterval)
		if err := sched.state.UpdateNextRun(sched.name, nextTime); err != nil {
			// Ignore error at lazy init
			log.Warnf("Scheduler '%s': Cannot update next run time: %s", sched.name, err.Error())
		}
		log.Infof("Scheduler %s: next run at %s (in %s)", sched.name, nextTime.String(), time.Until(nextTime).String())
		return nextTime
	}
	return nextTime
}

// Ready returns true if it's time for the next run
func (sched *Scheduler) Ready() bool {
	now := time.Now()
	next := sched.NextRun()
	return next.Before(now) || next.Equal(now)
}

// WaitChan returns a timer.Timer set to
// wait for the next run
func (sched *Scheduler) WaitChan() *time.Timer {
	if sched.Ready() {
		return time.NewTimer(0)
	}
	sleepTime := time.Until(sched.NextRun())
	log.Debugf("Scheduler %s: waiting for %s", sched.name, sleepTime.String())
	return time.NewTimer(sleepTime)
}

// Wait the necessary time before the next run becomes possible
func (sched *Scheduler) Wait(ctx context.Context) error {
	select {
	case <-sched.WaitChan().C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Ack acknowledges the last run.
// It does nothing if last run did not succeed, or
// has already been aknowledged
func (sched *Scheduler) Ack() error {
	if sched.lastTime == nil {
		return nil
	}
	next := sched.lastTime.Truncate(sched.GetInterval()).Add(sched.GetInterval())
	log.Debugf("Scheduler %s Ack: Updating next run time from %s to %s", sched.name, sched.NextRun().String(), next.String())
	if err := sched.state.UpdateNextRun(sched.name, next); err != nil {
		return err
	}
	sched.lastTime = nil
	log.Infof("Scheduler %s: Next run in %s", sched.name, time.Until(sched.NextRun()))
	return nil
}

// Step will execute the next round,
// starting from the last Acknowledged run
//
// It won't wait until it's time to be run. If it's not the time, it will return an error.
// For a blocking alternative, see StepWait()
func (sched *Scheduler) Step() (interface{}, error) {
	sched.lastTime = nil // Reset last successful query time
	if !sched.Ready() {
		return nil, ErrNotReady
	}
	now := time.Now()
	res, err := sched.job.Run(sched.NextRun(), now, sched.GetInterval())
	if err != nil {
		return nil, err
	}
	sched.lastTime = &now // Set last successful query time
	return res, nil
}

// StepWait execute the next round, waiting for the proper time if it's necessary
func (sched *Scheduler) StepWait() (interface{}, error) {
	sched.lastTime = nil // Reset last successful query time
	if err := sched.Wait(context.Background()); err != nil {
		return nil, err
	}
	return sched.Step()
}
