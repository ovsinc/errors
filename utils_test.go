package errors

import (
	"testing"
)

func CallHa() Objecter {
	return HandlerCaller()
}

func TestHandlerCaller(t *testing.T) {
	i := CallHa()
	if i.String() != "utils_test.go:8: CallHa()" {
		t.Errorf("Caller() = %v, want %v", i.String(), "utils_test.go:8: CallHa()")
	}
}
