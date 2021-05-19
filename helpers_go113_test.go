package errors_test

import (
	"fmt"
	"reflect"
	"testing"

	origerrors "errors"

	hashmultierr "github.com/hashicorp/go-multierror"
	"github.com/ovsinc/errors"
)

func TestIs(t *testing.T) {
	err1 := errors.New("1")
	err11 := errors.New("1")
	err2 := errors.New("2")
	err22 := err2
	erra := hashmultierr.Append(err1, err2)
	err3 := errors.New("3")
	errb := errors.Wrap(erra, err2)
	errc := errors.Append(err1, err3, err2)

	testCases := []struct {
		err    error
		target error
		match  bool
	}{
		{nil, nil, true},
		{err1, nil, false},
		{err1, err1, true},
		{err2, err22, true},
		{erra, err1, true},
		{errb, err1, true},
		{errc, err1, true},
		{err1, err11, false},
		{err1, err3, false},
		{erra, err3, false},
		{errb, err3, false},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.match {
				t.Errorf("Is(%v, %v) = %v, want %v", tc.err, tc.target, got, tc.match)
			}
		})
	}
}

func TestAs(t *testing.T) {
	err1 := errors.New("1")
	merr1 := errors.Append(err1)

	var merr1cast errors.Multierror
	origerrors.As(merr1, &merr1cast)

	var err2 error = err1

	var errE1 *errors.Error
	// var errE2 errors.Errorer
	var merrE1 errors.Multierror

	testCases := []struct {
		err    error
		target interface{}
		match  bool
		want   interface{} // value of target on match
	}{
		{
			nil,
			&err1,
			false,
			nil,
		},
		// {
		// 	err1,
		// 	&errE2,
		// 	true,
		// 	err1,
		// },
		{
			err1,
			&errE1,
			true,
			err1,
		},
		{
			err2,
			&errE1,
			true,
			err1,
		},
		{
			merr1,
			&merrE1,
			true,
			merr1cast,
		},
	}
	for i, tc := range testCases {
		tc := tc
		name := fmt.Sprintf("%d:As(Errorf(..., %v), %v)", i, tc.err, tc.target)
		// Clear the target pointer, in case it was set in a previous test.
		rtarget := reflect.ValueOf(tc.target)
		rtarget.Elem().Set(reflect.Zero(reflect.TypeOf(tc.target).Elem()))
		t.Run(name, func(t *testing.T) {
			match := errors.As(tc.err, tc.target)
			if match != tc.match {
				t.Fatalf("match: got %v; want %v", match, tc.match)
			}
			if !match {
				return
			}
			if got := rtarget.Elem().Interface(); got != tc.want {
				t.Fatalf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	e1 := origerrors.New("hello")
	e2 := errors.New("err two")

	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		want    string
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			wantNil: true,
		},
		{
			name: "one in wrap",
			args: args{
				err: errors.Wrap(e1, nil),
			},
			wantNil: true,
		},
		{
			name: "std err in wrap",
			args: args{
				err: errors.Wrap(e1, e2),
			},
			want: "hello",
		},
		{
			name: "err in wrap",
			args: args{
				err: errors.Wrap(e2, e1),
			},
			want: "err two",
		},
		{
			name: "nil in errs",
			args: args{
				err: errors.Append(nil, e2, nil, e1),
			},
			want: "err two",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := errors.Unwrap(tt.args.err); (err != nil) && !tt.wantNil && tt.want != err.Error() {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.want)
			} else if !tt.wantNil && err == nil {
				t.Errorf("Unwrap() error must be nil byt = %v", err)
			}
		})
	}
}
