package errors_test

import (
	"fmt"
	"log"
	"os"
	"time"

	stderrors "errors"

	"github.com/BurntSushi/toml"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ovsinc/multilog/golog"
	"golang.org/x/text/language"

	"github.com/ovsinc/errors"
)

// one err

func ExampleNew() {
	e := errors.New("hello world")

	fmt.Println(e.Error())

	// Output:
	// hello world
}

func ExampleLog() {
	errors.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))

	errors.Log(errors.New("hello world"))

	// Output:
	// ovsinc/errors hello world
}

func ExampleNewLog() {
	errors.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))

	_ = errors.NewLog("hello world")

	// Output:
	// ovsinc/errors hello world
}

func ExampleDownstreamDependencyTimedoutErr() {
	e := errors.DownstreamDependencyTimedoutErr("hello")

	fmt.Printf("%q\n", e)
	fmt.Printf("%v\n", e)

	// Output:
	// id: operation: error_type:DownstreamDependencyTimedout context_info: message:hello
	// (DownstreamDependencyTimedout) hello
}

func ExampleError_Marshal() {
	e := errors.NewWith(
		errors.SetID("myid"),
		errors.SetMsg("hello"),
		errors.AppendContextInfo(
			"hello",
			[]struct{ k, v interface{} }{{"1", 1}, {"10", 11}},
		),
		errors.AppendContextInfo("Joe", "Dow"),
		errors.SetOperation("test op"),
		errors.SetErrorType(errors.InputBody),
	)

	buf, _ := e.Marshal(&errors.MarshalJSON{})

	fmt.Println(string(buf))

	// Output:
	// {"id":"myid","operation":"test op","error_type":"InputBody","context":{"hello":"[{1 1} {10 11}]","Joe":"Dow"},"msg":"hello"}
}

func ExampleNewWith() {
	e := errors.NewWith(
		errors.SetID("myid"),
		errors.SetMsg("hello"),
		errors.AppendContextInfo(
			"hello",
			[]struct{ k, v interface{} }{{"1", 1}, {"10", 11}},
		),
		errors.AppendContextInfo("Joe", "Dow"),
		errors.SetOperation("test op"),
		errors.SetErrorType(errors.InputBody),
	)

	fmt.Printf("%q\n", e)
	fmt.Printf("%v\n", e)

	fmt.Printf("%c\n", e)
	fmt.Printf("%t\n", e)
	fmt.Printf("%o\n", e)
	fmt.Printf("%s\n", e)
	fmt.Printf("%j\n", e)
	fmt.Printf("%f\n", e)

	// Output:
	// id:myid operation:test op error_type:InputBody context_info:hello:[{1 1} {10 11}],Joe:Dow message:hello
	// (InputBody) [test op] {hello:[{1 1} {10 11}],Joe:Dow} hello
	// hello:[{1 1} {10 11}],Joe:Dow
	// InputBody
	// test op
	// hello
	// {"id":"myid","operation":"test op","error_type":"InputBody","context":{"hello":"[{1 1} {10 11}]","Joe":"Dow"},"msg":"hello"}
	// example_test.go:101: ExampleNewWith()
}

func ExampleError_WithOptions() {
	e1 := errors.New("hello")
	e2 := e1.WithOptions(
		errors.AppendContextInfo("hello", "world"),
		errors.SetOperation("test op"),
	)

	fmt.Println(e1.Error())
	fmt.Println(e2.Error())

	// Output:
	// hello
	// [test op] {hello:world} hello
}

// translate

func localizePrepare() *i18n.Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./_examples/translate/testdata/active.ru.toml")

	return i18n.NewLocalizer(bundle, "es", "ru", "en")
}

func localTransContext() errors.TranslateContext {
	var (
		unreadEmailCount = 5
		name             = "John Snow"
	)

	return errors.TranslateContext{
		TemplateData: map[string]interface{}{
			"Name":        name,
			"PluralCount": unreadEmailCount,
		},
		PluralCount: unreadEmailCount,
	}
}

