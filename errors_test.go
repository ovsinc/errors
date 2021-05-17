package errors

import (
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

var UnknownErrorType = "UNKNOWN_TYPE"

func TestNewNil(t *testing.T) {
	var err *Error

	assert.Nil(t, err.WithOptions(
		SetErrorType(UnknownErrorType),
		SetSeverity(SeverityError),
		SetMsg("hello"),
		SetContextInfo(CtxMap{"hello": "world"}),
	))
}

func TestNew(t *testing.T) {
	myerr1 := "some err"
	myerrType1 := "custom err type"
	myop1 := "read"
	myseverity := SeverityError

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
					SetMsg(myerr1),
				},
			},
			want: &Error{
				msg:        myerr1,
				errorType:  myerrType1,
				operations: []string{myop1},
				severity:   myseverity,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := New("", tt.args.ops...); got != nil && tt.want != nil && got.Error() != tt.want.Error() {
				t.Errorf("New() = %+v, want %+v.", got, tt.want)
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

			if tt.err != nil && tt.want != nil && tt.err.Error() != tt.want.Error() {
				t.Errorf("SetMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetFormatFn(t *testing.T) {
	myerr := &Error{}

	var testFormatFn FormatFn = func(w io.Writer, e *Error) {}

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
				DefaultFormatFn,
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

			if tt.err != nil && tt.want != nil && tt.err.Error() != tt.want.Error() {
				t.Errorf("SetFormatFn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		operations  []string
		errorType   string
		msg         string
		severity    Severity
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
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
				msg:        "",
			},
			want: "(UNKNOWN_TYPE)",
		},
		{
			name: "empty",
			fields: fields{
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
				msg:        "hello",
			},
			want: "(UNKNOWN_TYPE) -- hello",
		},
		{
			name: "with all params",
			fields: fields{
				operations:  []string{"write"},
				severity:    SeverityError,
				errorType:   "not found",
				msg:         "hello",
				contextInfo: CtxMap{"hello": "world", "hi": "there"},
			},
			want: "(not found)[write]<hello:world,hi:there> -- hello",
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
		operations  []string
		errorType   string
		msg         string
		severity    Severity
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
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
			},
			want: &Error{
				operations:  make([]string, 0),
				severity:    SeverityError,
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
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
			},
			want: &Error{
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
				msg:        err1,
			},
		},
		{
			name: "New. Oper options",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetErrorType("my type"),
					SetOperations("write", "read"),
					SetSeverity(SeverityWarn),
				},
			},
			fields: fields{
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
			},
			want: &Error{
				operations: []string{"write", "read"},
				severity:   SeverityWarn,
				errorType:  "my type",
				msg:        err1,
			},
		},
		{
			name: "Set cascade",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetErrorType("my type"),
					SetOperations("write", "read"),
					SetSeverity(SeverityWarn),
				},
			},
			fields: fields{
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
			},
			want: New("").WithOptions(
				SetOperations("write", "read"),
				SetSeverity(SeverityWarn),
				SetErrorType("my type"),
				SetMsg(err1),
			),
		},
		{
			name: "Set cascade 2",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetErrorType("my type"),
					SetOperations("write", "read"),
					SetSeverity(SeverityWarn),
				},
			},
			fields: fields{
				operations: make([]string, 0),
				severity:   SeverityError,
				errorType:  UnknownErrorType,
			},
			want: New("").
				WithOptions(
					SetOperations("write", "read"),
				).
				WithOptions(
					SetSeverity(SeverityWarn),
				).
				WithOptions(
					SetErrorType("my type"),
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
		want []string
		err  *Error
	}{
		{
			name: "New. set",
			err:  New("", SetOperations("new operation")),
			want: []string{"new operation"},
		},
		{
			name: "Set",
			err:  New("").WithOptions(SetOperations("new operation")),
			want: []string{"new operation"},
		},
		{
			name: "Set 2",
			err: New("").
				WithOptions(SetOperations("noe one")).
				WithOptions(AppendOperations()).
				WithOptions(SetOperations("new operation")),
			want: []string{"new operation"},
		},
		{
			name: "append",
			err: New("").
				WithOptions(SetOperations("new operation")).
				WithOptions(AppendOperations("noe one")),
			want: []string{"new operation", "noe one"},
		},
		{
			name: "Empty",
			err:  New(""),
			want: nil,
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
		want string
	}{
		{
			name: "empty",
			err:  &Error{},
			want: "",
		},
		{
			name: "New. Empty ",
			err:  New(""),
			want: "",
		},
		{
			name: "New. Set",
			err:  New("", SetErrorType("my type")),
			want: "my type",
		},
		{
			name: "Set",
			err:  New("").WithOptions(SetErrorType("my type")),
			want: "my type",
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
		want Severity
	}{
		{
			name: "empty",
			err:  &Error{},
			want: SeverityUnknown,
		},
		{
			name: "New",
			err:  New(""),
			want: SeverityError,
		},
		{
			name: "New. Set",
			err:  New("", SetSeverity(SeverityWarn)),
			want: SeverityWarn,
		},
		{
			name: "Set",
			err:  New("").WithOptions(SetSeverity(SeverityWarn)),
			want: SeverityWarn,
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
	var emptyErr *Error = &Error{}
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
	var mynil *Error = &Error{}
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

			if err != nil && tt.want != nil && err.Error() != tt.want.Error() {
				t.Errorf("Error.ErrorOrNil() error = _%v_, want _%v_", err, tt.want)
			}
		})
	}
}
