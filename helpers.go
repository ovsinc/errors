package errors

import (
	origerrors "errors"

	multierror "github.com/hashicorp/go-multierror"
)

// Is сообщает, соответствует ли ошибка err target-ошибке.
// Для multierr будет производится поиск в цепочке.
func Is(err, target error) bool {
	if err == nil {
		return target == nil
	}
	return origerrors.Is(err, target)
}

// As обнаруживает ошибку err, соответствующую target и устанавливает target в найденное значение.
func As(err error, target interface{}) bool {
	return origerrors.As(err, target)
}

// GetErrorType возвращает тип ошибки. Для НЕ *Error всегда будет UnknownErrorType.
func GetErrorType(err error) string {
	var errtype *Error

	if origerrors.As(err, &errtype) {
		return errtype.ErrorType()
	}

	return UnknownErrorType
}

// GetID возвращает ID ошибки. Для НЕ *Error всегда будет "".
func GetID(err error) (id string) {
	var idtype *Error

	if origerrors.As(err, &idtype) {
		return idtype.ID()
	}

	return
}

// ErrorOrNil вернет ошибку или nil
// Возможна обработка multierror или одиночной ошибки (*Error, error).
// Если хотя бы одна ошибка в цепочке является ошибкой, то она будет возвращена в качестве результата.
// В противном случае будет возвращен nil.
// Важно: *Error c Severity Warn не является ошибкой.
func ErrorOrNil(err error) error {
	var (
		multerr *multierror.Error
		myerr   *Error
	)

	switch {
	case origerrors.As(err, &myerr):
		return myerr.ErrorOrNil()

	case origerrors.As(err, &multerr):
		for _, e := range multerr.WrappedErrors() {
			// если это *Error и хотя бы одна из *Error не nil, вернуть ее
			if origerrors.As(e, &myerr) && myerr.ErrorOrNil() != nil {
				return myerr
			}
		}
		return nil
	}

	return err
}

// Cast преобразует тип error в *Error
// Если error не соответствует *Error, то будет создан *Error с сообщением err.Error().
// Для err == nil, вернется nil.
func Cast(err error) Errorer {
	if err == nil {
		return nil
	}

	var myerr *Error

	if origerrors.As(err, &myerr) {
		return myerr
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

	if unwrap, ok := err.(unwraper); ok { //nolint:errorlint
		return unwrap.Unwrap()
	}

	return nil
}

// UnwrapByID вернет ошибку (*Error) с указанным ID.
// Для multierror, функция вернет ошибку с указанным ID.
// Если ошибка с указанным ID не найдена, вернется nil.
func UnwrapByID(err error, id string) Errorer {
	getiderrFn := func(err error) (Errorer, bool) {
		var iderr *Error
		if origerrors.As(err, &iderr) {
			return iderr, iderr.ID() == id
		}
		return nil, false
	}

	var multerr *multierror.Error

	if origerrors.As(err, &multerr) {
		for _, e := range multerr.WrappedErrors() {
			if iderr, ok := getiderrFn(e); ok {
				return iderr
			}
		}
	}

	if iderr, ok := getiderrFn(err); ok {
		return iderr
	}

	return nil
}
