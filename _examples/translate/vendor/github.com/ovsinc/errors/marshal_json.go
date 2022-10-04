package errors

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

var (
	ErrUnknownType = New("unknown type")

	_ Marshaller = (*MarshalJSON)(nil)
)

type MarshalJSON struct{}

func (MarshalJSON) MarshalTo(i interface{}, dst io.Writer) error {
	switch t := i.(type) { //nolint:errorlint
	case nil:
		return nil
	case Multierror:
		jsonMultierrFormat(dst, t.Errors())
	case *Error:
		jsonFormat(dst, t)
	}
	return nil
}

func (m *MarshalJSON) Marshal(i interface{}) ([]byte, error) {
	if i == nil {
		return []byte{}, nil
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	_ = m.MarshalTo(i, buf)

	return buf.Bytes(), nil
}

func jsonFormat(buf io.Writer, e *Error) {
	if e == nil {
		_, _ = io.WriteString(buf, "null")
		return
	}

	_, _ = io.WriteString(buf, "{")

	// ID
	_, _ = io.WriteString(buf, "\"id\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.ID())
	_, _ = io.WriteString(buf, "\",")

	// Operation
	_, _ = io.WriteString(buf, "\"operation\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.Operation())
	_, _ = io.WriteString(buf, "\",")

	// ContextInfo
	_, _ = io.WriteString(buf, "\"context\":")
	if cxtInfo := e.ContextInfo(); len(cxtInfo) > 0 {
		_, _ = io.WriteString(buf, "{")
		ctxskeys := make([]string, 0, len(cxtInfo))
		for i := range cxtInfo {
			ctxskeys = append(ctxskeys, i)
		}
		sort.Strings(ctxskeys)
		_, _ = fmt.Fprintf(buf, "\"%s\":\"%v\"", ctxskeys[0], cxtInfo[ctxskeys[0]])
		for _, i := range ctxskeys[1:] {
			_, _ = buf.Write(_listSeparator)
			_, _ = fmt.Fprintf(buf, "\"%s\":\"%v\"", i, cxtInfo[i])
		}
		_, _ = io.WriteString(buf, "}")
	} else {
		_, _ = io.WriteString(buf, "null")
	}
	_, _ = io.WriteString(buf, ",")

	// Msg
	_, _ = io.WriteString(buf, "\"msg\":")
	_, _ = io.WriteString(buf, "\"")
	_, _ = buf.Write(e.Msg())
	_, _ = io.WriteString(buf, "\"")

	_, _ = io.WriteString(buf, "}")
}

// JSONMultierrFuncFormat функция форматирования вывода сообщения для multierr в виде JSON.
func jsonMultierrFormat(w io.Writer, es []*Error) {
	l := len(es)
	if l == 0 {
		_, _ = io.WriteString(w, "null")
		return
	}

	_, _ = io.WriteString(w, "{")

	_, _ = io.WriteString(w, "\"count\":")
	_, _ = io.WriteString(w, strconv.Itoa(l))
	_, _ = io.WriteString(w, ",")

	_, _ = io.WriteString(w, "\"messages\":")
	_, _ = io.WriteString(w, "[")
	switch l {
	case 1:
		jsonFormat(w, es[0])
	default:
		jsonFormat(w, es[0])
		for _, e := range es[1:] {
			_, _ = io.WriteString(w, ",")
			jsonFormat(w, e)
		}
	}
	_, _ = io.WriteString(w, "]")

	_, _ = io.WriteString(w, "}")
}
