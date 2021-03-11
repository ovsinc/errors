package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnwrapWithID(t *testing.T) {
	id1 := "myid"
	e1 := New(
		"e1",
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
				err: Append(New("first"), nil, e1, New("hello1"), nil, New("hello2", SetID("two"))),
				id:  id1,
			},
			want: e1,
		},
		{
			name: "not found",
			args: args{
				err: Append(New("first"), nil, New("hello1"), nil, New("hello2", SetID("two"))),
				id:  id1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := UnwrapWithID(tt.args.err, tt.args.id)
			if !Is(err, tt.want) {
				t.Errorf("UnwrapWithID() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkUnwrapWithID(b *testing.B) {
	id1 := "myid"
	e1 := New(
		"e1",
		SetID(id1),
	)

	err := Append(New("first"), e1, New("hello1"), nil, New("hello2", SetID("two")))

	findErr := UnwrapWithID(err, id1)
	require.NotNil(b, findErr)
	require.Equal(b, findErr.Error(), e1.Error())

	for i := 0; i < b.N; i++ {
		_ = UnwrapWithID(err, id1)
	}
}
