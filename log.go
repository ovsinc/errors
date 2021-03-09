package errors

import (
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

//

func AppendWithLog(err error, errs ...error) *multierror.Error {
	e := Append(err, errs...)
	Log(e)
	return e
}

func WrapWithLog(olderr error, err error) *multierror.Error {
	e := Wrap(olderr, err)
	Log(e)
	return e
}

func Log(err error, l ...logcommon.Logger) {
	severity := log.SeverityError
	type severitier interface {
		Severity() log.Severity
	}
	if customerr, ok := err.(severitier); ok {
		severity = customerr.Severity()
	}
	customlog(getLogger(l...), err, severity)
}

func NewWithLog(msg string, ops ...Options) *Error {
	e := New(msg, ops...)
	e.Log()
	return e
}

//

func (e *Error) Log(l ...logcommon.Logger) {
	customlog(getLogger(l...), e, e.Severity())
}
