package errors

import (
	"fmt"
	"strconv"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/valyala/bytebufferpool"
)

// DefaultMultierrFormatFunc функция форматирования для multierr ошибок.
var DefaultMultierrFormatFunc = StringMultierrFormatFunc

// JSONMultierrFuncFormat функция форматирования вывода сообщения для multierr в виде JSON.
func JSONMultierrFuncFormat(es []error) string {
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

// StringMultierrFormatFunc функция форматирования вывода сообщения для multierr в виде строки.
// Используется по-умолчанию.
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

// хелперы

// Append создаст или дополнит цепочку ошибок err с помощью errs.
// err и errs[N] могут быть nil.
func Append(err error, errs ...error) *multierror.Error {
	me := multierror.Append(err, errs...)
	me.ErrorFormat = DefaultMultierrFormatFunc
	return me
}

// Wrap обернет ошибку olderr в err и вернет цепочку.
// err и olderr могут быть nil.
func Wrap(olderr error, err error) *multierror.Error {
	me := Append(olderr, err)
	me.ErrorFormat = DefaultMultierrFormatFunc
	return me
}
