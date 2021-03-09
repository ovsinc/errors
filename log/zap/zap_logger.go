package zap

import (
	log "gitlab.com/ovsinc/errors/log/common"

	origzap "go.uber.org/zap"
)

type zaplogger struct {
	logger *origzap.Logger
}

// New конструтор интерфейс для использования логгера zap
func New(l *origzap.Logger) log.Logger {
	return &zaplogger{
		logger: l,
	}
}

func (l *zaplogger) Warn(err error)  { l.logger.Warn(err.Error()) }
func (l *zaplogger) Error(err error) { l.logger.Error(err.Error()) }
