package errors

import (
	"testing"
)

func CallHa() CallInfo {
	return HandlerCaller()
}

func TestHandlerCaller(t *testing.T) {
	i := CallHa()
	if i.FuncName != "CallHa" {
		t.Errorf("Caller() = %v, want %v", i.FuncName, "CallTe")
	}
}
