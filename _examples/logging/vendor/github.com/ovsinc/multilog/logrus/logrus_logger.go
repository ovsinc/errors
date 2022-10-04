// Package logrus implements the logrus logger.
//
// RU:
// Package logrus реализует логгер logrus.
package logrus

import (
	log "github.com/ovsinc/multilog/common"

	origlogrus "github.com/sirupsen/logrus"
)

// New constructor of a logger that wraps the original logger.
//
// RU:
// New конструтор интерфейс для использования логгера logrus
// Оборачивает logrus логгер l.
func New(l ...*origlogrus.Logger) log.Logger {
	logger := origlogrus.New()
	if len(l) > 0 {
		logger = l[0]
	}
	return &logruslogger{
		logger: logger,
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
