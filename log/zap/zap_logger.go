// Package zap реализует логгер zap.
package zap

import (
	log "gitlab.com/ovsinc/errors/log/common"
	origzap "go.uber.org/zap"
)

// New конструтор интерфейс для использования логгера zap
// Оборачивает логгер zap l.
func New(l *origzap.Logger) log.Logger {
	return &zaplogger{
		logger: l,
	}
}

type zaplogger struct {
	logger *origzap.Logger
}

func (l *zaplogger) Warn(err error)  { l.logger.Warn(err.Error()) }
func (l *zaplogger) Error(err error) { l.logger.Error(err.Error()) }
