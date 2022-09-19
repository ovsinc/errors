package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOnNil(t *testing.T) {
	var err *Error

	assert.Nil(t, err.WithOptions(
		SetMsg("hello"),
		SetContextInfo(CtxMap{"hello": "world"}),
	))
}

func TestNewWith(t *testing.T) { //nolint:funlen
	myop1 := "read"
	msg := "hello"
	myctx := CtxMap{
		"hello": "world",
	}

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
			want: &Error{},
		},
		{
			name: "err with operation",
			args: args{
				ops: []Options{
					SetOperation(myop1),
				},
			},
			want: &Error{
				operation: NewObjectFromString(myop1, _opDelimiterLeft, _opDelimiterRight),
			},
		},
		{
			name: "err with msg",
			args: args{
				ops: []Options{
					SetMsg(msg),
				},
			},
			want: &Error{
				msg: NewObjectFromString(msg, nil, nil),
			},
		},
		{
			name: "err with ctx",
			args: args{
				ops: []Options{
					SetContextInfo(myctx),
				},
			},
			want: &Error{
				contextInfo: myctx,
			},
		},
		{
			name: "err with reset msg",
			args: args{
				ops: []Options{
					SetMsg(msg),
					SetMsg("new"),
				},
			},
			want: &Error{
				msg: NewObjectFromString("new", nil, nil),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWith(tt.args.ops...); got == nil ||
				tt.want == nil ||
				got.Error() != tt.want.Error() {
				t.Errorf("New() = %+v, want %+v.", got, tt.want)
			}
		})
	}
}

//

type stru struct {
	msg string
}

func (a *stru) String() string {
	return a.msg
}

func TestNew(t *testing.T) { //nolint:funlen
	a := &stru{msg: "hello"}

	fn := func() string {
		return "hello"
	}

	helloMsg := NewObjectFromString("hello", nil, nil)

	type unknownType map[int]int

	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "nil",
			args: args{i: nil},
			want: &Error{},
		},
		{
			name: "string",
			args: args{i: "hello"},
			want: &Error{
				msg: helloMsg,
			},
		},
		{
			name: "error",
			args: args{i: New("hello")},
			want: &Error{
				msg: helloMsg,
			},
		},
		{
			name: "String() interface",
			args: args{i: a},
			want: &Error{
				msg: helloMsg,
			},
		},
		{
			name: "func() string",
			args: args{i: fn},
			want: &Error{
				msg: helloMsg,
			},
		},
		{
			name: "unknown",
			args: args{i: unknownType{}},
			want: &Error{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.i); (got == nil || tt.want == nil) ||
				!assert.Equal(t, got.Error(), tt.want.Error()) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
