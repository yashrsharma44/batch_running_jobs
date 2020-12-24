package pkg

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
)

func RegisterLogger() log.Logger {

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", time.Now().Format(time.RFC1123), "caller", log.DefaultCaller)

	return logger
}
