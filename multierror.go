package errors

import (
	"bytes"
	"fmt"

	"go.uber.org/atomic"
)

// Combine создаст цепочку ошибок из ошибок ...errors.
// Допускается использование `nil` в аргументах.
func Combine(errors ...error) error {
	return fromSlice(errors)
}

// CombineWithLog как и Combine создаст или дополнит цепочку ошибок err с помощью errs,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func CombineWithLog(errs ...error) error {
	e := Combine(errs...)
	Log(e)
	return e
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

	if _, ok := right.(*multiError); !ok { //nolint:errorlint
		if l, ok := left.(*multiError); ok && !l.copyNeeded.Swap(true) { //nolint:errorlint
			// Common case where the error on the left is constantly being
			// appended to.
			return &multiError{
				errors: append(l.errors, right),
				cur:    atomic.NewInt32(l.cur.Load()),
			}
		} else if !ok {
			// Both errors are single errors.
			return &multiError{
				errors: []error{left, right},
				cur:    atomic.NewInt32(int32(0)),
			}
		}
	}

	// Either right or both, left and right, are multiErrors. Rely on usual
	// expensive logic.
	errors := [2]error{left, right}
	return fromSlice(errors[0:])
}

// WrapWithLog обернет ошибку olderr в err и вернет цепочку,
// но при этом будет осуществлено логгирование с помощь логгера по-умолчанию.
func WrapWithLog(olderr error, err error) error {
	e := Wrap(olderr, err)
	Log(e)
	return e
}

var _ Multierror = (*multiError)(nil)

type Multierror interface {
	Errors() []error
	Error() string
	Format(f fmt.State, c rune)
	Marshal(fn ...Marshaller) ([]byte, error)
	Len() int
	Log(l ...Logger)
	Unwrap() error
	Last() error
	FindByID(id []byte) (error, bool)
}

type multiError struct {
	errors     []error
	cur        *atomic.Int32
	copyNeeded *atomic.Bool
}

// Errors returns the copy list of underlying errors.
func (merr *multiError) Errors() []error {
	if merr == nil || merr.errors == nil {
		return nil
	}
	return append(([]error)(nil), merr.errors...)
}

func (merr *multiError) Error() string {
	marshal := &MarshalString{}
	data, _ := marshal.Marshal(merr)
	return string(data)
}

func (merr *multiError) Marshal(fn ...Marshaller) ([]byte, error) {
	var marshal Marshaller
	switch {
	case len(fn) > 0:
		marshal = fn[0]
	case DefaultMarshaller != nil:
		marshal = DefaultMarshaller
	}
	return marshal.Marshal(merr)
}

func (merr *multiError) Format(f fmt.State, c rune) {
	var marshal Marshaller
	switch c {
	case 'w', 'v', 's':
		marshal = &MarshalString{}
	case 'j':
		marshal = &MarshalJSON{}
	}
	_ = marshal.MarshalTo(merr, f)
}

func (merr *multiError) Len() int {
	return len(merr.errors)
}

// Unwrap вернет самую новую ошибку в стеке
func (merr *multiError) Unwrap() error {
	return merr.Last()
}

// Last вернет самую новую (*Error) ошибку в стеке
func (merr *multiError) Last() error {
	if merr.errors == nil {
		return nil
	}
	return merr.errors[0]
}

func (merr *multiError) As(target interface{}) bool {
	if x, ok := target.(*Multierror); ok { //nolint:errorlint
		*x = merr
		return true
	}
	return false
}

func (merr *multiError) Is(target error) bool {
	if x, ok := target.(Multierror); ok { //nolint:errorlint
		return x == merr
	}
	return false
}

// Log выполнит логгирование ошибки с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func (merr *multiError) Log(l ...Logger) {
	Log(merr, l...)
}

func (merr *multiError) FindByID(id []byte) (error, bool) {
	if merr.errors == nil {
		return nil, false
	}

	if e, ok := merr.errors[merr.cur.Load()].(*Error); ok { //nolint:errorlint
		if bytes.Equal(e.ID(), id) {
			return e, true
		}
	}

	for i, e := range merr.errors {
		if ee, ok := e.(*Error); ok { //nolint:errorlint
			if bytes.Equal(ee.ID(), id) {
				merr.cur.Store(int32(i))
				return e, true
			}
		}
	}

	return nil, false
}

type inspectResult struct {
	// Number of top-level non-nil errors
	count int

	// Total number of errors including multiErrors
	capacity int

	// Index of the first non-nil error in the list. Value is meaningless if
	// Count is zero.
	firstErrorIdx int

	// Whether the list contains at least one multiError
	containsMultiError bool
}

// Inspects the given slice of errors so that we can efficiently allocate
// space for it.
func inspect(errors []error) (res inspectResult) {
	first := true
	for i, err := range errors {
		if err == nil {
			continue
		}

		res.count++
		if first {
			first = false
			res.firstErrorIdx = i
		}

		if merr, ok := err.(*multiError); ok { //nolint:errorlint
			res.capacity += len(merr.errors)
			res.containsMultiError = true
		} else {
			res.capacity++
		}
	}
	return
}

// fromSlice converts the given list of errors into a single error.
func fromSlice(errors []error) error {
	// Don't pay to inspect small slices.
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	}

	res := inspect(errors)
	switch res.count {
	case 0:
		return nil
	case 1:
		// only one non-nil entry
		return errors[res.firstErrorIdx]
	case len(errors):
		if !res.containsMultiError {
			// Error list is flat. Make a copy of it
			// Otherwise "errors" escapes to the heap
			// unconditionally for all other cases.
			// This lets us optimize for the "no errors" case.
			out := append(([]error)(nil), errors...)
			return &multiError{errors: out, cur: atomic.NewInt32(int32(res.firstErrorIdx))}
		}
	}

	nonNilErrs := make([]error, 0, res.capacity)
	for _, err := range errors[res.firstErrorIdx:] {
		if err == nil {
			continue
		}

		if nested, ok := err.(*multiError); ok { //nolint:errorlint
			nonNilErrs = append(nonNilErrs, nested.errors...)
		} else {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	return &multiError{errors: nonNilErrs, cur: atomic.NewInt32(int32(0))}
}
