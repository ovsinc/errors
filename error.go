package errors

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
)

var (
	ErrUnknownMarshaller = New("marshaller not set")

	_ error         = (*Error)(nil)
	_ fmt.Formatter = (*Error)(nil)
)

// CtxKV slice key-value контекста ошибки.
type CtxKV []struct {
	Key   string
	Value interface{}
}

// New конструктор *Error. Аргумент используется для указания сообщения.
// * In: interface{}
// Поддерживаемые варианты для NewWith(SetMsg(msg)):
// string, *Error -> as is
// error -> Error()
// interface{ String() string } -> String()
// func() string -> f()
// * Out: *Error
//
// ** *Error
func New(i interface{}) *Error {
	var msg string
	switch t := i.(type) {
	case string:
		msg = t
	case *Error:
		return t
	case error:
		msg = t.Error()
	case interface{ String() string }:
		msg = t.String()
	case func() string:
		msg = t()
	}
	return NewWith(SetMsg(msg))
}

// NewLog конструктор *Error, как и New,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func NewLog(i interface{}) *Error {
	e := New(i)
	e.Log()
	return e
}

// NewWith конструктор на необязательных параметрах
// * ops ...Options -- параметризация через функции-парметры.
// См. options.go
//
// ** *Error
func NewWith(ops ...Options) *Error {
	e := Error{}
	for _, op := range ops {
		op(&e)
	}
	return &e
}

// NewWithLog конструктор *Error, как и NewWith,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func NewWithLog(ops ...Options) *Error {
	e := NewWith(ops...)
	e.Log()
	return e
}

// Error структура кастомной ошибки.
type Error struct {
	id, msg, operation string
	contextInfo        CtxKV
	errorType          IErrType
}

// WithOptions производит параметризацию *Error с помощью функции-парметры Options.
// Допускается указывать произвольно количество ops.
// Возвращается новый экземпляр *Error с переопределением заданных параметров.
func (e *Error) WithOptions(ops ...Options) *Error {
	if e == nil {
		return nil
	}

	// copy *Error
	newerr := new(Error)
	*newerr = *e

	for _, op := range ops {
		op(newerr)
	}

	return newerr
}

// getters

// ID возвращает ID ошибки.
func (e *Error) ID() string {
	if e == nil {
		return ""
	}
	return e.id
}

// Msg возвращает исходное сообщение об ошибке.
func (e *Error) Msg() string {
	if e == nil {
		return ""
	}
	return e.msg
}

// Operations вернет список операций.
func (e *Error) Operation() string {
	if e == nil {
		return ""
	}
	return e.operation
}

// ErrorType вернет тип ошибки.
func (e *Error) ErrorType() IErrType {
	if e == nil {
		return defaultErrType
	}

	return e.errorType
}

// ContextInfo вернет контекст CtxKV ошибки.
func (e *Error) ContextInfo() CtxKV {
	return e.contextInfo
}

// методы форматирования

func mustMarshaler(fn ...Marshaller) Marshaller {
	var marshal Marshaller
	switch {
	case len(fn) > 0:
		marshal = fn[0]
	case DefaultMarshaller != nil:
		marshal = DefaultMarshaller
	default:
		panic(ErrUnknownMarshaller)
	}
	return marshal
}

// Marshal метод маршалит *Error.
// * fn ...Marshaller -- необязательный парамет для вызова кастомного маршалера
// если не указано, используется дефолтный.
// ** []byte, error
func (e *Error) Marshal(fn ...Marshaller) ([]byte, error) {
	marshal := mustMarshaler(fn...)
	return marshal.Marshal(e)
}

// Format производит форматирование строки, для поддержки fmt.Printf().
func (e *Error) Format(s fmt.State, verb rune) { //nolint:cyclop
	if e == nil {
		return
	}

	switch verb {
	case 'c':
		contextInfoFormat(s, e.ContextInfo(), false)

	case 'o':
		_, _ = s.Write(s2b(e.Operation()))

	case 't':
		_, _ = io.WriteString(s, e.ErrorType().String())

	case 'f':
		_, _ = io.WriteString(s, Caller(7)())

	case 's':
		if s.Flag('+') {
			// Translate в случае ошибки перевода
			// возвращает оригинальное сообщение
			_, _ = io.WriteString(s, DefaultTranslate(e))
			return
		}
		_, _ = s.Write(s2b(e.Msg()))

	case 'v':
		if s.Flag('#') {
			spew.Fdump(s, e)
			return
		}
		_ = mustMarshaler().MarshalTo(e, s)

	case 'q':
		// id
		_, _ = io.WriteString(s, "id:")
		_, _ = s.Write(s2b(e.ID()))
		// operation
		_, _ = io.WriteString(s, " operation:")
		_, _ = s.Write(s2b(e.Operation()))
		// errorType
		_, _ = io.WriteString(s, " error_type:")
		_, _ = io.WriteString(s, e.ErrorType().String())

		// errorType
		_, _ = io.WriteString(s, " context_info:")
		contextInfoFormat(s, e.ContextInfo(), false)
		// msg
		_, _ = io.WriteString(s, " message:")
		_, _ = s.Write(s2b(e.Msg()))

	case 'j':
		jmarshal := &MarshalJSON{}
		_ = jmarshal.MarshalTo(e, s)
	}
}

// дополнительные методы

// Sdump вернет текстовый дамп ошибки *Error.
func (e *Error) Sdump() string {
	if e == nil {
		return ""
	}
	return spew.Sdump(e)
}

// log

// Log выполнит логгирование ошибки с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func (e *Error) Log(l ...Logger) {
	Log(e, l...)
}

// Error methods

// Error возвращает строковое представление ошибки.
// Метод для реализации интерфейса error.
// Метод произволит перевод сообщения об ошибки, если localizer != nil.
// Для идентификации сообщения перевода используется ID ошибки.
func (e *Error) Error() string {
	data, err := e.Marshal()
	if err != nil {
		return ""
	}
	return b2s(data)
}

func (e *Error) Is(target error) bool {
	if x, ok := target.(*Error); ok { //nolint:errorlint
		return e == x
	}
	return false
}

func (e *Error) As(target interface{}) bool {
	switch x := target.(type) { //nolint:errorlint
	case **Error:
		*x = e

	default:
		return false
	}

	return true
}

func (e *Error) Unwrap() error {
	return nil
}
