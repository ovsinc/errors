package errors

import (
	"fmt"
	"strconv"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/valyala/bytebufferpool"
)

var DefaultMultierrFormatFunc = StringMultierrFormatFunc

func JsonMultierrFuncFormat(es []error) string {
	if len(es) == 0 {
		return "null"
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	_, _ = buf.WriteString("{")

	_, _ = buf.WriteString("\"count\":" + strconv.Itoa(len(es)) + ",")

	_, _ = buf.WriteString("\"messages\":")
	_, _ = buf.WriteString("[")
	writeErrFn := func(e error) {
		switch t := e.(type) { // nolint:errorlint
		case *Error:
			_, _ = buf.WriteString(JSONFormat(t))
		default:
			_, _ = buf.WriteString("\"" + fmt.Sprintf("%v", t) + "\"")
		}
	}
	switch len(es) {
	case 0:
	case 1:
		writeErrFn(es[0])
	default:
		writeErrFn(es[0])
		for _, e := range es[1:] {
			_, _ = buf.WriteString(",")
			writeErrFn(e)
		}
	}
	_, _ = buf.WriteString("]")

	_, _ = buf.WriteString("}")

	return buf.String()
}

func StringMultierrFormatFunc(es []error) string {
	if len(es) == 0 {
		return ""
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	writeErrFn := func(err error) {
		switch t := err.(type) { // nolint:errorlint
		case *Error:
			_, _ = buf.WriteString("* " + StringFormat(t))
		default:
			_, _ = buf.WriteString(fmt.Sprintf("* %v", t))
		}
	}

	for _, err := range es {
		writeErrFn(err)
		_, _ = buf.WriteString("\n")
	}

	return buf.String()
}

//

// Append proxy func for multierror.Append with Custom format
func Append(err error, errs ...error) *multierror.Error {
	me := multierror.Append(err, errs...)
	me.ErrorFormat = DefaultMultierrFormatFunc
	return me
}

// Wrap proxy func for multierror.Append with Custom format
func Wrap(olderr error, err error) *multierror.Error {
	me := Append(olderr, err)
	me.ErrorFormat = DefaultMultierrFormatFunc
	return me
}

// Unwrap позволяет получить оригинальную ошибку
// для multierror будет разернута цепочка ошибок
func Unwrap(err error) error {
	type unwraper interface {
		Unwrap() error
	}
	for err != nil {
		unwrap, ok := err.(unwraper)
		if !ok {
			break
		}
		err = unwrap.Unwrap()
	}
	return err
}
