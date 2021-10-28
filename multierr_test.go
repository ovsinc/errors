package errors_test

import (
	"testing"

	"github.com/ovsinc/errors"
	"github.com/stretchr/testify/assert"

	origerrors "errors"
)

var (
	me1 = errors.New(
		"hello1",
		errors.SetErrorType("not found"),
		errors.SetOperations("write"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	me2 = errors.New(
		"hello2",
		errors.SetErrorType("not found"),
		errors.SetOperations("read"),
		errors.SetSeverity(errors.SeverityError),
		errors.SetContextInfo(errors.CtxMap{"hello2": "world", "my2": "name"}),
	)

	errMe3 = origerrors.New("hello")
)

func TestWrapSimple(t *testing.T) {
	type args struct {
		left  error
		right error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantNil bool
	}{
		{
			name: "nil",
			args: args{
				left:  nil,
				right: nil,
			},
			wantNil: true,
		},
		{
			name: "nil left",
			args: args{
				left:  nil,
				right: me1,
			},
			want: "(not found)[write]{hello:world,my:name} -- hello1",
		},
		{
			name: "nil left std",
			args: args{
				left:  nil,
				right: errMe3,
			},
			want: "hello",
		},
		{
			name: "nil right",
			args: args{
				left:  me2,
				right: nil,
			},
			want: "(not found)[read]{hello2:world,my2:name} -- hello2",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if err := errors.Wrap(tt.args.left, tt.args.right); (err != nil) && !tt.wantNil && !assert.Equal(t, tt.want, err.Error()) {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.want)
			} else if err == nil && !tt.wantNil {
				t.Errorf("Wrap() error must be nill by = %v", err)
			}
		})
	}
}

func TestWrapMultierr(t *testing.T) {
	type args struct {
		left  error
		right error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantNil bool
	}{
		{
			name: "two",
			args: args{
				left:  me2,
				right: me1,
			},
			want: "the following errors occurred:\n\t#1 (not found)[read]{hello2:world,my2:name} -- hello2\n\t#2 (not found)[write]{hello:world,my:name} -- hello1\n",
		},
		{
			name: "two std",
			args: args{
				left:  errMe3,
				right: me1,
			},
			want: "the following errors occurred:\n\t#1 hello\n\t#2 (not found)[write]{hello:world,my:name} -- hello1\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := errors.Wrap(tt.args.left, tt.args.right); (err != nil) && !tt.wantNil && tt.want != err.Error() {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.want)
			} else if err == nil && !tt.wantNil {
				t.Errorf("Wrap() error must be nill by = %v", err)
			}
		})
	}
}

func TestCombine(t *testing.T) {
	type args struct {
		errors []error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantNil bool
	}{
		{
			name: "nil",
			args: args{
				errors: nil,
			},
			wantNil: true,
		},
		{
			name: "one",
			args: args{
				errors: []error{me1},
			},
			want: "the following errors occurred:\n\t#1 (not found)[write]{hello:world,my:name} -- hello1\n",
		},
		{
			name: "many with nil",
			args: args{
				errors: []error{nil, me1, nil, se2, nil, se3, nil},
			},
			want: "the following errors occurred:\n\t#1 (not found)[write]{hello:world,my:name} -- hello1\n\t#2 (not found)[read]{hello2:world,my2:name} -- hello2\n\t#3 (not found)[read]{hello3:world,my3:name} -- hello3\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := errors.Combine(tt.args.errors...); (err != nil) && !tt.wantNil && tt.want != err.Error() {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.want)
			} else if !tt.wantNil && err == nil {
				t.Errorf("Wrap() error must be nill by = %v", err)
			}
		})
	}
}
