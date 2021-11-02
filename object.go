package errors

import "bytes"

type Objecter interface {
	String() string
	Bytes() []byte
	Buffer() *bytes.Buffer
}

type object []byte

func NewObjectEmpty() Objecter {
	return object{}
}

func NewObjectFromBytes(v []byte) Objecter {
	return object(v)
}

func NewObjectFromString(s string) Objecter {
	return object(s)
}

func (o object) String() string {
	if len(o) == 0 {
		return ""
	}
	return string(o)
}

func (o object) Bytes() []byte {
	if len(o) == 0 {
		return []byte{}
	}
	return o
}

func (o object) Buffer() *bytes.Buffer {
	if len(o) == 0 {
		return &bytes.Buffer{}
	}
	return bytes.NewBuffer(o)
}
