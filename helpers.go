package errors

import (
	"bytes"
	origerrors "errors"
)

// GetID возвращает ID ошибки. Для НЕ *Error всегда будет "".
func GetID(err error) (id string) {
	if e, ok := simpleCast(err); ok { //nolint:errorlint
		return e.ID().String()
	}
	return
}

func findByID(err error, id []byte) (*Error, bool) {
	switch t := err.(type) { //nolint:errorlint
	case Multierror: // multiError
		for _, e := range t.Errors() {
			if bytes.Equal(e.ID().Bytes(), id) {
				return e, true
			}
		}
		return nil, false

	case *Error:
		return t, bytes.Equal(t.ID().Bytes(), []byte(id))
	}

	return nil, false
}

// UnwrapByID вернет ошибку (*Error) с указанным ID.
// Если ошибка с указанным ID не найдена, вернется nil.
func UnwrapByID(err error, id string) *Error {
	if e, ok := findByID(err, []byte(id)); ok {
		return e
	}
	return nil
}

// Contains проверит есть ли в цепочке ошибка с указанным ID.
// Допускается в качестве аргумента err указывать одиночную ошибку.
func Contains(err error, id string) bool {
	_, ok := findByID(err, []byte(id))
	return ok
}

// Is сообщает, соответствует ли ошибка err target-ошибке.
// Для multierr будет производится поиск в цепочке.
func Is(err, target error) bool {
	if err == nil {
		return target == nil
	}
	return origerrors.Is(err, target)
}

// As обнаруживает ошибку err, соответствующую типу target и устанавливает target в найденное значение.
func As(err error, target interface{}) bool {
	return origerrors.As(err, target)
}

func Unwrap(err error) error {
	return origerrors.Unwrap(err)
}
