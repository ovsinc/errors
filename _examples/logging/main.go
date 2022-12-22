package main

import (
	"time"

	"github.com/ovsinc/errors"
	"github.com/ovsinc/multilog/chain"
	"github.com/ovsinc/multilog/journald"
	"github.com/ovsinc/multilog/logrus"
)

func main() {
	now := time.Now()

	logrusLogger := logrus.New()
	errors.DefaultLogger = logrusLogger

	err := errors.NewWithLog(
		errors.SetMsg("hello error"),
		errors.AppendContextInfo("time", now.Format("2006-01-02T15:04:05-0700")),
	)

	err = err.WithOptions(
		errors.SetID("my id"),
		errors.AppendContextInfo("duration", time.Since(now).String()),
	)

	journalLogger := journald.New()

	chainLogger := chain.New(logrusLogger, journalLogger)

	err.Log(chainLogger)
}
