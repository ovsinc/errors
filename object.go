package errors

import "io"

// type Objecter interface {
// 	String() string
// 	Bytes() []byte
// 	Buffer() *bytes.Buffer
// }

type Object []byte

func NewObjectEmpty() Object {
	return Object{}
}

func NewObjectFromBytes(v []byte) Object {
	return Object(v)
}

func NewObjectFromString(s string) Object {
	return Object(s)
}

func (o Object) String() string {
	if o == nil {
		return ""
	}
	return string(o)
}

func (o Object) Bytes() []byte {
	if o == nil {
		return []byte{}
	}
	return o
}

func (o Object) Write(w io.Writer) (int, error) {
	if o == nil {
		return 0, nil
	}
	return w.Write(o)
}
