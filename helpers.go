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

func Find(err error, fn func(error) bool) error {
	switch t := err.(type) { //nolint:errorlint
	case *multiError:
		cure := t.errors[t.cur.Load()]
		if fn(cure) {
			return cure
		}
		for i, e := range t.errors {
			if fn(e) {
				t.cur.Store(int32(i))
				return e
			}
		}

	case *Error:
		if fn(t) {
			return t
		}

	case error:
		if fn(t) {
			return t
		}
	}

	return nil
}

// FindByID вернет ошибку (*Error) с указанным ID.
// Если ошибка с указанным ID не найдена, вернется nil.
func FindByID(err error, id string) error {
	return Find(err, func(e error) bool {
		ee, ok := e.(*Error) //nolint:errorlint
		return ok && bytes.Equal(ee.ID(), []byte(id))
	})
}

// FindByErr вернет ошибку (*Error) соответсвующую target или nil.
// Если ошибка не найдена, вернется nil.
func FindByErr(err error, target error) error {
	return Find(err, func(e error) bool {
		return origerrors.Is(e, target)
	})
}

// Contains проверит есть ли в цепочке целевая ошибка.
// Допускается в качестве аргумента err указывать одиночную ошибку.
func Contains(err error, fn func(error) bool) bool {
	switch t := err.(type) { //nolint:errorlint
	case *multiError:
		cure := t.errors[t.cur.Load()]
		if fn(cure) {
			return true
		}
		for i, e := range t.errors {
			if fn(e) {
				t.cur.Store(int32(i))
				return true
			}
		}
	}
	return fn(err)
}

// Contains проверит есть ли в цепочке ошибка с указанным ID.
// Допускается в качестве аргумента err указывать одиночную ошибку.
func ContainsByID(err error, id string) bool {
	return Contains(err, func(e error) bool {
		ee, ok := e.(*Error) //nolint:errorlint
		return ok && bytes.Equal(ee.ID(), []byte(id))
	})
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
