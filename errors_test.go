package errors

import (
	origerrors "errors"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ovsinc/errors/log"
)

func TestNewNil(t *testing.T) {
	var err *Error

	assert.Nil(t, err.WithOptions(
		SetErrorType(UnknownErrorType),
		SetOperations(Operation("")),
		SetSeverity(log.SeverityError),
		SetMsg("hello"),
		SetContextInfo(CtxMap{"hello": "world"}),
	))
}

func TestNew(t *testing.T) {
	myerr1 := "some err"
	myerrType1 := NewErrorType("custom err type")
	myop1 := NewOperation("read")
	myseverity := log.SeverityError

	type args struct {
		ops []Options
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "empty",
			args: args{
				ops: []Options{},
			},
			want: &Error{},
		},
		{
			name: "With err, error type, operation, severity",
			args: args{
				ops: []Options{
					SetErrorType(myerrType1),
					SetOperations(myop1),
					SetSeverity(myseverity),
				},
			},
			want: &Error{
				msg:        myerr1,
				errorType:  myerrType1,
				operations: []Operation{myop1},
				severity:   myseverity,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := New("", tt.args.ops...); origerrors.Is(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMsg(t *testing.T) {
	myerr := &Error{}

	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		err  *Error
		want *Error
	}{
		{
			name: "nil",
			err:  nil,
			want: nil,
			args: args{
				"hello",
			},
		},
		{
			name: "simple",
			err:  myerr,
			want: myerr,
			args: args{
				"hello",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := SetMsg(tt.args.msg)
			got(tt.err)

			if !origerrors.Is(tt.err, tt.want) {
				t.Errorf("SetMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetFormatFn(t *testing.T) {
	myerr := &Error{}

	var testFormatFn FormatFn = func(e *Error) string {
		return ""
	}

	type args struct {
		fn FormatFn
	}
	tests := []struct {
		name string
		args args
		err  *Error
		want *Error
	}{
		{
			name: "nil",
			err:  nil,
			want: nil,
			args: args{
				defaultFormatFn,
			},
		},
		{
			name: "simple",
			err:  myerr,
			want: myerr,
			args: args{
				testFormatFn,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := SetFormatFn(tt.args.fn)
			got(tt.err)

			if !origerrors.Is(tt.err, tt.want) {
				t.Errorf("SetFormatFn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		operations  []Operation
		errorType   ErrorType
		msg         string
		severity    log.Severity
		contextInfo CtxMap
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil",
			fields: fields{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
				msg:        "",
			},
			want: "[UNKNOWN_TYPE][ERROR]",
		},
		{
			name: "empty",
			fields: fields{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
				msg:        "hello",
			},
			want: "[UNKNOWN_TYPE][ERROR] -- hello",
		},
		{
			name: "with all params",
			fields: fields{
				operations:  []Operation{NewOperation("write")},
				severity:    log.SeverityError,
				errorType:   NewErrorType("not found"),
				msg:         "hello",
				contextInfo: CtxMap{"hello": "world", "hi": "there"},
			},
			want: "[not found][ERROR][write]<hello:world,hi:there> -- hello",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				operations:  tt.fields.operations,
				errorType:   tt.fields.errorType,
				msg:         tt.fields.msg,
				severity:    tt.fields.severity,
				contextInfo: tt.fields.contextInfo,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithOptions(t *testing.T) {
	err1 := "hello"

	type fields struct {
		operations  []Operation
		errorType   ErrorType
		msg         string
		severity    log.Severity
		contextInfo CtxMap
	}
	type args struct {
		ops []Options
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "New. context",
			args: args{
				ops: []Options{
					SetContextInfo(CtxMap{"duration": time.Second}),
				},
			},
			fields: fields{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
			},
			want: &Error{
				operations:  make([]Operation, 0),
				severity:    log.SeverityError,
				errorType:   UnknownErrorType,
				contextInfo: CtxMap{"duration": time.Second},
			},
		},
		{
			name: "New. msg",
			args: args{
				ops: []Options{
					SetMsg(err1),
				},
			},
			fields: fields{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
			},
			want: &Error{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
				msg:        err1,
			},
		},
		{
			name: "New. Oper options",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetErrorType(NewErrorType("my type")),
					SetOperations(NewOperation("write"), NewOperation("read")),
					SetSeverity(log.SeverityWarn),
				},
			},
			fields: fields{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
			},
			want: &Error{
				operations: []Operation{NewOperation("write"), NewOperation("read")},
				severity:   log.SeverityWarn,
				errorType:  NewErrorType("my type"),
				msg:        err1,
			},
		},
		{
			name: "Set cascade",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetErrorType(NewErrorType("my type")),
					SetOperations(NewOperation("write"), NewOperation("read")),
					SetSeverity(log.SeverityWarn),
				},
			},
			fields: fields{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
			},
			want: New("").WithOptions(
				SetOperations(NewOperation("write"), NewOperation("read")),
				SetSeverity(log.SeverityWarn),
				SetErrorType(NewErrorType("my type")),
				SetMsg(err1),
			),
		},
		{
			name: "Set cascade 2",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetErrorType(NewErrorType("my type")),
					SetOperations(NewOperation("write"), NewOperation("read")),
					SetSeverity(log.SeverityWarn),
				},
			},
			fields: fields{
				operations: make([]Operation, 0),
				severity:   log.SeverityError,
				errorType:  UnknownErrorType,
			},
			want: New("").
				WithOptions(
					SetOperations(NewOperation("write"), NewOperation("read")),
				).
				WithOptions(
					SetSeverity(log.SeverityWarn),
				).
				WithOptions(
					SetErrorType(NewErrorType("my type")),
				).
				WithOptions(
					SetMsg(err1),
				),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				operations:  tt.fields.operations,
				errorType:   tt.fields.errorType,
				msg:         tt.fields.msg,
				severity:    tt.fields.severity,
				contextInfo: tt.fields.contextInfo,
			}
			if got := e.WithOptions(tt.args.ops...); got.Error() != tt.want.Error() {
				t.Errorf("Error.WithOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Operations(t *testing.T) {
	tests := []struct {
		name string
		want []Operation
		err  *Error
	}{
		{
			name: "New. set",
			err:  New("", SetOperations(NewOperation("new operation"))),
			want: []Operation{NewOperation("new operation")},
		},
		{
			name: "Set",
			err:  New("").WithOptions(SetOperations(NewOperation("new operation"))),
			want: []Operation{NewOperation("new operation")},
		},
		{
			name: "Set 2",
			err: New("").
				WithOptions(SetOperations(NewOperation("noe one"))).
				WithOptions(AppendOperations()).
				WithOptions(SetOperations(NewOperation("new operation"))),
			want: []Operation{NewOperation("new operation")},
		},
		{
			name: "append",
			err: New("").
				WithOptions(SetOperations(NewOperation("new operation"))).
				WithOptions(AppendOperations(NewOperation("noe one"))),
			want: []Operation{NewOperation("new operation"), NewOperation("noe one")},
		},
		{
			name: "Empty",
			err:  New(""),
			want: []Operation{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Operations()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.Operations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ErrorType(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want ErrorType
	}{
		{
			name: "empty",
			err:  &Error{},
			want: "",
		},
		{
			name: "New. Empty ",
			err:  New(""),
			want: UnknownErrorType,
		},
		{
			name: "New. Set",
			err:  New("", SetErrorType(NewErrorType("my type"))),
			want: NewErrorType("my type"),
		},
		{
			name: "Set",
			err:  New("").WithOptions(SetErrorType(NewErrorType("my type"))),
			want: NewErrorType("my type"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.ErrorType(); got != tt.want {
				t.Errorf("Error.ErrorType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Severity(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want log.Severity
	}{
		{
			name: "empty",
			err:  &Error{},
			want: log.SeverityUnknown,
		},
		{
			name: "New",
			err:  New(""),
			want: log.SeverityError,
		},
		{
			name: "New. Set",
			err:  New("", SetSeverity(log.SeverityWarn)),
			want: log.SeverityWarn,
		},
		{
			name: "Set",
			err:  New("").WithOptions(SetSeverity(log.SeverityWarn)),
			want: log.SeverityWarn,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Severity(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.Severity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Sdump(t *testing.T) {
	emptyErr := &Error{}
	var mynil *Error
	e1 := New("")
	e2 := New("hello")

	tests := []struct {
		name string
		want string
		err  *Error
	}{
		{
			name: "empty",
			err:  emptyErr,
			want: spew.Sdump(emptyErr),
		},
		{
			name: "empty. new",
			err:  e1,
			want: spew.Sdump(e1),
		},
		{
			name: "empty. new with ops",
			err:  e2,
			want: spew.Sdump(e2),
		},
		{
			name: "nil",
			err:  mynil,
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Sdump(); got != tt.want {
				t.Errorf("Error.Sdump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ErrorOrNil(t *testing.T) {
	var mynil *Error
	mye1 := New("")

	tests := []struct {
		err     *Error
		name    string
		want    *Error
		wantNil bool
	}{
		{
			name:    "nil",
			err:     mynil,
			want:    nil,
			wantNil: true,
		},
		{
			name:    "not nil",
			err:     mye1,
			want:    mye1,
			wantNil: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.err.ErrorOrNil()

			if tt.wantNil {
				if err != nil {
					t.Errorf("Error.ErrorOrNil() want nil but do not")
				}
				return
			}

			if err.Error() != tt.want.Error() {
				t.Errorf("Error.ErrorOrNil() error = _%v_, want _%v_", err, tt.want)
			}
		})
	}
}
