package errors_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ovsinc/errors"
	"github.com/ovsinc/multilog/golog"
	"golang.org/x/text/language"
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

	errors.NewLog("hello world")

	// Output:
	// ovsinc/errors hello world
}

func ExampleNewWith() {
	e := errors.NewWith(
		errors.SetID("myid"),
		errors.SetMsg("hello"),
		errors.AppendContextInfo("hello", "world"),
		errors.AppendContextInfo("Joe", "Dow"),
		errors.SetOperation("test op"),
		errors.SetErrorType("mytype"),
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
	// id:myid operation:test op errorType:mytype contextInfo:map[Joe:Dow hello:world] msg:hello
	// (mytype) [test op] {Joe:Dow,hello:world} hello
	// Joe:Dow,hello:world
	// mytype
	// test op
	// hello
	// {"id":"myid","operation":"test op","context":{"Joe":"Dow","hello":"world"},"msg":"hello"}
	// example_test.go:63: ExampleNewWith()
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
		errors.SetContextInfo(errors.CtxMap{"hello": "world"}),
		errors.SetOperation("write"),
	)
}

func someFuncWithErr2() error {
	return errors.New(
		"connection error",
	)
}

func ExampleWrap() {
	err := someFuncWithErr()

	err = errors.Wrap(err, someFuncWithErr2())

	fmt.Printf("%v\n", err)

	// Output:
	// the following errors occurred:
	// 	#1 [write] {hello:world} connection error
	// 	#2 connection error
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

func someErrWithTimedCall() (err *errors.Error) {
	begin := time.Now()
	defer func() {
		err = err.WithOptions(
			errors.AppendContextInfo("duration", time.Since(begin).Round(time.Second)),
			errors.AppendContextInfo("call", errors.DefaultCaller()),
		)
	}()

	err = errors.New("some call")

	time.Sleep(1 * time.Second)

	return err
}

func ExampleCaller() {
	err := someErrWithTimedCall()

	fmt.Printf("%s\n", err.Error())

	// Output:
	// {call:example_test.go:263: ExampleCaller(),duration:1s} some call
}
