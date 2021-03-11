package errors

import (
	origerrors "errors"

	multierror "github.com/hashicorp/go-multierror"
)

// Is сообщает, соответствует ли ошибка err target-ошибке.
// Для multierr будет производится поиск в цепочке.
func Is(err, target error) bool {
	return origerrors.Is(err, target)
}

// As обнаруживает ошибку err, соответствующую target и устанавливает target в найденное значение.
func As(err error, target interface{}) bool {
	return origerrors.As(err, target)
}

// GetErrorType возвращает тип ошибки. Для НЕ *Error всегда будет UnknownErrorType.
func GetErrorType(err error) ErrorType {
	type errtyper interface {
		ErrorType() ErrorType
	}
	if etype, ok := err.(errtyper); ok {
		return etype.ErrorType()
	}
	return UnknownErrorType
}

// GetID возвращает ID ошибки. Для НЕ *Error всегда будет "".
func GetID(err error) (id string) {
	type idtyper interface {
		ID() string
	}
	if customerr, ok := err.(idtyper); ok {
		id = customerr.ID()
	}
	return
}

// ErrorOrNil вернет ошибку или nil
// Возможна обработка multierror или одиночной ошибки (*Error, error).
// Если хотя бы одна ошибка в цепочке является ошибкой, то она будет возвращена в качестве результата.
// В противном случае будет возвращен nil.
// Важно: *Error c Severity Warn не является ошибкой.
func ErrorOrNil(err error) error {
	switch t := err.(type) { // nolint:errorlint
	case *multierror.Error:
		for _, e := range t.WrappedErrors() {
			// если это *Error и хотя бы одна из *Error не nil, вернуть ее
			if customerr, ok := e.(*Error); ok && customerr.ErrorOrNil() != nil {
				return customerr
			}
		}
		return nil

	case *Error:
		return t.ErrorOrNil()

	}

	return err
}

// Cast преобразует тип error в *Error
// Если error не соответствует *Error, то будет создан *Error с сообщением err.Error().
// Для err == nil, вернется nil.
func Cast(err error) *Error {
	if err == nil {
		return nil
	}

	if t, ok := err.(*Error); ok { // nolint:errorlint
		return t
	}

	return New(err.Error())
}

// Unwrap позволяет получить оригинальную ошибку.
// Для этого тип err должен иметь метод `Unwrap() error`.
// Для multierror будет разернута цепочка ошибок.
// В противном случае будет возвращен nil.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	type unwraper interface {
		Unwrap() error
	}

	if unwrap, ok := err.(unwraper); ok {
		return unwrap.Unwrap()
	}
	return nil
}

// UnwrapWithID вернет ошибку с указанным ID.
// Для multierror, функция вернет ошибку с указанным ID.
// Если ошибка с указанным ID не найдена, вернется nil.
func UnwrapWithID(err error, id string) error {
	type idtyper interface {
		ID() string
	}

	getiderrFn := func(err error) (error, bool) {
		iderr, ok := err.(idtyper)
		if !ok {
			return nil, false
		}
		return err, iderr.ID() == id
	}

	if t, ok := err.(*multierror.Error); ok {
		for _, e := range t.WrappedErrors() {
			if _, ok := getiderrFn(e); ok {
				return e
			}
		}
	}

	if e, ok := getiderrFn(err); ok {
		return e
	}

	return nil
}
