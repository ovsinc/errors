package errors

import (
	"strings"

	"github.com/ovsinc/multilog"

	origerrors "errors"
)

// ErrNotValidSeverity ошибка валидации типа Severuty.
var ErrNotValidSeverity = origerrors.New("not a valid severity")

//

// ParseSeverityString парсит severity по строке.
// В случае ошибки парсинга, функция вернет SeverityUnknown и ошибку.
func ParseSeverityString(v string) (s Severity, err error) {
	switch strings.ToLower(v) {
	case "w", "warn", "warning":
		s = SeverityWarn

	case "e", "error", "err":
		s = SeverityError

	default:
		return SeverityUnknown, ErrNotValidSeverity
	}

	return s, nil
}

// ParseSeverityUint парсит severity по uint32.
// В случае ошибки парсинга, функция вернет SeverityUnknown и ошибку.
func ParseSeverityUint(v uint32) (s Severity, err error) {
	s = Severity(v)
	if !s.Valid() {
		return SeverityUnknown, ErrNotValidSeverity
	}
	return s, nil
}

//

// Severity ENUM тип определения Severity.
type Severity uint32

const (
	// SeverityUnknown не инициализированное значение, использовать не допускается.
	SeverityUnknown Severity = iota

	// SeverityWarn - предупреждение. Не является ошибкой по факту.
	SeverityWarn
	// SeverityError - ошибка.
	SeverityError

	// SeverityEnds терминирующее значение, использовать не допускается.
	SeverityEnds
)

// Uint32 конвертор в uint32
func (s Severity) Uint32() uint32 {
	return uint32(s)
}

// Valid проверка на валидность ENUM
func (s Severity) Valid() bool {
	return s > SeverityUnknown && s < SeverityEnds
}

// String получить строчное представление типа Severity.
// Для не корректных значение будет возврашено UNKNOWN.
func (s Severity) String() (str string) {
	switch s {
	case SeverityError:
		str = "ERROR"

	case SeverityWarn:
		str = "WARN"

	case SeverityEnds, SeverityUnknown:
		str = "UNKNOWN"

	default:
		str = "UNKNOWN"
	}
	return str
}

//

func customlog(l multilog.Logger, e error, severity Severity) {
	if e == nil {
		return
	}

	switch severity {
	case SeverityError:
		l.Errorf(e.Error())

	case SeverityWarn:
		l.Warnf(e.Error())

	case SeverityEnds, SeverityUnknown:
		l.Errorf(e.Error())

	default:
		l.Errorf(e.Error())
	}
}

func getLogger(l ...multilog.Logger) multilog.Logger {
	logger := multilog.DefaultLogger
	if len(l) > 0 {
		logger = l[0]
	}
	return logger
}

// хелперы

// AppendWithLog как и Append создаст или дополнит цепочку ошибок err с помощью errs,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func AppendWithLog(errs ...error) error {
	e := Append(errs...)
	Log(e)
	return e
}

// WrapWithLog обернет ошибку olderr в err и вернет цепочку,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func WrapWithLog(olderr error, err error) error {
	e := Wrap(olderr, err)
	Log(e)
	return e
}

// Log выполнить логгирование ошибки err с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func Log(err error, l ...multilog.Logger) {
	severity := SeverityError

	if errseverity, ok := simpleCast(err); ok {
		severity = errseverity.Severity()
	}
	customlog(getLogger(l...), err, severity)
}

// NewWithLog конструктор *Error, как и New,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func NewWithLog(msg string, ops ...Options) *Error {
	e := New(msg, ops...)
	e.Log()
	return e
}

// дополнительные методы *Error

// Log выполнит логгирование ошибки с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func (e *Error) Log(l ...multilog.Logger) {
	customlog(getLogger(l...), e, e.Severity())
}
