// Package chain implements a logging chain.
//
// RU:
// Package chain реализует цепочку логгирования.
package chain

import log "github.com/ovsinc/multilog/common"

// New constructor of a logger that implements a chain of loggers l.
// You can specify an arbitrary number.
// If l == nil, logging will not be performed.
//
// New конструктор логгера, реализующего цепочку из логгеров l.
// Можно указывать произвольное количество.
// Если l == nil, логгирование не будет осуществляться.
func New(l ...log.Logger) log.Logger {
	return &chainlog{
		loggers: append(make([]log.Logger, 0, len(l)), l...),
	}
}

type chainlog struct {
	loggers []log.Logger
}

func (l *chainlog) Debugf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		log.Debugf(format, args...)
	}
}

func (l *chainlog) Infof(format string, args ...interface{}) {
	for _, log := range l.loggers {
		log.Infof(format, args...)
	}
}

func (l *chainlog) Warnf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		log.Warnf(format, args...)
	}
}

func (l *chainlog) Errorf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		log.Errorf(format, args...)
	}
}
