package pkg

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/ulid"
)

type Task struct {
	id     ulid.ULID
	state  STATE
	logger log.Logger
	fn     func(int, log.Logger, time.Duration) error
	// Condition for running the long running task.
	// Can be substituted with a completion check once task is completed
	maxDuration time.Duration
	// duration for which task gets interrupted or yields
	delta time.Duration
}

func (t *Task) GetID() ulid.ULID {
	return t.id
}

func (t *Task) State() STATE {
	return t.state
}

func (t *Task) Run(signalState chan STATE, errorChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(errorChan)
	defer close(signalState)

	t.state = NOT_RUNNING
	for dr := 0; dr < int(t.maxDuration.Seconds()); {
		select {
		case t.state = <-signalState:
			switch t.state {
			case PLAY:
				level.Debug(t.logger).Log("msg", "task is running..")
			case PAUSE:
				level.Debug(t.logger).Log("msg", "task has been paused..")
			case STOP:
				level.Debug(t.logger).Log("msg", "task has been terminated..")
				errorChan <- fmt.Errorf("the task has been prematurely terminated")
				t.state = NOT_RUNNING
				return
			}
		default:
			if t.state == PAUSE {
				break
			}
			// Here the long running function should be non-blocking, and should/can yield at any delta second
			// so that processor can switch tasks.
			// This makes sure that the process when interrupted at time x, terminates after x + delta time
			// If the process is blocking there is no way to interrupt the process then. Will be terminated
			//  after completed.
			if err := t.fn(dr, t.logger, t.delta); err != nil {

				errorChan <- err
				level.Debug(t.logger).Log("task has been terminated due to error. err", err)
				return
			} // expected to block for 1 or delta second.
			dr++
		}
	}
	t.state = COMPLETED
}

func NewTask(id ulid.ULID, logger log.Logger, fn func(int, log.Logger, time.Duration) error, dur time.Duration, delta time.Duration) *Task {
	return &Task{
		id:          id,
		state:       NOT_RUNNING,
		logger:      log.With(logger, "task-id", id),
		fn:          fn,
		maxDuration: dur,
		delta:       delta,
	}
}
