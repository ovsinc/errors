package errors

import (
	"github.com/davecgh/go-spew/spew"
	"gitlab.com/ovsinc/errors/log"
)

var _ error = (*Error)(nil)

// Error структура кастомной ошибки
// Это НЕ потоко-безопасный объект
type Error struct {
	severity     log.Severity
	operations   []Operation
	formatFn     FormatFn
	errorType    ErrorType
	msg          string
	contextInfo  CtxMap
	translateMap translateMap
	lang         string
}

// New конструктор на необязательных параметрах
// * ops ...Options -- параметризация через функции-парметры
//
// ** *Error
func New(msg string, ops ...Options) *Error {
	e := &Error{
		operations: []Operation{},
		severity:   log.SeverityError,
		errorType:  UnknownErrorType,
		msg:        msg,
		lang:       "en",
	}
	for _, op := range ops {
		op(e)
	}
	return e
}

// setters

// WithOptions производит параметризацию *Error
// * ops ...Options - параметризация через функции-парметры
//
// ** *Error - возвращает модифицированный экземпляр *Error
func (e *Error) WithOptions(ops ...Options) *Error {
	for _, op := range ops {
		op(e)
	}
	return e
}

// getters

// Severity получить критичность ошибки
func (e *Error) Severity() log.Severity {
	return e.severity
}

// Msg получить исходное сообщение об ошибке
func (e *Error) Msg() string {
	return e.msg
}

// Error получить описание ошибки
// метод для реализации интерфейса error
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	fn := e.formatFn
	if e.formatFn == nil {
		fn = defaultFormatFn
	}
	return fn(e)
}

// дополнительные методы

// Sdump получить текстовый дамп ошибки
func (e *Error) Sdump() string {
	if e == nil {
		return ""
	}
	return spew.Sdump(e)
}

// ErrorOrNil получить ошибку или nil
// ошибкой считается *Error != nil и Severity == SeverityError
// т.е. SeverityWarn НЕ ошибка
func (e *Error) ErrorOrNil() error {
	if e == nil {
		return nil
	}
	if e.Severity() != log.SeverityError {
		return nil
	}
	return e
}
