package errors

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

var (
	_ctxDelimiterLeft  = []byte{'{'}
	_ctxDelimiterRight = []byte{'}'}
	_listSeparator     = []byte{','}

	_multilineIndent    = []byte("\t#")
	_multilineSeparator = []byte{'\n'}
	_multilinePrefix    = []byte("the following errors occurred:")

	_separator = []byte{' '}

	_opDelimiterLeft  = []byte{'['}
	_opDelimiterRight = []byte{']'}

	_errTypeDelimerLeft  = []byte{'('}
	_errTypeDelimerRight = []byte{')'}
)

var _ Marshaller = (*MarshalString)(nil)

type MarshalString struct{}

func (m *MarshalString) MarshalTo(i interface{}, dst io.Writer) error {
	switch t := i.(type) { //nolint:errorlint
	case nil:
		return nil
	case interface{ Errors() []error }: // multiError
		stringMultierrFormat(dst, t.Errors())
	case error: //one
		stringFormat(dst, t)
	}

	return nil
}

func (m *MarshalString) Marshal(i interface{}) ([]byte, error) {
	if i == nil {
		return nil, nil
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	_ = m.MarshalTo(i, buf)

	return buf.Bytes(), nil
}

//

func stringMultierrFormat(w io.Writer, es []error) {
	_, _ = w.Write(_multilinePrefix)
	_, _ = w.Write(_multilineSeparator)
	for i, err := range es {
		if err == nil {
			continue
		}
		_, _ = w.Write(_multilineIndent)
		_, _ = w.Write([]byte(strconv.Itoa(i + 1)))
		_, _ = w.Write([]byte(" "))
		stringFormat(w, err)
		_, _ = w.Write(_multilineSeparator)
	}
}

func contextInfoFormat(w io.Writer, ctxiptr *CtxMap, useDelimiter bool) {
	if ctxiptr == nil {
		return
	}

	ctxi := *ctxiptr
	if len(ctxi) < 1 {
		return
	}

	if useDelimiter {
		_, _ = w.Write(_ctxDelimiterLeft)
	}

	// отсортируем по ключам запишем в буфер в формате
	// {<key>:<value>,<key>:<value>}
	ctxskeys := make([]string, 0, len(ctxi))
	for i := range ctxi {
		ctxskeys = append(ctxskeys, i)
	}
	sort.Strings(ctxskeys)

	_, _ = fmt.Fprintf(w, "%s:%v", ctxskeys[0], ctxi[ctxskeys[0]])
	for _, i := range ctxskeys[1:] {
		_, _ = w.Write(_listSeparator)
		_, _ = fmt.Fprintf(w, "%s:%v", i, ctxi[i])
	}

	if useDelimiter {
		_, _ = w.Write(_ctxDelimiterRight)
		_, _ = w.Write(_separator)
	}
}

func stringFormat(w io.Writer, e error) {
	switch t := e.(type) { //nolint:errorlint
	case *Error:
		// id do not write

		// err type
		if t := t.ErrorType(); t != nil {
			_, _ = w.Write(_errTypeDelimerLeft)
			_, _ = w.Write(t)
			_, _ = w.Write(_errTypeDelimerRight)
			_, _ = w.Write(_separator)
		}

		// operation
		if op := t.Operation(); op != nil {
			_, _ = w.Write(_opDelimiterLeft)
			_, _ = w.Write(op)
			_, _ = w.Write(_opDelimiterRight)
			_, _ = w.Write(_separator)
		}

		// ctx
		ctxi := t.ContextInfo()
		contextInfoFormat(w, &ctxi, true)

		// msg
		_, _ = w.Write(t.Msg())

	default:
		_, _ = io.WriteString(w, t.Error())
	}
}
