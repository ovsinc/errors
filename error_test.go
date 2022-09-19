package errors

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

var UnknownErrorType = NewObjectFromString("UNKNOWN_TYPE", nil, nil)

func TestError_SetMsg(t *testing.T) {
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
			name: "nil err",
			err:  nil,
			want: nil,
			args: args{
				"hello",
			},
		},
		{
			name: "nil",
			err:  nil,
			want: nil,
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

func TestError_Error(t *testing.T) {
	type fields struct {
		operation   Objecter
		msg         Objecter
		contextInfo CtxMap
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "nil",
			fields: fields{},
			want:   "",
		},
		{
			name:   "empty",
			fields: fields{},
			want:   "",
		},
		{
			name: "only msg",
			fields: fields{
				msg: NewObjectFromString("hello", nil, nil),
			},
			want: "hello",
		},
		{
			name: "with all params",
			fields: fields{
				operation:   NewOperationFromString("write"),
				contextInfo: CtxMap{"hello": "world", "hi": "there"},
				msg:         NewMsgFromString("hello"),
			},
			want: "[write] {hello:world,hi:there} hello",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				operation:   tt.fields.operation,
				msg:         tt.fields.msg,
				contextInfo: tt.fields.contextInfo,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithOptions(t *testing.T) { //nolint:funlen
	err1 := "hello"

	type fields struct {
		operation   Objecter
		msg         Objecter
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
			fields: fields{},
			want: &Error{
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
			fields: fields{},
			want: &Error{
				msg: NewObjectFromString(err1, nil, nil),
			},
		},
		{
			name: "New. Oper options",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetOperation("read"),
				},
			},
			fields: fields{
				operation: NewObjectFromString("read", _opDelimiterLeft, _opDelimiterRight),
			},
			want: &Error{
				operation: NewObjectFromString("read", _opDelimiterLeft, _opDelimiterRight),
				msg:       NewObjectFromString(err1, nil, nil),
			},
		},
		{
			name: "Set cascade",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetOperation("write"),
				},
			},
			fields: fields{},
			want: New("").WithOptions(
				SetOperation("write"),
				SetMsg(err1),
			),
		},
		{
			name: "Set cascade 2",
			args: args{
				ops: []Options{
					SetMsg(err1),
					SetOperation("write"),
				},
			},
			fields: fields{},
			want: New("").
				WithOptions(
					SetOperation("read"),
				).
				WithOptions(
					SetOperation("write"),
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
				operation:   tt.fields.operation,
				msg:         tt.fields.msg,
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
		want Objecter
		err  *Error
	}{
		{
			name: "New. set",
			err:  NewWith(SetOperation("new operation")),
			want: NewObjectFromString("new operation", _opDelimiterLeft, _opDelimiterRight),
		},
		{
			name: "Set",
			err:  New("").WithOptions(SetOperation("new operation")),
			want: NewObjectFromString("new operation", _opDelimiterLeft, _opDelimiterRight),
		},
		{
			name: "Set 2",
			err: New("").
				WithOptions(SetOperation("new one")).
				WithOptions(SetOperation("new operation")),
			want: NewObjectFromString("new operation", _opDelimiterLeft, _opDelimiterRight),
		},
		{
			name: "Empty",
			err:  New(""),
			want: NewObjectEmpty(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Operation()
			assert.Equal(t, tt.want, got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Error.Operations() = %v, want %v", got, tt.want)
			// }
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
