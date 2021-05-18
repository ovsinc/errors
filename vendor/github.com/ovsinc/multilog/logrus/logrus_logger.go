// Package logrus реализует логгер logrus.
package logrus

import (
	log "github.com/ovsinc/multilog/common"

	origlogrus "github.com/sirupsen/logrus"
)

// New конструтор интерфейс для использования логгера logrus
// Оборачивает logrus логгер l.
func New(l *origlogrus.Logger) log.Logger {
	return &logruslogger{
		logger: l,
	}
}

type logruslogger struct {
	logger *origlogrus.Logger
}

func (l *logruslogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logruslogger) Infof(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logruslogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logruslogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
