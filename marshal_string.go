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
)

var _ Marshaller = (*MarshalString)(nil)

type MarshalString struct{}

func (MarshalString) MarshalTo(i interface{}, dst io.Writer) error {
	return stringMarshalTo(i, dst)
}

func stringMarshalTo(i interface{}, dst io.Writer) error {
	if i == nil {
		return nil
	}

	switch t := i.(type) { //nolint:errorlint
	case Multierror: // multiError
		stringMultierrFormat(dst, t.Errors())

	case *Error:
		stringFormat(dst, t)
	}

	return nil
}

func (MarshalString) Marshal(i interface{}) ([]byte, error) {
	if i == nil {
		return []byte{}, nil
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	_ = stringMarshalTo(i, buf)

	return buf.Bytes(), nil
}

func writeobject(buf io.Writer, o Objecter) {
	if n, _ := o.Write(buf); n > 0 {
		_, _ = buf.Write(_separator)
	}
}

func writecontext(buf io.Writer, ctxs CtxMap) {
	ctxslen := len(ctxs)
	if ctxslen == 0 {
		return
	}

	_, _ = buf.Write(_ctxDelimiterLeft)
	// отсортируем по ключам запишем в буфер в формате
	// {<key>:<value>,<key>:<value>}
	ctxskeys := make([]string, 0, ctxslen)
	for i := range ctxs {
		ctxskeys = append(ctxskeys, i)
	}
	sort.Strings(ctxskeys)
	_, _ = fmt.Fprintf(buf, "%s:%v", ctxskeys[0], ctxs[ctxskeys[0]])
	for _, i := range ctxskeys[1:] {
		_, _ = buf.Write(_listSeparator)
		_, _ = fmt.Fprintf(buf, "%s:%v", i, ctxs[i])
	}
	_, _ = buf.Write(_ctxDelimiterRight)
	_, _ = buf.Write(_separator)
}

func stringFormat(buf io.Writer, e *Error) {
	if e == nil {
		return
	}

	writeobject(buf, e.FileLine())  //nolint:ineffassign,staticcheck
	writeobject(buf, e.ErrorType()) //nolint:ineffassign,staticcheck
	writeobject(buf, e.Operation()) //nolint:ineffassign,staticcheck
	writecontext(buf, e.ContextInfo())

	_, _ = e.WriteTranslateMsg(buf)
}

func stringMultierrFormat(w io.Writer, es []*Error) {
	if len(es) == 0 {
		return
	}

	_, _ = w.Write(_multilinePrefix)
	_, _ = w.Write(_multilineSeparator)

	for i, err := range es {
		if err == nil {
			continue
		}
		_, _ = w.Write(_multilineIndent)
		_, _ = w.Write([]byte(strconv.Itoa(i + 1)))
		_, _ = w.Write([]byte(" "))

		_, _ = io.WriteString(w, err.Error())

		_, _ = w.Write(_multilineSeparator)
	}
}
