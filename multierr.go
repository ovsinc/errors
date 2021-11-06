// Используется оригинальный код проекта "go.uber.org/multierr" с частичным заимствованием.
// Код проекта "go.uber.org/multierr" распространяется под лицензией MIT (https://github.com/uber-go/multierr/blob/master/LICENSE.txt).

package errors

import (
	"fmt"
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
	if eg, ok := err.(Multierror); ok { //nolint:errorlint
		return eg.Errors()
	}

	return appendError(make([]*Error, 0, 1), err)
}

var _ Multierror = (*multiError)(nil)

type Multierror interface {
	Errors() []*Error
	Error() string
	Format(f fmt.State, c rune)
	Marshal(fn ...Marshaller) ([]byte, error)
	Len() int
}

type multiError struct {
	errors []*Error
	len    int
	last   int
}

// Errors returns the list of underlying errors.
//
// This slice MUST NOT be modified.
func (merr *multiError) Errors() []*Error {
	if merr == nil || merr.errors == nil {
		return []*Error{}
	}
	return merr.errors
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

	data, _ := marshal.Marshal(merr)
	_, _ = f.Write(data)
}

func (merr *multiError) Len() int {
	return merr.len
}

// Unwrap вернет самую новую ошибку в стеке
func (merr *multiError) Unwrap() error {
	es := merr.Errors()
	if len(es) == 0 {
		return nil
	}
	return es[merr.last]
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
