package errors

import (
	origerrors "errors"

	"github.com/hashicorp/go-multierror"
	"gitlab.com/ovsinc/errors/log"
	logcommon "gitlab.com/ovsinc/errors/log/common"
)

//

func customlog(l logcommon.Logger, e error, severity log.Severity) {
	if e == nil {
		return
	}

	switch severity {
	case log.SeverityError:
		l.Error(e)

	case log.SeverityWarn:
		l.Warn(e)

	case log.SeverityEnds, log.SeverityUnknown:
		l.Error(e)

	default:
		l.Error(e)
	}
}

func getLogger(l ...logcommon.Logger) logcommon.Logger {
	logger := log.DefaultLogger
	if len(l) > 0 {
		logger = l[0]
	}
	return logger
}

// хелперы

// AppendWithLog как и Append создаст или дополнит цепочку ошибок err с помощью errs,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func AppendWithLog(err error, errs ...error) *multierror.Error {
	e := Append(err, errs...)
	Log(e)
	return e
}

// WrapWithLog обернет ошибку olderr в err и вернет цепочку,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func WrapWithLog(olderr error, err error) *multierror.Error {
	e := Wrap(olderr, err)
	Log(e)
	return e
}

// Log выполнить логгирование ошибки err с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func Log(err error, l ...logcommon.Logger) {
	var errseverity *Error
	severity := log.SeverityError

	if origerrors.As(err, &errseverity) {
		severity = errseverity.Severity()
	}
	customlog(getLogger(l...), err, severity)
}

// NewWithLog конструктор *Error, как и New,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func NewWithLog(msg string, ops ...Options) Errorer {
	e := New(msg, ops...)
	e.Log()
	return e
}

// дополнительные методы *Error

// Log выполнит логгирование ошибки с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func (e *Error) Log(l ...logcommon.Logger) {
	customlog(getLogger(l...), e, e.Severity())
}
