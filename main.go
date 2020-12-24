package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/yashrsharma44/batch_running_jobs/api"
	"github.com/yashrsharma44/batch_running_jobs/pkg"
)

func handleSigterm(c chan os.Signal, cancel context.CancelFunc) {
	<-c
	cancel()
}

func main() {

	// Initialise the logger
	logger := pkg.RegisterLogger()

	// Set up the server
	level.Info(logger).Log("msg", "starting the server at 8000")
	r := mux.NewRouter().StrictSlash(true)
	srv := &http.Server{Handler: r,
		Addr: "127.0.0.1:8000",
	}
	batchJob := pkg.NewServerComponent(srv, logger)

	// Set up the routers
	api.RegisterRoutes(r, batchJob)

	if err := srv.ListenAndServe(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

}
