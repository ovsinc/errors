package errors

import (
	"fmt"
	"io"
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
		_, _ = io.WriteString(dst, "null")
		return nil
	case interface{ Errors() []error }:
		jsonMultierrFormat(dst, t.Errors())
	case error:
		jsonFormat(dst, t)
	}
	return nil
}

func (m *MarshalJSON) Marshal(i interface{}) ([]byte, error) {
	if i == nil {
		return []byte{}, nil
	}

	buf := bytebufferpool.Get()
	_ = m.MarshalTo(i, buf)
	data := buf.Bytes()
	bytebufferpool.Put(buf)

	return data, nil
}

func jsonFormat(buf io.Writer, e error) { //nolint:funlen
	switch t := e.(type) { //nolint:errorlint
	case *Error:
		_, _ = io.WriteString(buf, "{")

		// ID
		_, _ = io.WriteString(buf, "\"id\":")
		_, _ = io.WriteString(buf, "\"")
		_, _ = buf.Write(s2b(t.ID()))
		_, _ = io.WriteString(buf, "\",")

		// Operation
		_, _ = io.WriteString(buf, "\"operation\":")
		_, _ = io.WriteString(buf, "\"")
		_, _ = buf.Write(s2b(t.Operation()))
		_, _ = io.WriteString(buf, "\",")

		// ErrorType
		_, _ = io.WriteString(buf, "\"error_type\":")
		_, _ = io.WriteString(buf, "\"")
		_, _ = io.WriteString(buf, t.ErrorType().String())
		_, _ = io.WriteString(buf, "\",")

		// ContextInfo
		_, _ = io.WriteString(buf, "\"context\":")
		if cxtInfo := t.ContextInfo(); len(cxtInfo) > 0 {
			_, _ = io.WriteString(buf, "{")
			// 0
			_, _ = io.WriteString(buf, "\"")
			_, _ = buf.Write(s2b(cxtInfo[0].Key))
			_, _ = io.WriteString(buf, "\":\"")
			_, _ = fmt.Fprint(buf, cxtInfo[0].Value)
			_, _ = io.WriteString(buf, "\"")
			// other
			for _, i := range cxtInfo[1:] {
				_, _ = buf.Write(_listSeparator)
				_, _ = io.WriteString(buf, "\"")
				_, _ = buf.Write(s2b(i.Key))
				_, _ = io.WriteString(buf, "\":\"")
				_, _ = fmt.Fprint(buf, i.Value)
				_, _ = io.WriteString(buf, "\"")
			}
			_, _ = io.WriteString(buf, "}")
		} else {
			_, _ = io.WriteString(buf, "null")
		}
		_, _ = io.WriteString(buf, ",")

		// Msg
		_, _ = io.WriteString(buf, "\"msg\":")
		_, _ = io.WriteString(buf, "\"")
		_, _ = buf.Write(s2b(t.Msg()))
		_, _ = io.WriteString(buf, "\"")

		_, _ = io.WriteString(buf, "}")

	default:
		_, _ = io.WriteString(buf, "{")
		_, _ = io.WriteString(buf, "\"msg\":")
		_, _ = io.WriteString(buf, "\"")
		_, _ = io.WriteString(buf, t.Error())
		_, _ = io.WriteString(buf, "\"")
		_, _ = io.WriteString(buf, "}")
	}
}

// JSONMultierrFuncFormat функция форматирования вывода сообщения для multierr в виде JSON.
func jsonMultierrFormat(w io.Writer, es []error) {
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
