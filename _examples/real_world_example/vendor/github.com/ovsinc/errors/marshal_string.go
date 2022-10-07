package errors

import (
	"io"
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
	case error: // one
		stringFormat(dst, t)
	}

	return nil
}

func (m *MarshalString) Marshal(i interface{}) ([]byte, error) {
	if i == nil {
		return nil, nil
	}

	buf := bytebufferpool.Get()
	_ = m.MarshalTo(i, buf)
	data := buf.Bytes()
	bytebufferpool.Put(buf)

	return data, nil
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

func contextInfoFormat(w io.Writer, ctxi CtxKV, useDelimiter bool) {
	if len(ctxi) < 1 {
		return
	}

	if useDelimiter {
		_, _ = w.Write(_ctxDelimiterLeft)
	}

	// 0
	_, _ = w.Write(ctxi[0].Key)
	_, _ = io.WriteString(w, ":")
	_, _ = w.Write(ctxi[0].Value)
	// other
	for _, i := range ctxi[1:] {
		_, _ = w.Write(_listSeparator)
		_, _ = w.Write(i.Key)
		_, _ = io.WriteString(w, ":")
		_, _ = w.Write(i.Value)
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
		if et := t.ErrorType(); et != 0 {
			_, _ = w.Write(_errTypeDelimerLeft)
			_, _ = io.WriteString(w, t.ErrorType().String())

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
		contextInfoFormat(w, t.ContextInfo(), true)

		// msg
		_, _ = w.Write(t.Msg())

	default:
		_, _ = io.WriteString(w, t.Error())
	}
}
