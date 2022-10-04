package errors

import "github.com/ovsinc/multilog"

var DefaultLogger = NewLogger()

type Logger interface {
	Errorf(format string, args ...interface{})
}

func NewLogger(l ...multilog.Logger) Logger {
	logger := multilog.DefaultLogger
	if len(l) > 0 {
		logger = l[0]
	}
	return logger
}

// LOG-хелперы

// Log выполнить логгирование ошибки err с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func Log(err error, lg ...Logger) {
	l := DefaultLogger
	if len(lg) > 0 {
		l = lg[0]
	}
	l.Errorf(err.Error())
}
