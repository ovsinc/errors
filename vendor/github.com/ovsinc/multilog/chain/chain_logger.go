// Package chain реализаует цепочку логгирования.
package chain

import log "github.com/ovsinc/multilog/common"

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
