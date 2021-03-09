package errors

import (
	"testing"

	"gitlab.com/ovsinc/errors/log"
)

func TestLog(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.args.e)
		})
	}
}

func TestError_Log(t *testing.T) {
	type fields struct {
		severity    log.Severity
		contextInfo CtxMap
		operations  []Operation
		formatFn    FormatFn
		errorType   ErrorType
		msg         string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				severity:    tt.fields.severity,
				contextInfo: tt.fields.contextInfo,
				operations:  tt.fields.operations,
				formatFn:    tt.fields.formatFn,
				errorType:   tt.fields.errorType,
				msg:         tt.fields.msg,
			}
			e.Log()
		})
	}
}
