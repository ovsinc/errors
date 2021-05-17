package errors_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ovsinc/multilog"
	"github.com/ovsinc/multilog/golog"
	"gitlab.com/ovsinc/errors"
	"golang.org/x/text/language"
)

func itsOk() error {
	return nil
}

func itsErr(s string) error {
	return errors.New(s)
}

// Накопление результатов выполнения каких-либо функций спроверкой на НЕ nil
func ExampleErrorOrNil() {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	var err error
	err = errors.Append(
		nil,
		itsOk(),
		itsErr("one"),
		errors.New("two", errors.SetSeverity(errors.SeverityWarn)),
		itsOk(),
	)
	err = errors.Wrap(err, errors.New("three", errors.SetSeverity(errors.SeverityWarn)))

	fmt.Printf("%v\n", errors.ErrorOrNil(err))
	// Output:
	// one
}

// Добавление ошибок в mutierror с логгированием
// тут изменена функция форматирования вывода -- испльзуется json
func ExampleAppendWithLog() {
	multilog.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat

	_ = errors.AppendWithLog(
		nil,
		itsOk(),
		itsErr("one"),
		itsErr("two"),
		itsOk(),
	)

	// Output:
	// ovsinc/errors ERR: {"count":2,"messages":[{"id":"","error_type":"","severity":"ERROR","operations":[],"context":null,"msg":"one"},{"id":"","error_type":"","severity":"ERROR","operations":[],"context":null,"msg":"two"}]}
}

func someFuncWithErr() error {
	return errors.New(
		"connection error",
		errors.SetContextInfo(errors.CtxMap{"hello": "world"}),
		errors.AppendOperations("write"),
		errors.SetSeverity(errors.SeverityUnknown),
		errors.SetErrorType(""),
	)
}

func someFuncWithErr2() error {
	return errors.New(
		"connection error",
		errors.SetSeverity(errors.SeverityUnknown),
		errors.SetErrorType(""),
	)
}

func ExampleWrap() {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := someFuncWithErr()

	err = errors.Wrap(err, someFuncWithErr2())

	fmt.Printf("%v\n", err)

	// Output:
	// the following errors occurred:
	// * [write]<hello:world> -- connection error
	// * connection error
}

func ExampleNewWithLog() {
	multilog.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	_ = errors.Append(
		nil,
		itsOk(),
		errors.NewWithLog("one"),
		errors.NewWithLog("two"),
		itsOk(),
	)

	_ = errors.NewWithLog("three")

	// Output:
	// ovsinc/errors ERR: one
	// ovsinc/errors ERR: two
	// ovsinc/errors ERR: three
}

func someErrFunc() error {
	return errors.New("connection error", errors.SetErrorType("NOT_FOUND"))
}

func ExampleGetErrorType() {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := someErrFunc()

	switch errors.GetErrorType(err) {
	case "NOT_FOUND":
		fmt.Printf("Got error with type NOT_FOUND")
	case errors.UnknownErrorType:
		fmt.Printf("Got error with type %s", errors.UnknownErrorType)
	default:
		fmt.Printf("Got some unknown")
	}

	// Output:
	// Got error with type NOT_FOUND
}

func someTimedCast() (err error) {
	begin := time.Now()
	defer func() {
		err = errors.Cast(err).
			WithOptions(errors.SetContextInfo(
				errors.CtxMap{
					"duration": time.Since(begin).Round(time.Second),
					"call":     errors.DefaultCaller(),
				},
			))
	}()

	err = errors.New("some call")

	time.Sleep(1 * time.Second)

	return err
}

func ExampleLog() {
	multilog.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	errors.Log(someTimedCast())

	// Output:
	// ovsinc/errors ERR: <call:example_test.go:162,duration:1s> -- some call
}

func localizePrepare() *i18n.Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./testdata/active.ru.toml")

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

func ExampleError_WithOptions() {
	e1 := errors.New("hello")

	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			_ = e1.WithOptions(
				errors.SetMsg("new error " + strconv.Itoa(i)),
			)
		}()
	}

	wg.Wait()

	e2 := e1.WithOptions(
		errors.AppendContextInfo("hello", "world"),
		errors.AppendOperations("test op"),
		errors.SetErrorType("test type"),
	)

	fmt.Println(e1.Error())
	fmt.Println(e2.Error())

	// Output:
	// hello
	// (test type)[test op]<hello:world> -- hello
}

func ExampleError_TranslateMsg() {
	errEmailsUnreadMsg := localTransContext()
	localizer := localizePrepare()

	e1 := errors.New(
		"fallback message",
		errors.SetID("ErrEmailsUnreadMsg"),
		errors.SetErrorType("not found"),
		errors.SetTranslateContext(&errEmailsUnreadMsg),
		errors.SetLocalizer(localizer),
	)

	fmt.Printf("%v\n", e1)
	fmt.Println(e1.Error())
	fmt.Print(e1.TranslateMsg())

	// Output:
	// (not found) -- У John Snow имеется 5 непрочитанных сообщений.
	// (not found) -- У John Snow имеется 5 непрочитанных сообщений.
	// У John Snow имеется 5 непрочитанных сообщений.
}

func ExampleNew() {
	e := errors.New("hello world")

	fmt.Println(e)

	// Output:
	// hello world
}
