package errors_test

import (
	"testing"

	"github.com/ovsinc/errors"
	"github.com/stretchr/testify/require"
)

var (
	je1 = errors.NewWith(
		errors.SetMsg("hello1"),
		errors.SetOperation("write"),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "my": "name"}),
	)

	je2 = errors.NewWith(
		errors.SetMsg("hello2"),
		errors.SetOperation("read"),
		errors.SetContextInfo(errors.CtxMap{"hello2": "world", "my2": "name"}),
	)

	je3 = errors.NewWith(
		errors.SetMsg("hello3"),
		errors.SetOperation("read"),
		errors.SetContextInfo(errors.CtxMap{"hello3": "world", "my3": "name"}),
	)
)

func BenchmarkJsonFn(b *testing.B) {
	e := errors.NewWith(
		errors.SetMsg("hello"),
		errors.SetOperation("write"),
		errors.SetContextInfo(errors.CtxMap{"hello": "world", "hi": "there"}),
	)

	data, err := e.Marshal(&errors.MarshalJSON{})
	require.Nil(b, err)

	require.JSONEq(b, string(data), "{\"id\":\"\",\"operation\":\"write\",\"context\":{\"hello\":\"world\",\"hi\":\"there\"},\"msg\":\"hello\"}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat3Errs(b *testing.B) {
	e := errors.Combine(je1, je2, je3)

	marshal := &errors.MarshalJSON{}
	data, err := marshal.Marshal(e)
	require.Nil(b, err)

	require.JSONEq(b, string(data), "{\"count\":3,\"messages\":[{\"id\":\"\",\"operation\":\"write\",\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"id\":\"\",\"operation\":\"read\",\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"},{\"id\":\"\",\"operation\":\"read\",\"context\":{\"hello3\":\"world\",\"my3\":\"name\"},\"msg\":\"hello3\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat2Errs(b *testing.B) {
	e := errors.Wrap(je1, je2)

	marshal := &errors.MarshalJSON{}
	data, err := marshal.Marshal(e)
	require.Nil(b, err)

	require.JSONEq(b, string(data), "{\"count\":2,\"messages\":[{\"id\":\"\",\"operation\":\"write\",\"context\":{\"hello\":\"world\",\"my\":\"name\"},\"msg\":\"hello1\"},{\"id\":\"\",\"operation\":\"read\",\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkJsonMultierrFuncFormat1Err(b *testing.B) {
	e := errors.Wrap(nil, je2)

	marshal := &errors.MarshalJSON{}
	data, err := marshal.Marshal(e)
	require.Nil(b, err)

	require.JSONEq(b, string(data), "{\"count\":1,\"messages\":[{\"id\":\"\",\"operation\":\"read\",\"context\":{\"hello2\":\"world\",\"my2\":\"name\"},\"msg\":\"hello2\"}]}")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}
