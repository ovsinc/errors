package errors

import (
	"io"
)

var (
	_opDelimiterLeft     = []byte{'['}
	_opDelimiterRight    = []byte{']'}
	_errTypeDelimerLeft  = []byte{'('}
	_errTypeDelimerRight = []byte{')'}
	_separator           = []byte{' '}
)

// One time object
type Objecter interface {
	String() string
	Bytes() []byte
	Write(w io.Writer) (int, error)
	Len() int
}

type object struct {
	data []byte
}

func NewObjectEmpty() Objecter {
	return &object{}
}

func NewIDFromString(s string) Objecter {
	return NewObjectFromString(s, nil, nil)
}

func NewIDFromBytes(s []byte) Objecter {
	return NewObjectFromBytes(s, nil, nil)
}

func NewOperationFromString(s string) Objecter {
	return NewObjectFromString(s, _opDelimiterLeft, _opDelimiterRight)
}

func NewOperationFromBytes(s []byte) Objecter {
	return NewObjectFromBytes(s, _opDelimiterLeft, _opDelimiterRight)
}

func NewMsgFromString(s string) Objecter {
	return NewObjectFromString(s, nil, nil)
}

func NewMsgFromBytes(s []byte) Objecter {
	return NewObjectFromBytes(s, nil, nil)
}

func NewErrorTypeFromString(s string) Objecter {
	return NewObjectFromString(s, _errTypeDelimerLeft, _errTypeDelimerRight)
}

func NewErrorTypeFromBytes(s []byte) Objecter {
	return NewObjectFromBytes(s, _errTypeDelimerLeft, _errTypeDelimerRight)
}

func NewFileLineFromBytes(s []byte) Objecter {
	return NewObjectFromBytes(s, nil, nil)
}

func NewFileLineFromString(s string) Objecter {
	return NewObjectFromString(s, nil, nil)
}

func NewObjectFromString(s string, leftDelimiter, rightDelimiter []byte) Objecter {
	return NewObjectFromBytes([]byte(s), leftDelimiter, rightDelimiter)
}

func NewObjectFromBytes(d, leftDelimiter, rightDelimiter []byte) Objecter {
	data := append(make([]byte, 0), leftDelimiter...)
	data = append(data, d...)
	data = append(data, rightDelimiter...)
	return &object{data: data}
}

func (o *object) String() string {
	if o == nil || o.data == nil {
		return ""
	}
	return string(o.data)
}

func (o *object) Bytes() []byte {
	if o == nil {
		return nil
	}
	return o.data
}

func (o *object) Write(w io.Writer) (int, error) {
	if o == nil {
		return 0, nil
	}
	return w.Write(o.data)
}

func (o *object) Len() int {
	if o == nil {
		return 0
	}
	return len(o.data)
}
