package errors

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/valyala/bytebufferpool"
	"gitlab.com/ovsinc/errors/log"
	logcommon "gitlab.com/ovsinc/errors/log/common"
)

var (
	_ interface{ Error() string } = (*Error)(nil)
	_ Errorer                     = (*Error)(nil)
)

// // _bufferPool is a pool of bytes.Buffers.
// var _bufferPool = sync.Pool{
// 	New: func() interface{} {
// 		return &bytes.Buffer{}
// 	},
// }

// Errorer итерфейс кастомной ошибки.
type Errorer interface {
	WithOptions(ops ...Options) Errorer
	ID() string
	Severity() log.Severity
	Msg() string
	Error() string
	Sdump() string
	ErrorOrNil() error

	ErrorType() string

	Operations() []Operation
	Format(s fmt.State, verb rune)

	TranslateContext() *TranslateContext
	Localizer() *i18n.Localizer

	Log(l ...logcommon.Logger)
}

// Error структура кастомной ошибки.
// Внимание. Это НЕ потоко-безопасный объект.
type Error struct {
	severity         log.Severity
	operations       []Operation
	formatFn         FormatFn
	contextInfo      CtxMap
	translateContext *TranslateContext
	localizer        *i18n.Localizer
	errorType        string
	msg              string
	id               string
}

// New конструктор на необязательных параметрах
// * ops ...Options -- параметризация через функции-парметры
//
// ** *Error
func New(msg string, ops ...Options) Errorer {
	e := &Error{
		operations: []Operation{},
		severity:   log.SeverityError,
		errorType:  UnknownErrorType,
		msg:        msg,
	}
	for _, op := range ops {
		op(e)
	}
	return e
}

// setters

// WithOptions производит параметризацию *Error с помощью функции-парметры Options.
// Допускается указывать произвольно количество ops.
// Возвращается модифицированный экземпляр *Error.
func (e *Error) WithOptions(ops ...Options) Errorer {
	for _, op := range ops {
		op(e)
	}
	return e
}

// getters

// ID возвращает ID ошибки.
func (e *Error) ID() string {
	return e.id
}

// Severity возвращает критичность ошибки
func (e *Error) Severity() log.Severity {
	return e.severity
}

// Msg возвращает исходное сообщение об ошибке
func (e *Error) Msg() string {
	return e.msg
}

// Error возвращает строковое представление ошибки.
// Метод для реализации интерфейса error.
// Метод произволит перевод сообщения об ошибки, если localizer != nil.
// Для идентификации сообщения перевода используется ID ошибки.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	fn := e.formatFn
	if e.formatFn == nil {
		fn = defaultFormatFn
	}

	// buf := _bufferPool.Get().(*bytes.Buffer)
	// buf.Reset()
	// defer _bufferPool.Put(buf)

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	fn(buf, e)

	return buf.String()
}

// дополнительные методы

// Sdump вернет текстовый дамп ошибки *Error.
func (e *Error) Sdump() string {
	if e == nil {
		return ""
	}
	return spew.Sdump(e)
}

// ErrorOrNil вернет ошибку или nil.
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
