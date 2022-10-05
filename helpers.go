package errors

import (
	"bytes"
	origerrors "errors"
)

// GetID возвращает ID ошибки. Для НЕ *Error всегда будет "".
func GetID(err error) (id string) {
	if e, ok := err.(*Error); ok { //nolint:errorlint
		return string(e.ID())
	}
	return
}

func findByErr(err error, target error) (error, bool) {
	switch t := err.(type) { //nolint:errorlint
	case Multierror:
		for _, e := range t.Errors() {
			if Is(e, target) {
				return e, true
			}
		}
		return nil, false

	case *Error:
		if Is(t, target) {
			return t, true
		}

	case error:
		if Is(t, target) {
			return New(t), true
		}
	}

	return nil, false
}

func findByID(err error, id []byte) (error, bool) {
	switch t := err.(type) { //nolint:errorlint
	case interface{ FindByID([]byte) (error, bool) }: // multiError
		return t.FindByID(id)

	case *Error:
		return t, bytes.Equal(t.ID(), []byte(id))
	}

	return nil, false
}

// UnwrapByID вернет ошибку (*Error) с указанным ID.
// Если ошибка с указанным ID не найдена, вернется nil.
func UnwrapByID(err error, id string) error {
	if e, ok := findByID(err, []byte(id)); ok {
		return e
	}
	return nil
}

// UnwrapByErr вернет ошибку (*Error) соответсвующую target или nil.
// Если ошибка не найдена, вернется nil.
func UnwrapByErr(err error, target error) error {
	if e, ok := findByErr(err, target); ok {
		return e
	}
	return nil
}

// Contains проверит есть ли в цепочке целевая ошибка.
// Допускается в качестве аргумента err указывать одиночную ошибку.
func Contains(err error, target error) bool {
	_, ok := findByErr(err, target)
	return ok
}

// Contains проверит есть ли в цепочке ошибка с указанным ID.
// Допускается в качестве аргумента err указывать одиночную ошибку.
func ContainsByID(err error, id string) bool {
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

// Cast преобразует тип error в *Error
// Если error не соответствует *Error, то будет создан *Error с сообщением err.Error().
// Для err == nil, вернется nil.
func Cast(err error) (*Error, bool) {
	switch t := err.(type) { //nolint:errorlint
	case nil:
		return nil, false

	case *Error:
		return t, true

	case interface{ Last() error }:
		return New(t.Last()), true

	default:
		return New(err), true
	}
}

// CastMultierr преобразует тип error в *Multierror
// Если error не соответствует Multierror, то будет создан Multierror с сообщением err.Error().
// Для err == nil, вернется nil.
func CastMultierr(err error) (Multierror, bool) {
	switch t := err.(type) { //nolint:errorlint
	case nil:
		return nil, false

	case Multierror:
		return t, true

	default:
		return Combine(err).(Multierror), true ////nolint:errorlint
	}
}
