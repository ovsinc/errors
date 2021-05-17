package errors

// Options опции из параметра ошибки.
type Options func(e *Error)

// SetFormatFn установит пользовательскую функцию-форматирования
func SetFormatFn(fn FormatFn) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.formatFn = fn
	}
}

// SetMsg установит сообщение об ошибке.
func SetMsg(msg string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.msg = msg
	}
}

// SetSeverity устновит Severity.
func SetSeverity(severity Severity) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.severity = severity
	}
}

// SetID установит ID ошибки.
func SetID(id string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.id = id
	}
}

// SetErrorType установит тип ошибки
func SetErrorType(etype string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = etype
	}
}
