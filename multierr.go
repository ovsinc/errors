// Используется оригинальный код проекта "go.uber.org/multierr" с частичным заимствованием.
// Код проекта "go.uber.org/multierr" распространяется под лицензией MIT (https://github.com/uber-go/multierr/blob/master/LICENSE.txt).

package errors

import (
	"fmt"
	"io"

	"github.com/valyala/bytebufferpool"
)

// Errors returns a slice containing zero or more errors that the supplied
// error is composed of. If the error is nil, a nil slice is returned.
//
// 	err := multierr.Append(r.Close(), w.Close())
// 	errors := multierr.Errors(err)
//
// If the error is not composed of other errors, the returned slice contains
// just the error that was passed in.
//
// Callers of this function are free to modify the returned slice.
func Errors(err error) []*Error {
	if eg, ok := err.(multiError); ok { //nolint:errorlint
		return eg
	}

	return appendError(make([]*Error, 0, 1), err)
}

var _ Multierror = (multiError)(nil)

type Multierror interface {
	Errors() []*Error
	Error() string
	Format(f fmt.State, c rune)
}

// multiError is an error that holds one or more errors.
//
// An instance of this is guaranteed to be non-empty and flattened. That is,
// none of the errors inside multiError are other multiErrors.
//
// multiError formats to a semi-colon delimited list of error messages with
// %v and with a more readable multi-line format with %+v.
type multiError []*Error

// Errors returns the list of underlying errors.
//
// This slice MUST NOT be modified.
func (merr multiError) Errors() []*Error {
	return merr
}

func (merr multiError) Error() string {
	if len(merr) == 0 {
		return ""
	}

	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)

	merr.writeLines(buff)

	return buff.String()
}

func (merr multiError) Format(f fmt.State, c rune) {
	switch c {
	case 'w', 'v', 's':
		merr.writeLines(f)
	case 'j':
		JSONMultierrFuncFormat(f, merr)
	}
}

func (merr multiError) writeLines(w io.Writer) {
	if DefaultMultierrFormatFunc == nil {
		StringMultierrFormatFunc(w, merr)
		return
	}
	DefaultMultierrFormatFunc(w, merr)
}

type inspectResult struct {
	// Number of top-level non-nil errors
	Count int

	// Total number of errors including multiErrors
	Capacity int
}

// Inspects the given slice of errors so that we can efficiently allocate
// space for it.
func inspect(errors []error) (res inspectResult) {
	for _, err := range errors {
		if err == nil {
			continue
		}

		res.Count++
		if merr, ok := err.(multiError); ok { //nolint:errorlint
			res.Capacity += len(merr)
		} else {
			res.Capacity++
		}
	}
	return
}

func appendError(errors multiError, err error) []*Error {
	switch t := err.(type) {
	case nil:
		return nil

	case *Error:
		return append(errors, t)

	case Multierror:
		return append(errors, t.Errors()...)

	case error:
		return append(errors, New(t.Error()))
	}

	return errors
}

// fromSlice converts the given list of errors into a single error.
func fromSlice(errors []error) Multierror {
	res := inspect(errors)
	if res.Count == 0 {
		return nil
	}

	nonNilErrs := make(multiError, 0, res.Capacity)
	for _, err := range errors {
		if err == nil {
			continue
		}
		nonNilErrs = appendError(nonNilErrs, err)
	}

	return nonNilErrs
}

// Combine создаст цепочку ошибок из ошибок ...errors.
// Допускается использование `nil` в аргументах.
func Combine(errors ...error) error {
	return fromSlice(errors)
}

// Wrap обернет ошибку `left` ошибкой `right`, получив цепочку.
// Допускается использование `nil` в одном из аргументов.
func Wrap(left error, right error) error {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	}
	return fromSlice([]error{left, right})
}

// Unwrap вернет самую новую ошибку в стеке
func (merr multiError) Unwrap() error {
	if len(merr) == 0 {
		return nil
	}
	return merr[len(merr)-1]
}

func (merr multiError) As(target interface{}) bool {
	if x, ok := target.(*Multierror); ok { //nolint:errorlint
		*x = merr
		return true
	}
	return false
}

func (merr multiError) Is(target error) bool {
	if x, ok := target.(Multierror); ok { //nolint:errorlint
		return x == &merr
	}
	return false
}