func ExampleTranslate() {
	errEmailsUnreadMsg := localTransContext()
	localizer := localizePrepare()

	e1 := errors.NewWith(
		errors.SetMsg("fallback message"),
		errors.SetID("ErrEmailsUnreadMsg"),
	)

	errors.DefaultLocalizer = localizePrepare()
	defer func() {
		errors.DefaultLocalizer = nil
	}()

	msg, _ := errors.Translate(e1, localizer, &errEmailsUnreadMsg)
	fmt.Println(msg)

	msg, _ = errors.Translate(e1, nil, &errEmailsUnreadMsg)
	fmt.Println(msg)

	eunknown := errors.NewWith(
		errors.SetMsg("fallback unknown message"),
		errors.SetID("ErrUnknownErrorMsg"),
	)

	fmt.Printf("%+s\n", eunknown)

	// Output:
	// У John Snow имеется 5 непрочитанных сообщений.
	// У John Snow имеется 5 непрочитанных сообщений.
	// Неизвестная ошибка
}

// multierror

func itsOk() error {
	return nil
}

func itsErr(s string) error {
	return errors.New(s)
}

func ExampleCombine() {
	e := errors.Combine(
		nil,
		itsOk(),
		itsErr("one"),
		itsErr("two"),
		itsOk(),
	)
	fmt.Println(e)

	// Output:
	// the following errors occurred:
	// 	#1 one
	// 	#2 two
}

// Добавление ошибок в mutierror с логгированием
// тут изменена функция форматирования вывода -- используется json
func ExampleCombineWithLog() {
	errors.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))

	_ = errors.CombineWithLog(
		nil,
		itsOk(),
		itsErr("one"),
		itsErr("two"),
		itsOk(),
	)

	// Output:
	// ovsinc/errors the following errors occurred:
	// 	#1 one
	// 	#2 two
}

func someFuncWithErr() error {
	return errors.NewWith(
		errors.SetMsg("connection error"),
		errors.SetContextInfo(errors.CtxKV{{"hello", "world"}}),
		errors.SetOperation("write"),
		errors.SetID("someid"),
	)
}

func someFuncWithErr2() error {
	return errors.New(
		"connection error",
	)
}

func ExampleWrap() {
	e := stderrors.New("hello world")

	err := someFuncWithErr()
	err = errors.Wrap(err, someFuncWithErr2())
	err = errors.Wrap(err, e)

	fmt.Printf("%v", err)

	oneErr := errors.FindByID(err, "someid")

	fmt.Printf(
		"err with id 'someid': %v; id: %s; op: %s\n",
		oneErr,
		errors.GetID(oneErr),
		errors.GetOperation(oneErr),
	)

	fmt.Println(errors.FindByErr(err, e))

	// Output:
	// the following errors occurred:
	//	#1 [write] {hello:world} connection error
	//	#2 connection error
	//	#3 hello world
	// err with id 'someid': [write] {hello:world} connection error; id: someid; op: write
	// hello world
}

func ExampleNewWithLog() {
	errors.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))

	_ = errors.Combine(
		nil,
		itsOk(),
		errors.NewLog("one"),
		errors.NewLog("two"),
		itsOk(),
	)

	_ = errors.NewLog("three")

	// Output:
	// ovsinc/errors one
	// ovsinc/errors two
	// ovsinc/errors three
}

func someErrWithTimedCall() error {
	err := errors.NewWith(
		errors.SetMsg("some call"),
		errors.AppendContextInfo("call", errors.RuntimeCaller()),
		errors.AppendContextInfo("duration", time.Second),
	)
	time.Sleep(1 * time.Second)
	return err
}

func ExampleCaller() {
	err := someErrWithTimedCall()

	fmt.Printf("%s\n", err.Error())

	// Output:
	// {call:example_test.go:306: ExampleCaller(),duration:1s} some call
}
