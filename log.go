package errors

import "github.com/ovsinc/multilog"

// DefaultLogger дефолтный логгер в пакете.
var DefaultLogger Logger = multilog.DefaultLogger //nolint:gochecknoglobals

// Logger используемый в пакете интерфейс логгера.
type Logger interface {
	Errorf(format string, args ...interface{})
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
