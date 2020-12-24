package api

import (
	"github.com/gorilla/mux"
	"github.com/yashrsharma44/batch_running_jobs/pkg"
)

func RegisterRoutes(r *mux.Router, srvComp *pkg.ServerComponent) {
	r.Use(mux.CORSMethodMiddleware(r))

	// General Server API
	r.HandleFunc("/new-job/", pkg.PermissiveCORSMiddleware(srvComp.CreateNewJob)).Methods("GET", "OPTIONS")
	r.HandleFunc("/modify/{worker-id}/{new-job-status}", pkg.PermissiveCORSMiddleware(srvComp.ModifyWorkerStatus)).Methods("GET", "OPTIONS")
	r.HandleFunc("/status/{worker-id}", pkg.PermissiveCORSMiddleware(srvComp.GetWorkerStatus)).Methods("GET", "OPTIONS")

	// Upload API
	r.HandleFunc("/upload/new/{worker-id}", pkg.PermissiveCORSMiddleware(srvComp.CreateWorkerHandler)).Methods("POST", "OPTIONS")

	// // Download API
	r.HandleFunc("/download/new/{worker-id}", pkg.PermissiveCORSMiddleware(srvComp.DownloadWorkerHandler)).Methods("POST", "OPTIONS")

	// // Running a task API
	r.HandleFunc("/job/new/{worker-id}", pkg.PermissiveCORSMiddleware(srvComp.HandleJob)).Methods("POST", "OPTIONS")

	// Handler for main page
	r.PathPrefix("/").HandlerFunc(pkg.PermissiveCORS).Methods("GET", "OPTIONS")
}
