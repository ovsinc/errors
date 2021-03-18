package errors

import (
	origerrors "errors"
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

// GetErrorType возвращает тип ошибки. Для НЕ *Error всегда будет "".
func GetErrorType(err error) string {
	var errtype *Error

	if origerrors.As(err, &errtype) {
		return errtype.ErrorType()
	}

	return ""
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
		multerr *multiError
		myerr   *Error
	)

	switch {
	case origerrors.As(err, &myerr):
		return myerr.ErrorOrNil()

	case origerrors.As(err, &multerr):
		for _, e := range multerr.Errors() {
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

func findByID(err error, id string) (Errorer, bool) {
	var (
		merr  *multiError
		myerr *Error
	)

	getiderrFn := func(err error) (Errorer, bool) {
		var iderr *Error
		if origerrors.As(err, &iderr) {
			return iderr, iderr.ID() == id
		}
		return nil, false
	}

	switch {
	case origerrors.As(err, &merr):
		for _, e := range merr.Errors() {
			if iderr, ok := getiderrFn(e); ok {
				return iderr, true
			}
		}

	case origerrors.As(err, &myerr):
		if myerr.ID() == id {
			return myerr, true
		}
	}

	return nil, false
}

// UnwrapByID вернет ошибку (*Error) с указанным ID.
// Если ошибка с указанным ID не найдена, вернется nil.
func UnwrapByID(err error, id string) Errorer {
	if e, ok := findByID(err, id); ok {
		return e
	}
	return nil
}

// Contains проверит есть ли в цепочке ошибка с указанным ID.
// Допускается в качестве аргумента err указывать одиночную ошибку.
func Contains(err error, id string) bool {
	_, ok := findByID(err, id)
	return ok
}
