package errors_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/ovsinc/errors"
	customlog "gitlab.com/ovsinc/errors/log"
	"gitlab.com/ovsinc/errors/log/golog"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
		errors.New("two", errors.SetSeverity(customlog.SeverityWarn)),
		itsOk(),
	)
	err = errors.Wrap(err, errors.New("three", errors.SetSeverity(customlog.SeverityWarn)))

	fmt.Printf("%v\n", errors.ErrorOrNil(err))
	//Output:
	//[ERROR] -- one
}

// Добавление ошибок в mutierror с логгированием
// тут изменена функция форматирования вывода -- испльзуется json
func ExampleAppendWithLog() {
	logger := log.New(os.Stdout, "ovsinc/errors ", 0)
	customlog.DefaultLogger = golog.New(logger)
	errors.DefaultMultierrFormatFunc = errors.JSONMultierrFuncFormat

	_ = errors.AppendWithLog(
		nil,
		itsOk(),
		itsErr("one"),
		itsErr("two"),
		itsOk(),
	)

	// Output:
	// ovsinc/errors {"count":2,"messages":[{"error_type":"","severity":"ERROR","operations":[],"context":null,"msg":"one"},{"error_type":"","severity":"ERROR","operations":[],"context":null,"msg":"two"}]}
}

func someFuncWithErr() error {
	return errors.New(
		"connection error",
		errors.SetContextInfo(errors.CtxMap{"hello": "world"}),
		errors.AppendOperations("write"),
		errors.SetSeverity(customlog.SeverityUnknown),
		errors.SetErrorType(""),
	)
}

func someFuncWithErr2() error {
	return errors.New(
		"connection error",
		errors.SetSeverity(customlog.SeverityUnknown),
		errors.SetErrorType(""),
	)
}

func ExampleWrap() {
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	err := someFuncWithErr()

	err = errors.Wrap(err, someFuncWithErr2())

	fmt.Printf("%v\n", err)

	//Output:
	//* [write]<hello:world> -- connection error
	//* connection error
}

func ExampleNewWithLog() {
	logger := log.New(os.Stdout, "ovsinc/errors ", 0)
	customlog.DefaultLogger = golog.New(logger)
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
	// ovsinc/errors [ERROR] -- one
	// ovsinc/errors [ERROR] -- two
	// ovsinc/errors [ERROR] -- three
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
	logger := log.New(os.Stdout, "ovsinc/errors ", 0)
	customlog.DefaultLogger = golog.New(logger)
	errors.DefaultMultierrFormatFunc = errors.StringMultierrFormatFunc

	errors.Log(someTimedCast())

	// Output:
	// ovsinc/errors [ERROR]<call:example_test.go:163,duration:1s> -- some call
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

func ExampleError_translateMsg() {
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
	fmt.Print(e1.Error())

	//Output:
	//[not found][ERROR] -- У John Snow имеется 5 непрочитанных сообщений.
	//[not found][ERROR] -- У John Snow имеется 5 непрочитанных сообщений.
}

func ExampleNew() {
	e := errors.New("hello world")

	fmt.Println(e)

	//Output:
	//[ERROR] -- hello world
}
