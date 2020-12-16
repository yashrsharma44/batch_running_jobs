package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
)

func main() {

	// Initialise the logger
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", time.Now().Format(time.RFC1123), "caller", log.DefaultCaller)

	// Set up the app logic
	batch :=

	// Set up the routers
	r := mux.NewRouter()
	r.HandleFunc("/upload/new/{worker_id}", batch.WorkerHandler)
	r.HandleFunc("/upload/{worker_id}/{new_job_status}", batch.ModifyWokerStatus)
	r.HandleFunc("/", batch.MainPageHandler)

	// Set up the server
	level.Info(logger).Log("msg", "starting the server at 9090")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
	}

	if err := srv.ListenAndServe(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	// Cleanup and shutdown

}
