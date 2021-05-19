// +build !window

// Package journal implements the journal logger.
//
//
// RU:
// Package journald реализует логгер journald.
package journald

import (
	"fmt"

	log "github.com/ovsinc/multilog/common"

	pkgjournal "github.com/coreos/go-systemd/journal"
)

// New constructor of the custom journald logger.
//
//
// RU:
// New конструктор journald логгера.
func New() log.Logger {
	return &journaldlogger{}
}

type journaldlogger struct{}

func (l *journaldlogger) log(pri pkgjournal.Priority, format string, args ...interface{}) {
	_ = pkgjournal.Send(fmt.Sprintf(format, args...), pri, nil)
}

func (l *journaldlogger) Debugf(format string, args ...interface{}) {
	l.log(pkgjournal.PriDebug, format, args...)
}

func (l *journaldlogger) Infof(format string, args ...interface{}) {
	l.log(pkgjournal.PriInfo, format, args...)
}

func (l *journaldlogger) Warnf(format string, args ...interface{}) {
	l.log(pkgjournal.PriWarning, format, args...)
}

func (l *journaldlogger) Errorf(format string, args ...interface{}) {
	l.log(pkgjournal.PriErr, format, args...)
}
