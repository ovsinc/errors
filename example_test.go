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
	"github.com/ovsinc/errors"
	"github.com/ovsinc/multilog"
	"github.com/ovsinc/multilog/golog"
	"golang.org/x/text/language"
)

var UnknownErrorType = errors.NewObjectFromString("UNKNOWN_TYPE")

func itsOk() error {
	return nil
}

func itsErr(s string) error {
	return errors.New(s)
}

// Добавление ошибок в mutierror с логгированием
// тут изменена функция форматирования вывода -- используется json
func ExampleCombineWithLog() {
	multilog.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))

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
	// 	#1 write: {hello:world} -- connection error
	// 	#2 connection error
}

func ExampleNewWithLog() {
	multilog.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))

	_ = errors.Combine(
		nil,
		itsOk(),
		errors.NewWithLog("one"),
		errors.NewWithLog("two"),
		itsOk(),
	)

	_ = errors.NewWithLog("three")

	// Output:
	// ovsinc/errors one
	// ovsinc/errors two
	// ovsinc/errors three
}

func someTimedCast() (err error) {
	begin := time.Now()
	defer func() {
		err = errors.Cast(err).
			WithOptions(errors.SetContextInfo(
				errors.CtxMap{
					"duration": time.Since(begin).Round(time.Second),
					"call":     errors.HandlerCaller().FilePosition,
				},
			))
	}()

	err = errors.New("some call")

	time.Sleep(1 * time.Second)

	return err
}

func ExampleLog() {
	multilog.DefaultLogger = golog.New(log.New(os.Stdout, "ovsinc/errors ", 0))

	errors.Log(someTimedCast())

	// Output:
	// ovsinc/errors {call:example_test.go:101,duration:1s} -- some call
}

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
		errors.SetOperation("test op"),
	)

	fmt.Println(e1.Error())
	fmt.Println(e2.Error())

	// Output:
	// hello
	// test op: {hello:world} -- hello
}

func ExampleError_TranslateMsg() {
	errEmailsUnreadMsg := localTransContext()
	localizer := localizePrepare()

	e1 := errors.NewWith(
		errors.SetMsg("fallback message"),
		errors.SetID("ErrEmailsUnreadMsg"),
		errors.SetTranslateContext(&errEmailsUnreadMsg),
		errors.SetLocalizer(localizer),
	)

	fmt.Printf("%v\n", e1)
	fmt.Println(e1.Error())
	fmt.Print(e1.TranslateMsg())

	// Output:
	// У John Snow имеется 5 непрочитанных сообщений.
	// У John Snow имеется 5 непрочитанных сообщений.
	// У John Snow имеется 5 непрочитанных сообщений.
}

func ExampleNew() {
	e := errors.New("hello world")

	fmt.Println(e)

	// Output:
	// hello world
}
