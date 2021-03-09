package log

import (
	origerrors "errors"
	"strings"
)

// ErrNotValidSeverity ошибка валидации типа Severuty
var ErrNotValidSeverity = origerrors.New("not a valid severity")

//

// ParseSeverityString парсить severity по строке
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

// ParseSeverityUint парсить severity по uint32
func ParseSeverityUint(v uint32) (s Severity, err error) {
	s = Severity(v)
	if !s.Valid() {
		return SeverityUnknown, ErrNotValidSeverity
	}
	return s, nil
}

//

// Severity ENUM тип определения Severity
type Severity uint32

const (
	// SeverityUnknown не инициализированное значение, использовать не допускается
	SeverityUnknown Severity = iota

	// SeverityWarn ...
	SeverityWarn
	// SeverityError ...
	SeverityError

	// SeverityEnds терминирующее значение, использовать не допускается
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

// String получить строчное предстваление типа Severity
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
