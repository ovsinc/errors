// +build !windows,!plan9

// Package syslog реализует логгер syslog.
package syslog

import (
	gosystemsyslog "log/syslog"

	log "gitlab.com/ovsinc/errors/log/common"
)

// New конструтор интерфейс для использования логгера syslog.
// Оборачивает writer интерфейс.
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
