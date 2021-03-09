package golog

import (
	gosystemlog "log"

	log "gitlab.com/ovsinc/errors/log/common"
)

// New конструтор интерфейс для использования системного логгера log
func New(l *gosystemlog.Logger) log.Logger {
	return &systemlog{
		logger: l,
	}
}

type systemlog struct {
	logger *gosystemlog.Logger
}

func (l *systemlog) Warn(err error)  { l.logger.Print(err.Error()) }
func (l *systemlog) Error(err error) { l.logger.Print(err.Error()) }
