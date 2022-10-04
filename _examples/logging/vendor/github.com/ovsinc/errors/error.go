package errors

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
)

const verboseStr = "id:%s operation:%s errorType:%s contextInfo:%v msg:%s"

var (
	ErrUnknownMarshaller = New("marshaller not found")

	_ error         = (*Error)(nil)
	_ fmt.Formatter = (*Error)(nil)
)

// CtxMap map контекста ошибки.
// В качестве ключа всегда должна быть строка, а значение - любой тип.
// При преобразовании ошибки в строку CtxMap может использоваться различные методы.
// Для функции JSONFormat CtxMap будет преобразовываться с помощью JSON marshall.
// Для функции StringFormat CtxMap будет преобразовываться с помощью fmt.Sprintf.
type CtxMap map[string]interface{}

// New конструктор на необязательных параметрах
// * ops ...Options -- параметризация через функции-парметры.
// См. options.go
//
// ** *Error
func New(i interface{}) *Error {
	var msg string
	switch t := i.(type) {
	case string:
		msg = t
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

func NewWith(ops ...Options) *Error {
	e := Error{}
	for _, op := range ops {
		op(&e)
	}
	return &e
}

func NewWithLog(ops ...Options) *Error {
	e := NewWith(ops...)
	e.Log()
	return e
}

// Error структура кастомной ошибки.
// Это потоко-безопасный объект.
type Error struct {
	id, msg, operation, errorType []byte
	// type like a:
	// http - https://cs.opensource.google/go/go/+/refs/tags/go1.19.1:src/net/http/status.go;l=9
	// grpc - https://pkg.go.dev/google.golang.org/grpc/codes
	contextInfo CtxMap
}

// WithOptions производит параметризацию *Error с помощью функции-парметры Options.
// Допускается указывать произвольно количество ops.
// Возвращается новый экземпляр *Error.
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
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) ID() []byte {
	if e == nil {
		return nil
	}
	return e.id
}

// Msg возвращает исходное сообщение об ошибке.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) Msg() []byte {
	if e == nil {
		return nil
	}
	return e.msg
}

// Operations вернет список операций.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) Operation() []byte {
	if e == nil {
		return nil
	}
	return e.operation
}

// ErrorType вернет тип ошибки.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) ErrorType() []byte {
	if e == nil {
		return nil
	}

	return e.errorType
}

// ContextInfo вернет контекст CtxMap ошибки.
func (e *Error) ContextInfo() CtxMap {
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

func (e *Error) Marshal(fn ...Marshaller) ([]byte, error) {
	marshal := mustMarshaler(fn...)
	return marshal.Marshal(e)
}

// Format производит форматирование строки, для поддержки fmt.Printf().
func (e *Error) Format(s fmt.State, verb rune) {
	if e == nil {
		return
	}

	switch verb {
	case 'c':
		ctxi := e.ContextInfo()
		contextInfoFormat(s, &ctxi, false)

	case 'o':
		_, _ = s.Write(e.Operation())

	case 't':
		_, _ = s.Write(e.ErrorType())

	case 'f':
		_, _ = io.WriteString(s, Caller(7)())

	case 's':
		if s.Flag('+') {
			// Translate в случае ошибки перевода
			// возвращает оригинальное сообщение
			io.WriteString(s, DefaultTranslate(e))
			return
		}
		_, _ = s.Write(e.Msg())

	case 'v':
		if s.Flag('#') {
			spew.Fdump(s, e)
			return
		}
		mustMarshaler().MarshalTo(e, s)

	case 'q':
		fmt.Fprintf(
			s,
			verboseStr,
			e.ID(),
			e.Operation(),
			e.ErrorType(),
			e.ContextInfo(),
			e.Msg(),
		)

	case 'j':
		jmarshal := &MarshalJSON{}
		jmarshal.MarshalTo(e, s)
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
	return string(data)
}

func (e *Error) Is(target error) bool {
	switch x := target.(type) { //nolint:errorlint
	case *Error:
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
