package errors

func NewErrorType(s string) ErrorType {
	return ErrorType(s)
}

// ErrorType тип ошибки
type ErrorType string

func (e ErrorType) String() string {
	return string(e)
}

// UnknownErrorType неизвестная ошибка - дефолтное значение типа ошибки
var UnknownErrorType = NewErrorType("UNKNOWN_TYPE")

//

// SetErrorType установить тип ошибки
func SetErrorType(etype ErrorType) Options {
	return func(e *Error) {
		if e == nil {
			return
		}
		e.errorType = etype
	}
}

//

// ErrorType получить тип ошибки
func (e *Error) ErrorType() ErrorType {
	return e.errorType
}
