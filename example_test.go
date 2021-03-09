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
		itsErr("two"),
		itsOk(),
	)
	err = errors.Wrap(err, itsErr("three"))

	fmt.Printf("%v\n", errors.ErrorOrNil(err))
	//Output:
	//[UNKNOWN_TYPE][ERROR] -- one
}

// Добавление ошибок в mutierror с логгированием
// тут изменена функция форматирования вывода -- испльзуется json
func ExampleAppendWithLog() {
	logger := log.New(os.Stdout, "ovsinc/errors ", 0)
	customlog.DefaultLogger = golog.New(logger)
	errors.DefaultMultierrFormatFunc = errors.JsonMultierrFuncFormat

	_ = errors.AppendWithLog(
		nil,
		itsOk(),
		itsErr("one"),
		itsErr("two"),
		itsOk(),
	)

	//Output:
	//ovsinc/errors {"count":2,"messages":[{"error_type":"UNKNOWN_TYPE","severity":"ERROR","operations":[],"context":null,"msg":"one"},{"error_type":"UNKNOWN_TYPE","severity":"ERROR","operations":[],"context":null,"msg":"two"}]}
}

func someFuncWithErr() error {
	return errors.New(
		"connection error",
		errors.SetContextInfo(errors.CtxMap{"hello": "world"}),
		errors.AppendOperations("write"),
		errors.SetSeverity(customlog.SeverityUnknown),
		errors.SetErrorType(errors.NewErrorType("")),
	)
}

func someFuncWithErr2() error {
	return errors.New(
		"connection error",
		errors.SetSeverity(customlog.SeverityUnknown),
		errors.SetErrorType(errors.NewErrorType("")),
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

	//Output:
	//ovsinc/errors [UNKNOWN_TYPE][ERROR] -- one
	//ovsinc/errors [UNKNOWN_TYPE][ERROR] -- two
	//ovsinc/errors [UNKNOWN_TYPE][ERROR] -- three
}

func someErrFunc() error {
	return errors.New("connection error", errors.SetErrorType(errors.NewErrorType("NOT_FOUND")))
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

	//Output:
	//Got error with type NOT_FOUND
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

	//Output:
	//ovsinc/errors [UNKNOWN_TYPE][ERROR]<call:example_test.go:163,duration:1s> -- some call
}

func localizePrepare() *i18n.Localizer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./testdata/active.ru.toml")

	return i18n.NewLocalizer(bundle, "ru")
}

func localMsg() errors.TranslateMessage {
	var (
		unreadEmailCount = 5
		name             = "John Snow"
	)

	return errors.TranslateMessage{
		ID: "ErrEmailsUnreadMsg",
		TemplateData: map[string]interface{}{
			"Name":        name,
			"PluralCount": unreadEmailCount,
		},
		PluralCount: unreadEmailCount,
	}
}

func ExampleError_TranslateMsg() {
	errEmailsUnreadMsg := localMsg()
	localizer := localizePrepare()

	e1 := errors.New(
		"fallback message",
		errors.SetErrorType(errors.NewErrorType("not found")),
		errors.AddTranslatedMessage("ru", &errors.TranslateMessage{
			TemplateData: errEmailsUnreadMsg.TemplateData,
			PluralCount:  errEmailsUnreadMsg.PluralCount,
			Localizer:    localizer,
			ID:           errEmailsUnreadMsg.ID,
		}),
		errors.SetLang("ru"),
	)

	fmt.Printf("%v\n", e1)
	fmt.Printf("%v\n", e1.WithOptions(errors.SetLang("en")))

	//Output:
	//[not found][ERROR] -- У John Snow имеется 5 непрочитанных сообщений.
	//[not found][ERROR] -- fallback message
}

func ExampleTranslate() {
	var e1 error

	errEmailsUnreadMsg := localMsg()
	localizer := localizePrepare()

	e1 = errors.New(
		"fallback msg",
		errors.SetErrorType(errors.NewErrorType("not found")),
		errors.AddTranslatedMessage("ru", &errors.TranslateMessage{
			TemplateData: errEmailsUnreadMsg.TemplateData,
			PluralCount:  errEmailsUnreadMsg.PluralCount,
			Localizer:    localizer,
			ID:           errEmailsUnreadMsg.ID,
		}),
		errors.SetLang("ru"),
	)

	fmt.Println(errors.Translate(e1, "ru"))
	fmt.Printf("%v\n", errors.Translate(fmt.Errorf("Hello there"), "ru"))

	//Output:
	//[not found][ERROR] -- У John Snow имеется 5 непрочитанных сообщений.
	//Hello there
}

func ExampleNew() {
	e := errors.New("hello world")

	fmt.Println(e)

	//Output:
	//[UNKNOWN_TYPE][ERROR] -- hello world
}
