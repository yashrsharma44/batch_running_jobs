package main

import "net/http"

type BatchJob struct {
}

func (job *BatchJob) WorkerHandler() {

}

func (job *BatchJob) ModifyWorkerStatus() {

}

func (job *BatchJob) MainPageHandler(w http.ResponseWriter, r *http.Request) {

}

func NewBatchJob() *BatchJob {
	return &BatchJob{}
}
