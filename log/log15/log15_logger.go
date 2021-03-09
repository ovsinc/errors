package log15

import (
	log15orig "github.com/inconshreveable/log15"

	log "gitlab.com/ovsinc/errors/log/common"
)

func New(l log15orig.Logger) log.Logger {
	return &log15logger{
		logger: l,
	}
}

type log15logger struct {
	logger log15orig.Logger
}

func (l *log15logger) Warn(err error)  { l.logger.Warn(err.Error()) }
func (l *log15logger) Error(err error) { l.logger.Error(err.Error()) }
