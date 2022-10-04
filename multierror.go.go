package errors

import "fmt"

// Combine создаст цепочку ошибок из ошибок ...errors.
// Допускается использование `nil` в аргументах.
func Combine(errors ...error) Multierror {
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
func Wrap(left error, right error) Multierror {
	return fromSlice([]error{left, right})
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
	Errors() []*Error
	Error() string
	Format(f fmt.State, c rune)
	Marshal(fn ...Marshaller) ([]byte, error)
	Len() int
	Log(l ...Logger)
	Unwrap() error
	Last() *Error
}

type multiError struct {
	errors []*Error
	len    int
	last   int
}

// Errors returns the copy list of underlying errors.
func (merr *multiError) Errors() []*Error {
	if merr == nil || merr.errors == nil {
		return []*Error{}
	}
	result := make([]*Error, len(merr.errors))
	copy(result, merr.errors)
	return result
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
	return merr.len
}

// Unwrap вернет самую новую ошибку в стеке
func (merr *multiError) Unwrap() error {
	return merr.Last()
}

// Last вернет самую новую (*Error) ошибку в стеке
func (merr *multiError) Last() *Error {
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

// Log выполнит логгирование ошибки с ипользованием логгера l[0].
// Если l не указан, то в качестве логгера будет использоваться логгер по-умолчанию.
func (merr *multiError) Log(l ...Logger) {
	Log(merr, l...)
}

// Нужна оптимизация по примеру uber https://github.com/uber-go/multierr/blob/master/error.go#L333

// append err to []*Error
// errors must not be nil
func appendError(errors []*Error, err interface{}) []*Error {
	switch t := err.(type) {
	case nil:
		return nil

	case *Error:
		return append(errors, t)

	case Multierror:
		return append(errors, t.Errors()...)

	case error:
		return append(errors, New(t))
	}

	return errors
}

// fromSlice converts the given list of errors into a single error.
func fromSlice(errors []error) Multierror {
	nonNilErrs := make([]*Error, 0, len(errors)+1)
	for _, err := range errors {
		if err == nil {
			continue
		}
		nonNilErrs = appendError(nonNilErrs, err)
	}

	last := 0
	len := len(nonNilErrs)
	if len > 0 {
		last = len - 1
	}

	return &multiError{
		errors: nonNilErrs,
		len:    len,
		last:   last,
	}
}
