package errors

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	json "github.com/goccy/go-json"
)

var (
	// DefaultFormatFn функция форматирования, используемая по-умолчанию
	DefaultFormatFn FormatFn //nolint:gochecknoglobals

	// DefaultMultierrFormatFunc функция форматирования для multierr ошибок.
	DefaultMultierrFormatFunc MultierrFormatFn //nolint:gochecknoglobals

	_multilinePrefix    = []byte("the following errors occurred:") //nolint:gochecknoglobals
	_multilineSeparator = []byte("\n")                             //nolint:gochecknoglobals
	_multilineIndent    = []byte("* ")                             //nolint:gochecknoglobals
	_msgSeparator       = []byte(" -- ")                           //nolint:gochecknoglobals
)

type (
	// FormatFn тип функции форматирования.
	FormatFn func(w io.Writer, e *Error)

	// MultierrFormatFn типу функции морматирования для multierr.
	MultierrFormatFn func(w io.Writer, es []error)
)

// JSONFormat функция форматирования вывода сообщения *Error в JSON.
func JSONFormat(buf io.Writer, e *Error) {
	if e == nil {
		_, _ = io.WriteString(buf, "null")
		return
	}

	_, _ = io.WriteString(buf, "{")

	// ID
	_, _ = io.WriteString(buf, "\"id\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = io.WriteString(buf, e.ID())
	_, _ = io.WriteString(buf, "\",")

	// ErrorType
	_, _ = io.WriteString(buf, "\"error_type\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = io.WriteString(buf, e.ErrorType())
	_, _ = io.WriteString(buf, "\",")

	// Severity
	_, _ = io.WriteString(buf, "\"severity\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = io.WriteString(buf, e.Severity().String())
	_, _ = io.WriteString(buf, "\",")

	// Operations
	_, _ = io.WriteString(buf, "\"operations\":[")
	ops := e.Operations()
	if len(ops) > 0 {
		op0 := ops[0]
		_, _ = io.WriteString(buf, "\"")
		_, _ = io.WriteString(buf, op0)
		_, _ = io.WriteString(buf, "\"")
		for _, opN := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = io.WriteString(buf, "\"")
			_, _ = io.WriteString(buf, opN)
			_, _ = io.WriteString(buf, "\"")
		}
	}
	_, _ = io.WriteString(buf, "],")

	// ContextInfo
	_, _ = io.WriteString(buf, "\"context\":")
	cxtInfo := e.ContextInfo()
	if len(cxtInfo) > 0 {
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "")
		_ = enc.Encode(e.ContextInfo())
	} else {
		_, _ = io.WriteString(buf, "null")
	}
	_, _ = io.WriteString(buf, ",")

	// Msg
	_, _ = io.WriteString(buf, "\"msg\":")
	_, _ = io.WriteString(buf, "\"")
	if len(e.Msg()) == 0 {
		_, _ = io.WriteString(buf, "")
	} else {
		_, _ = e.WriteTranslateMsg(buf)
	}
	_, _ = io.WriteString(buf, "\"")

	_, _ = io.WriteString(buf, "}")
}

// StringFormat функция форматирования вывода сообщения *Error в виде строки.
// Используется по-умолчанию.
func StringFormat(buf io.Writer, e *Error) { //nolint:cyclop
	if e == nil {
		return
	}

	writeDelim := false

	if et := e.ErrorType(); len(et) > 0 {
		_, _ = io.WriteString(buf, "(")
		_, _ = io.WriteString(buf, et)
		_, _ = io.WriteString(buf, ")")
		writeDelim = true
	}

	if ops := e.Operations(); len(ops) > 0 {
		_, _ = io.WriteString(buf, "[")
		op0 := ops[0]
		_, _ = io.WriteString(buf, op0)
		for _, opN := range ops[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = io.WriteString(buf, opN)
		}
		_, _ = io.WriteString(buf, "]")
		writeDelim = true
	}

	if ctxs := e.ContextInfo(); len(ctxs) > 0 {
		_, _ = io.WriteString(buf, "<")
		ctxskeys := make([]string, 0, len(ctxs))
		for i := range ctxs {
			ctxskeys = append(ctxskeys, i)
		}
		sort.Strings(ctxskeys)
		_, _ = fmt.Fprintf(buf, "%s:%v", ctxskeys[0], ctxs[ctxskeys[0]])
		for _, i := range ctxskeys[1:] {
			_, _ = io.WriteString(buf, ",")
			_, _ = fmt.Fprintf(buf, "%s:%v", i, ctxs[i])
		}
		_, _ = io.WriteString(buf, ">")
		writeDelim = true
	}

	if msg := e.Msg(); writeDelim && len(msg) > 0 {
		_, _ = buf.Write(_msgSeparator)
	}

	_, _ = e.WriteTranslateMsg(buf)
}

//

// multierr

// StringMultierrFormatFunc функция форматирования вывода сообщения для multierr в виде строки.
// Используется по-умолчанию.
func StringMultierrFormatFunc(w io.Writer, es []error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "")
		return
	}

	_, _ = w.Write(_multilinePrefix)
	_, _ = w.Write(_multilineSeparator)

	for _, err := range es {
		if err == nil {
			continue
		}
		_, _ = w.Write(_multilineIndent)
		_, _ = io.WriteString(w, err.Error())
		_, _ = w.Write(_multilineSeparator)
	}
}

// JSONMultierrFuncFormat функция форматирования вывода сообщения для multierr в виде JSON.
func JSONMultierrFuncFormat(w io.Writer, es []error) {
	if len(es) == 0 {
		_, _ = io.WriteString(w, "null")
	}

	_, _ = io.WriteString(w, "{")

	_, _ = io.WriteString(w, "\"count\":")
	_, _ = io.WriteString(w, strconv.Itoa(len(es)))
	_, _ = io.WriteString(w, ",")

	_, _ = io.WriteString(w, "\"messages\":")
	_, _ = io.WriteString(w, "[")
	writeErrFn := func(e error) {
		if e == nil {
			return
		}
		if myerr, ok := simpleCast(e); ok {
			JSONFormat(w, myerr)
			return
		}
		_, _ = fmt.Fprintf(w, "\"%v\"", e)
	}
	switch len(es) {
	case 1:
		writeErrFn(es[0])
	default:
		writeErrFn(es[0])
		for _, e := range es[1:] {
			_, _ = io.WriteString(w, ",")
			writeErrFn(e)
		}
	}
	_, _ = io.WriteString(w, "]")

	_, _ = io.WriteString(w, "}")
}
