package errors

import (
	origerrors "errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnwrapByID(t *testing.T) {
	id1 := "myid"
	e1 := NewWith(
		SetMsg("e1"),
		SetID(id1),
	)

	type args struct {
		err error
		id  string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want: nil,
		},
		{
			name: "one",
			args: args{
				err: e1,
				id:  id1,
			},
			want: e1,
		},
		{
			name: "multi",
			args: args{
				err: Combine(
					New("first"), nil, e1,
					New("hello1"), nil,
					NewWith(SetMsg("hello2"), SetID("two"))),
				id: id1,
			},
			want: e1,
		},
		{
			name: "not found",
			args: args{
				err: Combine(New("first"), nil,
					New("hello1"), nil,
					NewWith(SetMsg("hello2"), SetID("two"))),
				id: id1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := UnwrapByID(tt.args.err, tt.args.id)
			if err != nil && tt.want != nil && err.Error() != tt.want.Error() {
				t.Errorf("UnwrapByID() error = %#v, want %#v", err, tt.want)
			}
		})
	}
}

func BenchmarkUnwrapByID(b *testing.B) {
	id1 := "myid"
	e1 := NewWith(
		SetMsg("e1"),
		SetID(id1),
	)

	err := Combine(
		New("first"), e1, New("hello1"), nil,
		NewWith(SetMsg("hello2"), SetID("two")))

	findErr := UnwrapByID(err, id1)
	require.NotNil(b, findErr)
	require.Equal(b, findErr.Error(), e1.Error())

	for i := 0; i < b.N; i++ {
		_ = UnwrapByID(err, id1)
	}
}

func TestIs(t *testing.T) {
	err1 := New("1")

	err2 := New("2")
	err22 := err2

	erra := Combine(err1, nil)
	errb := Combine(nil, err1)

	errc := Wrap(nil, err2)
	errd := Wrap(err2, nil)

	erre := errc

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
		{errc, err2, true},
		{errd, err2, true},
		{erre, errc, true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(t *testing.T) {
			if got := Is(tc.err, tc.target); got != tc.match {
				t.Errorf("Is(%v, %v) = %v, want %v", tc.err, tc.target, got, tc.match)
			}
		})
	}
}

func TestAs(t *testing.T) { //nolint:funlen
	err1 := New("1")

	var err2 error = err1

	var errE1 *Error
	var errE2 errorer

	merr1 := Combine(err1, err2)
	var merr1cast Multierror
	merr2cast := merr1.(*multiError)

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
		{
			err1,
			&errE2,
			true,
			err1,
		},
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
			&merr1cast,
			true,
			merr2cast,
		},
	}
	for i, tc := range testCases {
		tc := tc
		name := fmt.Sprintf("%d:As(Errorf(..., %v), %v)", i, tc.err, tc.target)
		// Clear the target pointer, in case it was set in a previous test.
		rtarget := reflect.ValueOf(tc.target)
		rtarget.Elem().Set(reflect.Zero(reflect.TypeOf(tc.target).Elem()))
		t.Run(name, func(t *testing.T) {
			match := As(tc.err, tc.target)
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
	e2 := New("err two")

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
				err: Wrap(e1, nil),
			},
			wantNil: true,
		},
		{
			name: "std err in wrap",
			args: args{
				err: Wrap(e1, e2),
			},
			want: "err two",
		},
		{
			name: "nil in errs",
			args: args{
				err: Combine(nil, e2, nil, e1),
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := Unwrap(tt.args.err); (err != nil) && !tt.wantNil && tt.want != err.Error() {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.want)
			} else if !tt.wantNil && err == nil {
				t.Errorf("Unwrap() error must be nil byt = %v", err)
			}
		})
	}
}
