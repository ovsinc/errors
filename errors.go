package errors

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/valyala/bytebufferpool"
	"gitlab.com/ovsinc/errors/log"
	logcommon "gitlab.com/ovsinc/errors/log/common"
)

var (
	_ interface{ Error() string } = (*Error)(nil)
	_ error                       = (*Error)(nil)
	_ errorer                     = (*Error)(nil)
)

type errorer interface {
	WithOptions(ops ...Options) *Error
	ID() string
	Severity() log.Severity
	Msg() string
	Error() string
	Sdump() string
	ErrorOrNil() error
	ContextInfo() CtxMap

	ErrorType() string

	Operations() []string
	Format(s fmt.State, verb rune)

	TranslateContext() *TranslateContext
	Localizer() *i18n.Localizer
	WriteTranslateMsg(w io.Writer) (int, error)
	TranslateMsg() string

	Log(l ...logcommon.Logger)
}

// Error структура кастомной ошибки.
// Это потоко-безопасный объект.
type Error struct {
	severity         log.Severity
	operations       []string
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
func New(msg string, ops ...Options) *Error {
	e := &Error{
		severity: log.SeverityError,
		msg:      msg,
	}
	for _, op := range ops {
		op(e)
	}
	return e
}

// setters

// WithOptions производит параметризацию *Error с помощью функции-парметры Options.
// Допускается указывать произвольно количество ops.
// Возвращается новый экземпляр *Error.
func (e *Error) WithOptions(ops ...Options) *Error {
	if e == nil {
		return nil
	}

	newerr := &Error{
		severity:         e.severity,
		operations:       e.operations,
		formatFn:         e.formatFn,
		contextInfo:      e.contextInfo,
		translateContext: e.translateContext,
		localizer:        e.localizer,
		errorType:        e.errorType,
		msg:              e.msg,
		id:               e.id,
	}

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

// Severity возвращает критичность ошибки
func (e *Error) Severity() log.Severity {
	if e == nil {
		return 0
	}

	return e.severity
}

// Msg возвращает исходное сообщение об ошибке
func (e *Error) Msg() string {
	if e == nil {
		return ""
	}

	return e.msg
}

// ErrorType вернет тип ошибки
func (e *Error) ErrorType() string {
	if e == nil {
		return ""
	}

	return e.errorType
}

// Error возвращает строковое представление ошибки.
// Метод для реализации интерфейса error.
// Метод произволит перевод сообщения об ошибки, если localizer != nil.
// Для идентификации сообщения перевода используется ID ошибки.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	var fn FormatFn
	switch {
	case e.formatFn != nil:
		fn = e.formatFn
	case DefaultFormatFn != nil:
		fn = DefaultFormatFn
	default:
		fn = StringFormat
	}

	// buf := _bufferPool.Get().(*bytes.Buffer)
	// buf.Reset()
	// defer _bufferPool.Put(buf)

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	fn(buf, e)

	return buf.String()
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
		fmt.Fprintf(s, "%v\n", e.Operations())

	case 'l':
		_, _ = io.WriteString(s, e.Severity().String())

	case 't':
		_, _ = io.WriteString(s, e.ErrorType())

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

// ErrorOrNil вернет ошибку или nil.
// ошибкой считается *Error != nil и Severity == SeverityError
// т.е. SeverityWarn НЕ ошибка
func (e *Error) ErrorOrNil() error {
	if e != nil && e.Severity() == log.SeverityError {
		return e
	}
	return nil
}
