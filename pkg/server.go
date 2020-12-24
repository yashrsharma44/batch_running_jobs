package pkg

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid"
)

type ServerComponent struct {
	srv    *http.Server
	logger log.Logger
	schd   *Scheduler
}

// General Server API Handler
func (comp *ServerComponent) CreateNewJob(w http.ResponseWriter, r *http.Request) {

	// Write Header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Slot is available, return a worker id which is enabled for 5 minutes.
	// Validity of the worker id checked when a job is submitted.
	worker_id := getULID().String()

	NewResponse(worker_id, MapStatetoMsg[NOT_RUNNING], http.StatusOK, "Slot is available, ID will expire in 5 minutes.", &w)
}

func (comp *ServerComponent) GetWorkerStatus(w http.ResponseWriter, r *http.Request) {

	// Write Header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := mux.Vars(r)
	isValid, message, ok := true, "", false

	// Parse ULID from string.
	worker_id, err := ulid.ParseStrict(id["worker-id"])
	if err != nil {
		isValid = false
		message = fmt.Sprintf("ULID parsing failed. err=%v", err)
	}

	var task *Task
	// Check if ULID has a running task, if not return error
	if isValid {
		task, ok = comp.schd.taskList[worker_id]
		if !ok {
			isValid = false
			message = fmt.Sprintf("No task running with this id=%v", id["worker-id"])
		}
	}

	// Get worker status of a running task. If the task is not running, or invalid, return error
	if isValid {
		NewResponse(id["worker-id"], MapStatetoMsg[task.State()], http.StatusOK, "Task running with the following state", &w)
		return
	}

	NewResponse(id["worker-id"], "Invalid", http.StatusBadRequest, message, &w)
}

func (comp *ServerComponent) ModifyWorkerStatus(w http.ResponseWriter, r *http.Request) {
	// Write Header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Make sure worker id is valid
	id := mux.Vars(r)
	worker_id := id["worker-id"]

	worker_ulid, err := ulid.ParseStrict(worker_id)
	if err != nil {
		NewResponse(worker_id, "", http.StatusBadRequest, fmt.Sprintf("worker id is invalid. err=%v", err), &w)
		return
	}

	NewState := MapParamtoState[id["new-job-status"]]
	OldState, err := comp.schd.GetState(worker_ulid)

	// Check if the worker id hasn't been used after being terminated.
	if err != nil {
		NewResponse(worker_id, MapStatetoMsg[NOT_RUNNING], http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err), &w)
		return
	}

	// Make sure we dont allow multiple requests for state change.
	if NewState == OldState {
		NewResponse(worker_id, MapStatetoMsg[NewState], http.StatusBadRequest, "new state of the task is the same as that of the old state", &w)
		return
	}

	// Modify the state change, return error if any.
	if err := comp.schd.ModifyState(worker_ulid, NewState); err != nil {
		NewResponse(worker_id, MapStatetoMsg[STOP], http.StatusBadRequest, fmt.Sprintf("%v", err), &w)
		return
	}

	NewResponse(worker_id, MapStatetoMsg[comp.schd.taskList[worker_ulid].State()], http.StatusOK, "changed the state of the task", &w)

	// If the task is terminated, remove it from the scheduler and increase task count.
	if NewState == STOP {
		_ = comp.schd.RemoveTask(worker_ulid)
		comp.schd.IncrementTaskCount()
	}
}

// All the task handlers share the same logic.
func (comp *ServerComponent) taskHelper(w *http.ResponseWriter, r *http.Request, fn runnerFunction) {

	// Write Header
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Make sure worker id is valid
	id := mux.Vars(r)
	worker_id := id["worker-id"]

	worker_ulid, err := ulid.ParseStrict(worker_id)
	if err != nil {
		NewResponse(worker_id, "", http.StatusBadRequest, fmt.Sprintf("worker id is invalid. err=%v", err), w)
		return
	}

	// Check if any slot is available
	totalWorker := comp.schd.TotalTask()
	fmt.Printf("SLOTS %v\n", totalWorker)
	if totalWorker == 0 || totalWorker > MAXWORKER {
		NewResponse("", MapStatetoMsg[NOT_RUNNING], http.StatusServiceUnavailable, "Server has no available slots, try again later.", w)
		return
	}

	// Check if the task has been started before
	if ok := comp.schd.CheckTaskPresent(worker_ulid); ok {
		state, _ := comp.schd.GetState(worker_ulid)
		NewResponse(worker_id, MapStatetoMsg[state], http.StatusBadRequest, "the task has already been started", w)
		return
	}

	// Check if the worker id has not expired
	if err := checkTimeValid(worker_ulid); err != nil {
		NewResponse(worker_id, MapStatetoMsg[NOT_RUNNING], http.StatusBadRequest, "worker id has expired, please use another one", w)
		return
	}

	// Initialise the count for the new task
	comp.schd.DecrementTaskCount()

	// Create a task
	// Assume we need to run a task which has 10,000 entries to be uploaded, and each entries take 1 second.
	// This is totally modifiable from the client side, although for demonstration purposes, we have simulated this
	// The total time can be changed to completion of the task, although it should be interruptible.
	taskObj := NewTask(worker_ulid, comp.logger, fn, 10000*time.Second, 1*time.Second)

	// Schedule it
	// Wait for completion
	// If issues return error
	if err := comp.schd.ScheduleTask(taskObj); err != nil {
		NewResponse(worker_id, MapStatetoMsg[STOP], http.StatusInternalServerError, fmt.Sprintf("Internal Server Error. err=%v", err), w)
		return
	}

	// else return success and remove the task from the scheduler
	comp.schd.IncrementTaskCount()
	state, _ := comp.schd.GetState(worker_ulid)
	_ = comp.schd.RemoveTask(worker_ulid)
	NewResponse(worker_id, MapStatetoMsg[state], http.StatusOK, "Task has been successfully completed", w)

}

// Upload API Handler
func (comp *ServerComponent) CreateWorkerHandler(w http.ResponseWriter, r *http.Request) {
	comp.taskHelper(&w, r, uploadWorker)
}

// Download API Handler
func (comp *ServerComponent) DownloadWorkerHandler(w http.ResponseWriter, r *http.Request) {
	comp.taskHelper(&w, r, downloadWorker)
}

// Long Running task API Handler
func (comp *ServerComponent) HandleJob(w http.ResponseWriter, r *http.Request) {
	comp.taskHelper(&w, r, longRunningWorker)
}

func NewServerComponent(srv *http.Server, logger log.Logger) *ServerComponent {
	return &ServerComponent{
		srv:    srv,
		logger: logger,
		schd:   NewScheduler(),
	}
}
