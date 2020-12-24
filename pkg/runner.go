package pkg

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const MAXDURATION int = 10000

type runnerFunction func(int, log.Logger, time.Duration) error

// Contains all the function that that imitates the long running job

func uploadWorker(dur int, logger log.Logger, delta time.Duration) error {
	level.Debug(logger).Log("msg", fmt.Sprintf("sleeping currently, duration %v/%v", dur+1, MAXDURATION))
	time.Sleep(delta)

	return nil
}

func downloadWorker(dur int, logger log.Logger, delta time.Duration) error {
	return uploadWorker(dur, logger, delta)
}

func longRunningWorker(dur int, logger log.Logger, delta time.Duration) error {
	return uploadWorker(dur, logger, delta)
}
