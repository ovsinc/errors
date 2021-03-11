// Package chain реализаует цепочку логгирования.
package chain

import log "gitlab.com/ovsinc/errors/log/common"

// New конструктор логгера, реализующего цепочку из логгеров l.
// Можно указывать произвольное значение.
// Если l == nil, логгирование не будет осуществляться.
func New(l ...log.Logger) log.Logger {
	return &chainlog{
		loggers: append(make([]log.Logger, 0, len(l)), l...),
	}
}

type chainlog struct {
	loggers []log.Logger
}

func (l *chainlog) Warn(err error) {
	for _, log := range l.loggers {
		log.Warn(err)
	}
}

func (l *chainlog) Error(err error) {
	for _, log := range l.loggers {
		log.Error(err)
	}
}
