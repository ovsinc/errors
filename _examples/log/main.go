package main

import (
	"time"

	"github.com/ovsinc/multilog"
	"github.com/ovsinc/multilog/chain"
	"github.com/ovsinc/multilog/journald"
	"github.com/ovsinc/multilog/logrus"
	origlogrus "github.com/sirupsen/logrus"
	"gitlab.com/ovsinc/errors"
)

func main() {
	now := time.Now()

	logrusLogger := logrus.New(origlogrus.New())

	multilog.DefaultLogger = logrusLogger

	err := errors.NewWithLog(
		"hello error",
		errors.SetSeverity(errors.SeverityWarn),
		errors.SetContextInfo(
			errors.CtxMap{
				"time": now,
			},
		),
	)

	err = err.WithOptions(
		errors.SetID("my id"),
		errors.AppendContextInfo("duration", time.Since(now)),
	)

	journalLogger := journald.New()

	chainLogger := chain.New(logrusLogger, journalLogger)

	err.Log(chainLogger)
}
