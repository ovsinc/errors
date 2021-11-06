package errors

import (
	"github.com/ovsinc/multilog"
)

func getLogger(l ...multilog.Logger) multilog.Logger {
	logger := multilog.DefaultLogger
	if len(l) > 0 {
		logger = l[0]
	}
	return logger
}

// LOG-хелперы

// CombineWithLog как и Combine создаст или дополнит цепочку ошибок err с помощью errs,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func CombineWithLog(errs ...error) error {
	e := Combine(errs...)
	Log(e)
	return e
}

// WrapWithLog обернет ошибку olderr в err и вернет цепочку,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func WrapWithLog(olderr error, err error) error {
	e := Wrap(olderr, err)
	Log(e)
	return e
}

// Log выполнить логгирование ошибки err с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func Log(err error, l ...multilog.Logger) {
	loggger := getLogger(l...)
	if err == nil || loggger == nil {
		return
	}
	loggger.Errorf(err.Error())
}
