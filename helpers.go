package errors

// GetErrorType возвращает тип ошибки. Для НЕ *Error всегда будет "".
func GetErrorType(err error) string {
	if errtype, ok := err.(*Error); ok { //nolint:errorlint
		return errtype.ErrorType()
	}

	return ""
}

// GetID возвращает ID ошибки. Для НЕ *Error всегда будет "".
func GetID(err error) (id string) {
	if idtype, ok := err.(*Error); ok { //nolint:errorlint
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
	if e, ok := cast(err); ok {
		return e.ErrorOrNil()
	}

	return err
}

func errsFn(errs []error) Errorer {
	for _, e := range errs {
		if myerr, ok := simpleCast(e); ok {
			return myerr
		}
	}
	return nil
}

func simpleCast(err error) (Errorer, bool) {
	e, ok := err.(Errorer) //nolint:errorlint
	return e, ok
}

func cast(err error) (Errorer, bool) {
	switch t := err.(type) { //nolint:errorlint
	case interface{ Errors() []error }: // *multiError
		return errsFn(t.Errors()), true

	case interface{ WrappedErrors() []error }: // *github.com/hashicorp/go-multierror.Error
		return errsFn(t.WrappedErrors()), true

	case Errorer:
		return t, true
	}

	return nil, false
}

// Cast преобразует тип error в *Error
// Если error не соответствует *Error, то будет создан *Error с сообщением err.Error().
// Для err == nil, вернется nil.
func Cast(err error) Errorer {
	if err == nil {
		return nil
	}

	if e, ok := cast(err); ok {
		return e
	}

	return New(err.Error())
}

func findByID(err error, id string) (Errorer, bool) {
	checkIDFn := func(errs []error) Errorer {
		for _, err := range errs {
			if e, ok := simpleCast(err); ok && e.ID() == id {
				return e
			}
		}
		return nil
	}

	switch t := err.(type) { //nolint:errorlint
	case interface{ Errors() []error }: // *multiError
		e := checkIDFn(t.Errors())
		return e, e != nil

	case interface{ WrappedErrors() []error }: // *github.com/hashicorp/go-multierror.Error
		e := checkIDFn(t.WrappedErrors())
		return e, e != nil

	case Errorer:
		return t, t.ID() == id
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
