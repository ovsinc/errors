package errors

import (
	"gitlab.com/ovsinc/errors/log"
)

// Options опиции из параметра ошибки
type Options func(e *Error)

// SetFormatFn установить пользовательскую функцию-форматирования
func SetFormatFn(fn FormatFn) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.formatFn = fn
	}
}

// SetMsg установить сообщение
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = msg
	}
}

// SetSeverity ...
func SetSeverity(severity log.Severity) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.severity = severity
	}
}
