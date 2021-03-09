package errors

import (
	origerrors "errors"

	multierror "github.com/hashicorp/go-multierror"
)

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return origerrors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool {
	return origerrors.As(err, target)
}

func GetErrorType(err error) ErrorType {
	type errtyper interface {
		ErrorType() ErrorType
	}
	if etype, ok := err.(errtyper); ok {
		return etype.ErrorType()
	}
	return UnknownErrorType
}

func GetID(err error) (id string) {
	type idtyper interface {
		ID() string
	}
	if customerr, ok := err.(idtyper); ok {
		id = customerr.ID()
	}
	return
}

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

func Cast(err error) *Error {
	switch t := err.(type) { // nolint:errorlint
	case *Error:
		return t
	}
	return New(err.Error())
}
