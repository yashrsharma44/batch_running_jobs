package pkg

import (
	"fmt"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

type Scheduler struct {
	taskList     map[ulid.ULID]*Task
	signalList   map[ulid.ULID]chan STATE
	errorList    map[ulid.ULID]chan error
	totalWorkers uint64

	mtx sync.RWMutex
}

func (schd *Scheduler) TotalTask() uint64 {
	schd.mtx.RLock()
	defer schd.mtx.RUnlock()

	return schd.totalWorkers
}

func (schd *Scheduler) IncrementTaskCount() {
	schd.mtx.RLock()
	defer schd.mtx.RUnlock()

	schd.totalWorkers += 1
}

func (schd *Scheduler) DecrementTaskCount() {
	schd.mtx.RLock()
	defer schd.mtx.RUnlock()

	schd.totalWorkers -= 1
}

func (schd *Scheduler) ScheduleTask(task *Task) error {
	// Initialise the task in the task list and create a channel mapping for the same id
	id := task.GetID()

	schd.taskList[id] = task
	schd.signalList[id] = make(chan STATE, 1)
	schd.errorList[id] = make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	go task.Run(schd.signalList[id], schd.errorList[id], &wg)

	schd.signalList[id] <- PLAY

	wg.Wait()
	// Check if any errors
	for val := range schd.errorList[id] {
		return val
	}

	return nil
}

func (schd *Scheduler) ModifyState(id ulid.ULID, newState STATE) error {

	if ok := schd.CheckTaskPresent(id); !ok {
		return fmt.Errorf("task for this id doesn't exist")
	}

	if schd.taskList[id].State() == STOP {
		return fmt.Errorf("task has been terminated, no change of state is allowed")
	}

	schd.signalList[id] <- newState
	// Wait for 2 * delta for the task to yield, so that we can get the new state applied
	time.Sleep(2 * schd.taskList[id].delta)
	return nil
}

// Check if the task is present in the scheduler
func (schd *Scheduler) CheckTaskPresent(id ulid.ULID) bool {
	_, ok := schd.taskList[id]
	return ok
}

func (schd *Scheduler) GetState(id ulid.ULID) (STATE, error) {
	if ok := schd.CheckTaskPresent(id); !ok {
		return NOT_RUNNING, fmt.Errorf("no task exists for the ID=%v", id)
	}
	return schd.taskList[id].State(), nil
}

func (schd *Scheduler) RemoveTask(id ulid.ULID) error {
	if ok := schd.CheckTaskPresent(id); !ok {
		return fmt.Errorf("cannot remove a task which is not present")
	}
	delete(schd.taskList, id)
	delete(schd.signalList, id)
	delete(schd.errorList, id)

	return nil
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		taskList:     make(map[ulid.ULID]*Task),
		signalList:   make(map[ulid.ULID]chan STATE),
		errorList:    make(map[ulid.ULID]chan error),
		totalWorkers: MAXWORKER,
	}
}
