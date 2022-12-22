package errors

import (
	"io"
)

type Marshaller interface {
	Marshal(interface{}) ([]byte, error)
	MarshalTo(interface{}, io.Writer) error
}

var DefaultMarshaller = &MarshalString{} //nolint:gochecknoglobals
