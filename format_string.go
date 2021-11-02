package errors

import (
	"fmt"
	"io"
	"sort"
	"strconv"
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
)

// StringFormat функция форматирования вывода сообщения *Error в виде строки.
// Используется по-умолчанию.
func StringFormat(buf io.Writer, e *Error) { //nolint:cyclop
	if e == nil {
		return
	}

	writeDelim := false

	if op := e.Operation().Bytes(); len(op) > 0 {
		_, _ = buf.Write(op)
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
			buf.Write(_listSeparator)
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

//

// multierr

// StringMultierrFormatFunc функция форматирования вывода сообщения для multierr в виде строки.
// Используется по-умолчанию.
func StringMultierrFormatFunc(w io.Writer, es []*Error) {
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
