package main

import (
	"time"

	origlogrus "github.com/sirupsen/logrus"
	"gitlab.com/ovsinc/errors"
	"gitlab.com/ovsinc/errors/log"
	"gitlab.com/ovsinc/errors/log/chain"
	"gitlab.com/ovsinc/errors/log/journald"
	"gitlab.com/ovsinc/errors/log/logrus"
)

func main() {
	now := time.Now()

	logrusLogger := logrus.New(origlogrus.New())

	log.DefaultLogger = logrusLogger

	err := errors.NewWithLog(
		"hello error",
		errors.SetSeverity(log.SeverityWarn),
		errors.SetContextInfo(
			errors.CtxMap{
				"time": now,
			},
		),
	)

	err = err.WithOptions(
		errors.SetSeverity(log.SeverityError),
		errors.AppendContextInfo("duration", time.Since(now)),
	)

	journalLogger := journald.New()

	chainLogger := chain.New(logrusLogger, journalLogger)

	err.Log(chainLogger)
}
