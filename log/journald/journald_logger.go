package journald

import (
	pkgjournal "github.com/coreos/go-systemd/journal"

	log "gitlab.com/ovsinc/errors/log/common"
)

func New() log.Logger {
	return &journaldlogger{}
}

type journaldlogger struct{}

func (l *journaldlogger) Warn(err error) {
	_ = pkgjournal.Send(err.Error(), pkgjournal.PriWarning, nil)
}

func (l *journaldlogger) Error(err error) {
	_ = pkgjournal.Send(err.Error(), pkgjournal.PriErr, nil)
}
