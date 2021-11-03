package errors

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

var (
	_multilinePrefix    = []byte("the following errors occurred:") //nolint:gochecknoglobals
	_multilineSeparator = []byte("\n")                             //nolint:gochecknoglobals
	_multilineIndent    = []byte("\t#")                            //nolint:gochecknoglobals
	_msgSeparator       = []byte(" -- ")                           //nolint:gochecknoglobals
	_listSeparator      = []byte(",")                              //nolint:gochecknoglobals
	_opDelimiter        = []byte(": ")                             //nolint:gochecknoglobals

	_ctxDelimiterLeft  = []byte("{") //nolint:gochecknoglobals
	_ctxDelimiterRight = []byte("}") //nolint:gochecknoglobals

	_ Marshaller = (*MarshalString)(nil)
)

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

func stringFormat(buf io.Writer, e *Error) { //nolint:cyclop
	if e == nil {
		return
	}

	writeDelim := false

	n, _ := e.Operation().Write(buf)
	if n > 0 {
		_, _ = buf.Write(_opDelimiter)
		writeDelim = true
	}

	if ctxs := e.ContextInfo(); len(ctxs) > 0 {
		_, _ = buf.Write(_ctxDelimiterLeft)
		// отсортируем по ключам запишем в буфер в формате
		// {<key>:<value>,<key>:<value>}
		ctxskeys := make([]string, 0, len(ctxs))
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
		writeDelim = true
	}

	if writeDelim && len(e.Msg().Bytes()) > 0 {
		_, _ = buf.Write(_msgSeparator)
	}

	_, _ = e.WriteTranslateMsg(buf)
}

func stringMultierrFormat(w io.Writer, es []*Error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "")
		return
	}

	_, _ = w.Write(_multilinePrefix)
	_, _ = w.Write(_multilineSeparator)

	for i, err := range es {
		if err == nil {
			continue
		}
		_, _ = w.Write(_multilineIndent)
		_, _ = io.WriteString(w, strconv.Itoa(i+1))
		_, _ = io.WriteString(w, " ")
		_, _ = io.WriteString(w, err.Error())
		_, _ = w.Write(_multilineSeparator)
	}
}
