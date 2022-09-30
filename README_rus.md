# Errors

В процессе работы приложения часто приходится возвращать и обрабатывать ошибки. Стандартный пакет ошибок `errors` достаточно беден в плане возможностей. Пакет `github.com/pkg/errors` более интересен, но также не лишен недостатков.

Этот пакет призван добавить возможностей к обработке ошибок. Для удобства использования используется стратегия, принятая в `github.com/pkg/errors` и в целом в [golang](https://golang.org/). Он совместим со стандартным пакетом `errors`.

## Оглавление

1. [Установка](#Установка)
2. [Миграция](#Миграция)
3. [Тестирование](#Тестирование)
  - [Производительность](#Производительность)
4. [Сценарии использования](#Сценарии-использования)
  - [Замена стандартной errors](#Замена-стандартной-errors)
  - [Дополнительные опции](#Дополнительные-опции)
  - [Логгирование](#Логгирование)
  - [Перевод сообщения ошибки](#Переводсообщения-ошибки)
  - [Цепочка ошибок](#Цепочка-ошибок)
5. [Особенности использования](#Особенности-использования)
  - [Управление логгированием ошибки](#Управление-логгированием-ошибки)
  - [Настройка перевода сообщения ошибки](#Настройка-перевода-сообщения-ошибки)
  - [Сериализация сообщения ошибки](#Сериализация-сообщения-ошибки)
6. [Список задач](#Список-задач)
7. [Лицензия](#Лицензия)

____

## Установка

```text
go get github.com/ovsinc/errors
```

Для простого использования достаточно будет импортировать пакет в своём приложении:

```golang
package main

import (
    "fmt"
    "github.com/ovsinc/errors"
)

func main() {
    fmt.Printf("%v\n", errors.New("hello error"))
}

```

[К оглавлению](#Оглавление)

## Миграция

Поскольку `github.com/ovsinc/errors` совместим с `errors`, то в общем случае миграция достаточно проста.

```golang
package main

import (
    "fmt"
    // "errors"
    "github.com/ovsinc/errors"
)

func main() {
    fmt.Printf("%v\n", errors.New("hello error"))
}

```

[К оглавлению](#Оглавление)

## Тестирование

Склонируйте репозиторий:

```text
git clone https://github.com/ovsinc/errors
cd errors
```

Для запуска юнит-тестов перейдите каталог репозитория в выполните:

```text
make test
```

### Производительность

Для запуска теста производительности перейдите каталог репозитория в выполните:

```text
make bench
```

Для сравнения с аналогичными решениями выполните:

```text
make bench_vendors
```

[К оглавлению](#Оглавление)

## Сценарии использования

### Замена стандартной errors

Тут всё просто: нужно заменить импорт `errors` на `github.com/ovsinc/errors` и пользоваться как было привычно:

```golang
package main

import (
    "fmt"
    // "errors"
    "github.com/ovsinc/errors"
)

func main() {
    fmt.Println(errors.New("hello error"))
}
```

При использовании форматирования возможно более широкое использование использование:

| Глагол и флаг | Описание |
| ------------- | -------- |
| `v` | Строка, сериализованная ошибка в строку. |
| `#v`| Дамп *Error |
| `j` | Строка, сериализованная ошибка в JSON (содержит все поля). |
| `q` | Строка, сериализованная ошибка в строку (с указанием ID). |
| `s` | Строка, сообщение. |
| `+s` | Строка, сообщение с переводом (если возможно). |
| `t` | Строка, тип. |
| `o` | Строка, операция. |
| `с` | Строка, сериализованный контекст. |
| `f` | Строка, место вызова. |

Пример использования:

```golang
package main

import (
    "fmt"
    "github.com/ovsinc/errors"
)

func main() {
    e := errors.NewWith(
        errors.SetMsg("hello error"),
        errors.SetOperation("store to db"),
        errors.SetID("<myid>"),
        errors.SetErrorType("internal"),
        errors.AppendContextInfo("host", "localhost"),
        errors.AppendContextInfo("db", "postgres"),
    )
    
    fmt.Printf("id: %i, op: %o, type: %t ctx: %c, msg: %s", e, e, e, e, e)

    fmt.Printf("%+v", e)

    fmt.Printf("full str: %q", e)

    fmt.Printf("full str: %j", e)
}
```

### Расширенное использование

#### Дополнительные опции

Вызов `NewWith` позволяет создать ошибку с нужными свойствами в стиле функций-параметров.

| Опция | Описание |
| ----- | -------- |
| `Error.SetMsg(string)` | Установить сообщение об ошибке. |
| `Error.SetOperation(string)` | Установит операцию. |
| `Error.SetErrorType(string)` | Установит тип. |
| `Error.SetID(string)` | Установит идентификатор. |
| `Error.SetContextInfo(CtxMap)` | Установит контекст. |
| `Error.AppendContextInfo(string, interface{})` | Добавит контекст к имеющимуся. Если контекст не был создан, создаст. |

#### Логгирование

Существует возможность логгирования ошибки.

Возможные варианты вызова:

- из конструтороа `NewLog` (аналогично конструтору `New`, но с логгированием) или `NewWithLog` (аналогично `NewWith`, но с логгированием);
- вызов метода `Error.Log(...Logger)`;
- хелпер `Log(error, ...Logger)`.

### Перевод сообщения ошибки

Сообщение об ошибке можно перевести. Для корректного выполнения перевода в `*Error` должен быть установлен идентификатор.

Возможные варианты вызова:

- вызов метода `Error.Translate(...Translater) (string, error)`;
- хелпер `Translate(error, ...Translater) (string, error)`;
- форматированныый вывод `Printf` с руной `#s` (используется дефолтный контекст перевода).

В случае ошибки перевода методы вернут оригинальное сообщение.

### Цепочка ошибок

Иногда требуется в одном месте собрать несколько ошибок в цепочке вызовов.
Цепочка вызовов, пример:

```golang
package main

import (
    "net/http"
    "github.com/ovsinc/errors"
)

var (
    ErrModel      = errors.New("some *model* error")
    ErrController = errors.New("some *control* error")
)

type Myhandler struct{}

func (*Myhandler) ModelFunc() error {
    return ErrModel
}

func (h *Myhandler) ControlFunc() error {
    err := h.ModelFunc()
    if err != nil {
        return errors.Wrap(ErrController, err)
    }
    return nil
}

func (h *Myhandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
    code := http.StatusOK
    msg := "Hello world"

    err := h.ControlFunc()
    if err != nil {
        errors.Log(err)
        code = http.StatusInternalServerError
        msg = "Some errors occured\n"
    }

    w.WriteHeader(code)
    w.Write([]byte(msg))
}

func main() {
    h := new(Myhandler)
    http.HandleFunc("/", h.HandleFunc)
    http.ListenAndServe(":8000", nil)
}
```

Или когда нужно в цикле обработать однотипные выводы и в конце вынести общий вердикт.
Общий вердикт, пример:

```golang
package main

import (
    "database/sql"
    "github.com/ovsinc/errors"
)

func main() {
    var (
        srvs = []string{
            "myhsot",
            "localhost",
        }

        err    error
        client *sql.DB
    )

    for _, connStr := range srvs {
        db, e := sql.Open("postgres", connStr)
        if err == nil {
            client = db
            break
        }

        err = errors.Combine(err, e)
    }

    if client == nil {
        errors.Log(err)
    }
}
```

Использование в разных потоках может быть не безопасным!
В горутинах лучше использовать [errgroup](https://pkg.go.dev/golang.org/x/sync@v0.0.0-20220923202941-7f9b1623fab7/errgroup).

## Особенности использования

### Управление логгированием ошибки

Логгирование в пакете реализовано с помощью библиотеки [multilog](https://github.com/ovsinc/multilog).

```golang
type Logger interface {
    Errorf(format string, args ...interface{})
}
```

В пакете установлен логгер по-умолчанию, установленный на использование стандартного для Go логгера `log`.

При необходимости его можно легко переопределить на более подходящее значение из пакета [multilog](https://github.com/ovsinc/multilog).

Возможные варианты вызова:

- метод `*Error.Log(l ...multilog.Logger)`;
- хелпер `errors.Log(error, ...Logger)`;
- методы-конструкторы: `CombineWithLog`,`WrapWithLog`, `NewLog`, `NewWithLog`.

Ниже приведен пример использования `github.com/ovsinc/errors` c логгированием:

```golang
package main

import (
    "github.com/ovsinc/multilog/journald"
    "github.com/ovsinc/multilog/logrus"
    origlogrus "github.com/sirupsen/logrus"
    "github.com/ovsinc/errors"
)

func main() {
    logrusLogger := logrus.New(origlogrus.New())
    errors.DefaultLogger = logrusLogger

    err := errors.NewLog("hello error")

    journalLogger := journald.New()
    errors.Log(err, journalLogger)
}
```

### Настройка перевода сообщения ошибки

Перевод сообщения ошибки реализован с помощью библиотеки `github.com/nicksnyder/go-i18n/v2/i18n`.

Важно чтобы каждый объект `*Error` должен содержать ID. Он используется в `go-i18n` для поиска переводимого сообщения.
В случае перевода сообщения ошибки, заданное сообщение будет использоваться при срабатывании fallback сценария, т.е. при возникновении ошибки при переводе.

Для работы переводов нужно выполнить подготовку:

- инициализировать локализатор go-i18n (должен соответсвовать интерфейсу `errors.Localizer`);
См. подробности в пакете [i18n](https://github.com/nicksnyder/go-i18n).
- при создании сообщения `*Error`, использовать конструктор `NewWith` с заданием ID.

Может оказаться удобным установить локализатор `errors.DefaultLocalizer` для всего приложения. Тогда, конечно, локализатор должен содержать весь набор переводимых сообщений и настроен на использование требуемых языков.

Получение переведенного сообщения:

- хелпер `errors.TranslateMsg(error, Localizer, *TranslateContext)`;
- форматированный вывода с руной `#s`.

В простых случаях, если задан `errors.DefaultLocalizer` для всего приложения, то можно использовать форматированный вывода из составка `fmt` с руной `#s`. При этом необходимо учитывать, что плюральные форммы требуют задание контекста для каждого переводимого сообщения, что в случае с форматированным выводом сделать нельзя.

Пример использование перевода в сообщении ошибки:

```golang
package main

import (
    _ "embed"
    "fmt"

    "github.com/BurntSushi/toml"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "github.com/ovsinc/errors"
    "golang.org/x/text/language"
)

//go:embed testdata/active.ru.toml
var translationRu []byte

func main() {
    bundle := i18n.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
    bundle.MustParseMessageFileBytes(translationRu, "testdata/active.ru.toml")

    localizer := i18n.NewLocalizer(bundle, "ru")
    locCtx := errors.TranslateContext{
            TemplateData: map[string]interface{}{
            "Name":        "John Snow",
            "PluralCount": 5,
        },
        PluralCount: 5,
    }

    err := errors.NewWith(
        errors.SetMsg("fallback message"),
        errors.SetID("ErrEmailsUnreadMsg"),
    )

    fmt.Println(errors.TranslateMsg(err, localizer, &locCtx))
}
```

---------------------------------

### Сериализация сообщения ошибки

В пакете представлены по паре (JSON, string) функций форматирования для единичного сообщения и цепочке сообщений ошибки.

Для цепочки сообщений изменение функции-форматера осуществляется через переменную `DefaultMultierrFormatFunc`. Для неё определено значение по-умолчанию `var DefaultMultierrFormatFunc = StringMultierrFormatFunc`.

Multierror-сообщения форматируются в пакете следующими функциями:

- для вывода в формате JSON - `JSONMultierrFuncFormat(w io.Writer, es []error)`;
- для строкового вывода - `StringMultierrFormatFunc(w io.Writer, es []error)`.

Для сообщений с типом `*Error` используются функции-форматеры типа `type FormatFn func(e *Error) string`. Задать требуемую функцию форматирования можно с помощью функции-параметра `SetFormatFn` в конструкторе или изменить это значение с помощью метода `WithOptions`. Можно задать функцию-форматирования по-умолчанию через переменную `DefaultFormatFn`.

В пакете представлены следующие функции-форматеры:

- для вывода в формате JSON - `JSONFormat(buf io.Writer, e *Error)`;
- для строкового вывода - `StringFormat(buf io.Writer, e *Error)`.

Внимание! При использовании форматирования цепочки сообщения `JSONMultierrFuncFormat` функция форматирование `*Error` по-умолчанию переключается на `JSONFormat`.

Все функций форматирования используют `github.com/valyala/bytebufferpool`, что хорошо сказывается на общей производительности и уменьшает потребление памяти.

[К оглавлению](#Оглавление)


### Хелперы

Is(err, target error) bool
As(err error, target interface{}) bool
Unwrap(err error) error

GetID(err error) (id string)
Cast(err error) (*Error, bool)

CastMultierr(err error) (Multierror, bool)
UnwrapByID(err error, id string) *Error
UnwrapByErr(err error, target error) *Error
Contains(err error, target error) bool
ContainsByID(err error, id string) bool

## Список задач

- [ ] Повысить покрытие тестами;
- [ ] Более подробные комментарии для описания методов и функций;
- [ ] Перевод README на en.
- [ ] Проработать сценарии использования в handler (HTTP, GRPC,..)

[К оглавлению](#Оглавление)

## Лицензия

Код пакета распространяется под лицензией [Apache 2.0](http://directory.fsf.org/wiki/License:Apache2.0).
