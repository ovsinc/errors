package logrus

import (
	origlogrus "github.com/sirupsen/logrus"

	log "gitlab.com/ovsinc/errors/log/common"
)

// New конструтор интерфейс для использования логгера logrus
func New(l *origlogrus.Logger) log.Logger {
	return &logruslogger{
		logger: l,
	}
}

type logruslogger struct {
	logger *origlogrus.Logger
}

func (l *logruslogger) Warn(err error)  { l.logger.Warn(err.Error()) }
func (l *logruslogger) Error(err error) { l.logger.Error(err.Error()) }
