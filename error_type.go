package errors

// UnknownErrorType неизвестная ошибка - дефолтное значение типа ошибки
var UnknownErrorType = "UNKNOWN_TYPE" //nolint:gochecknoglobals

//

// SetErrorType установит тип ошибки
func SetErrorType(etype string) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = etype
	}
}

//

// ErrorType вернет тип ошибки
func (e *Error) ErrorType() string {
	return e.errorType
}
