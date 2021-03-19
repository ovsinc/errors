// Используется оригинальный код проекта "go.uber.org/multierr" с частичным заимствованием.
// Код проекта "go.uber.org/multierr" распространяется под лицензией MIT (https://github.com/uber-go/multierr/blob/master/LICENSE.txt).

package errors

import (
	"fmt"
	"io"

	"github.com/valyala/bytebufferpool"
)

type errorGroup interface {
	Errors() []error
}

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
func Errors(err error) []error {
	if err == nil {
		return nil
	}

	// Note that we're casting to multiError, not errorGroup. Our contract is
	// that returned errors MAY implement errorGroup. Errors, however, only
	// has special behavior for multierr-specific error objects.
	//
	// This behavior can be expanded in the future but I think it's prudent to
	// start with as little as possible in terms of contract and possibility
	// of misuse.
	if eg, ok := err.(*multiError); ok { //nolint:errorlint
		errors := eg.Errors()
		result := make([]error, len(errors))
		copy(result, errors)
		return result
	}

	return []error{err}
}

// multiError is an error that holds one or more errors.
//
// An instance of this is guaranteed to be non-empty and flattened. That is,
// none of the errors inside multiError are other multiErrors.
//
// multiError formats to a semi-colon delimited list of error messages with
// %v and with a more readable multi-line format with %+v.
type multiError struct {
	errors []error
}

var _ errorGroup = (*multiError)(nil)

// Errors returns the list of underlying errors.
//
// This slice MUST NOT be modified.
func (merr *multiError) Errors() []error {
	if merr == nil {
		return nil
	}
	return merr.errors
}

func (merr *multiError) Error() string {
	if merr == nil {
		return ""
	}

	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)

	merr.writeLines(buff)

	return buff.String()
}

func (merr *multiError) Format(f fmt.State, c rune) {
	merr.writeLines(f)
}

func (merr *multiError) writeLines(w io.Writer) {
	if DefaultMultierrFormatFunc == nil {
		StringMultierrFormatFunc(w, merr.errors)
		return
	}
	DefaultMultierrFormatFunc(w, merr.errors)
}

type inspectResult struct {
	// Number of top-level non-nil errors
	Count int

	// Total number of errors including multiErrors
	Capacity int

	// Index of the first non-nil error in the list. Value is meaningless if
	// Count is zero.
	FirstErrorIdx int

	// Whether the list contains at least one multiError
	ContainsMultiError bool
}

// Inspects the given slice of errors so that we can efficiently allocate
// space for it.
func inspect(errors []error) (res inspectResult) {
	first := true
	for i, err := range errors {
		if err == nil {
			continue
		}

		res.Count++
		if first {
			first = false
			res.FirstErrorIdx = i
		}

		if merr, ok := err.(*multiError); ok { //nolint:errorlint
			res.Capacity += len(merr.errors)
			res.ContainsMultiError = true
		} else {
			res.Capacity++
		}
	}
	return
}

// fromSlice converts the given list of errors into a single error.
func fromSlice(errors []error) error {
	res := inspect(errors)
	switch res.Count {
	case 0:
		return nil
	// case 1:
	// 	// only one non-nil entry
	// 	return errors[res.FirstErrorIdx]
	case len(errors):
		if !res.ContainsMultiError {
			// already flat
			return &multiError{errors: errors}
		}
	}

	nonNilErrs := make([]error, 0, res.Capacity)
	// for _, err := range errors[res.FirstErrorIdx:] {
	for _, err := range errors {
		if err == nil {
			continue
		}

		if nested, ok := err.(*multiError); ok { //nolint:errorlint
			nonNilErrs = append(nonNilErrs, nested.errors...)
		} else {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	return &multiError{errors: nonNilErrs}
}

// Append создаст цепочку ошибок из ошибок ...errors. Допускается использование `nil` в аргументах.
func Append(errors ...error) error {
	return fromSlice(errors)
}

// Wrap обернет ошибку `left` ошибкой `right`, получив цепочку. Допускается использование `nil` в обоих аргументах.
func Wrap(left error, right error) error {
	return fromSlice([]error{left, right})
}
