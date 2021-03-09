package syslog

import (
	gosystemsyslog "log/syslog"

	log "gitlab.com/ovsinc/errors/log/common"
)

// New конструтор интерфейс для использования системного логгера log
func New(w *gosystemsyslog.Writer) log.Logger {
	return &sysloglogger{
		writer: w,
	}
}

type sysloglogger struct {
	writer *gosystemsyslog.Writer
}

func (l *sysloglogger) Warn(err error)  { _ = l.writer.Warning(err.Error()) }
func (l *sysloglogger) Error(err error) { _ = l.writer.Err(err.Error()) }
