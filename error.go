package errors

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/ovsinc/multilog"
)

var (
	_ error   = (*Error)(nil)
	_ errorer = (*Error)(nil)
)

var ErrUnknownMarshaller = New("marshaller not found")

type errorer interface {
	WithOptions(ops ...Options) *Error

	ID() Object
	Msg() Object
	ContextInfo() CtxMap
	Operation() Object

	Sdump() string

	Format(s fmt.State, verb rune)
	Marshal(fn ...Marshaller) ([]byte, error)

	Error() string

	TranslateMsg() string
	WriteTranslateMsg(w io.Writer) (int, error)

	Log(l ...multilog.Logger)
}

// CtxMap map контекста ошибки.
// В качестве ключа всегда должна быть строка, а значение - любой тип.
// При преобразовании ошибки в строку CtxMap может использоваться различные методы.
// Для функции JSONFormat CtxMap будет преобразовываться с помощью JSON marshall.
// Для функции StringFormat CtxMap будет преобразовываться с помощью fmt.Sprintf.
type CtxMap map[string]interface{}

// Error структура кастомной ошибки.
// Это потоко-безопасный объект.
type Error struct {
	id               Object
	msg              Object
	operation        Object
	translateContext *TranslateContext
	localizer        Localizer
	contextInfo      CtxMap
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
func (e *Error) ID() Object {
	if e == nil || e.id == nil {
		return NewObjectEmpty()
	}

	return e.id
}

// Msg возвращает исходное сообщение об ошибке.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) Msg() Object {
	if e == nil || e.msg == nil {
		return NewObjectEmpty()
	}

	return e.msg
}

// Operations вернет список операций.
// Это безопасный метод, всегда возвращает не nil.
func (e *Error) Operation() Object {
	if e == nil || e.operation == nil {
		return NewObjectEmpty()
	}
	return e.operation
}

// TranslateContext вернет *TranslateContext.
func (e *Error) TranslateContext() *TranslateContext {
	return e.translateContext
}

// методы форматирования

func (e *Error) Marshal(fn ...Marshaller) ([]byte, error) {
	var marshal Marshaller
	switch {
	case len(fn) > 0:
		marshal = fn[0]
	case DefaultMarshaller != nil:
		marshal = DefaultMarshaller
	default:
		return nil, ErrUnknownMarshaller
	}
	return marshal.Marshal(e)
}

// Format производит форматирование строки, для поддержки fmt.Printf().
func (e *Error) Format(s fmt.State, verb rune) {
	if e == nil {
		return
	}

	switch verb {
	case 'c':
		fmt.Fprintf(s, "%v\n", e.ContextInfo())

	case 'o':
		fmt.Fprintf(s, "%v\n", e.Operation())

	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.Sdump())
			return
		}
		_, _ = io.WriteString(s, e.Error())

	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

// context info

// ContextInfo вернет контекст CtxMap ошибки.
func (e *Error) ContextInfo() CtxMap {
	return e.contextInfo
}

// дополнительные методы

// Sdump вернет текстовый дамп ошибки *Error.
func (e *Error) Sdump() string {
	if e == nil {
		return ""
	}

	if e == nil {
		return ""
	}
	return spew.Sdump(e)
}

// log

// Log выполнит логгирование ошибки с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func (e *Error) Log(l ...multilog.Logger) {
	logger := getLogger(l...)
	if logger == nil {
		return
	}
	logger.Errorf(e.Error())
}

// Error methods

// Error возвращает строковое представление ошибки.
// Метод для реализации интерфейса error.
// Метод произволит перевод сообщения об ошибки, если localizer != nil.
// Для идентификации сообщения перевода используется ID ошибки.
func (e *Error) Error() string {
	marshal := &MarshalString{}
	data, _ := marshal.Marshal(e)
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
