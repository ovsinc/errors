# Errors

В процессе работы приложения часто приходится возвращать и обрабатывать ошибки. Стандартный пакет ошибок `errors` достаточно беден в плане возможностей. Пакет `github.com/pkg/errors` более интересен, но также не лишен недостатков.

Этот пакет призван добавить возможностей к обработке ошибок. Для удобства использования используется стратегия, принятая в `github.com/pkg/errors` и в целом в [golang](https://golang.org/). Он совместим со стандартным пакетом `errors`.

## Оглавление

1. [Установка](#установка)
2. [Миграция](#миграция)
3. [Тестирование](#тестирование)
4. [Сценарии использования](#сценарии-использования)
5. [Особенности использования](#особенности-использования)
6. [Список задач](#список-задач)
7. [Лицензия](#лицензия)

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

[К оглавлению](#оглавление)

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

[К оглавлению](#оглавление)

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

Сравнение основных возможностей для похожих решений, одна ошибка:

```text
go test -benchmem -run=^$ -bench "^(BenchmarkVendorStandartError|BenchmarkVendorStandartConstructor|BenchmarkVendorXerrors|BenchmarkVendorXerrorsConstructor|BenchmarkVendorMyNewFull|BenchmarkVendorMyNewFullConstructor|BenchmarkVendorMyNewSimple|BenchmarkVendorMyNewSimpleConstructor)$"
goos: linux
goarch: amd64
pkg: github.com/ovsinc/errors
cpu: Intel(R) Core(TM) i7-10850H CPU @ 2.70GHz
BenchmarkVendorStandartError-12                 713995772                1.698 ns/op           0 B/op          0 allocs/op
BenchmarkVendorStandartConstructor-12           1000000000               0.3875 ns/op          0 B/op          0 allocs/op
BenchmarkVendorXerrors-12                       718511691                1.648 ns/op           0 B/op          0 allocs/op
BenchmarkVendorXerrorsConstructor-12             2052255               588.9 ns/op             0 B/op          0 allocs/op
BenchmarkVendorMyNewSimple-12                   15405230                77.94 ns/op            0 B/op          0 allocs/op
BenchmarkVendorMyNewSimpleConstructor-12        18630474                80.98 ns/op           80 B/op          1 allocs/op
BenchmarkVendorMyNewFull-12                      6755743               158.3 ns/op             0 B/op          0 allocs/op
BenchmarkVendorMyNewFullConstructor-12           2980184               429.6 ns/op           416 B/op          3 allocs/op
```

Сравнение основных возможностей для похожих решений, две шибки в цепочке:

```text
go test -benchmem -run=^$ -bench "^(BenchmarkVendorMyMulti2StdErr|BenchmarkVendorMyMulti2MySimple|BenchmarkVendorHashiMulti2StdErr|BenchmarkVendorUberMulti2StdErr)$"
goos: linux
goarch: amd64
pkg: github.com/ovsinc/errors
cpu: Intel(R) Core(TM) i7-10850H CPU @ 2.70GHz
BenchmarkVendorMyMulti2StdErr-12         2840337               462.3 ns/op            64 B/op          5 allocs/op
BenchmarkVendorMyMulti2MySimple-12       2606301               476.8 ns/op            64 B/op          5 allocs/op
BenchmarkVendorHashiMulti2StdErr-12      1405405               939.4 ns/op           136 B/op          6 allocs/op
BenchmarkVendorUberMulti2StdErr-12       9235027               128.3 ns/op            16 B/op          1 allocs/op
```

[К оглавлению](#оглавление)

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
        errors.SetErrorType(errors.NotFound),
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

#### Дополнительные свойства

Вызов `NewWith` позволяет создать ошибку с нужными свойствами в стиле функций-параметров.

| Опция | Описание |
| ----- | -------- |
| `SetMsg(string)` | Установить сообщение об ошибке. |
| `SetOperation(string)` | Установит операцию. |
| `SetErrorType(errType)` | Установит тип. |
| `SetID(string)` | Установит идентификатор. |
| `SetContextInfo(CtxKV)` | Установит контекст. |
| `AppendContextInfo(string, string)` | Добавит контекст к имеющимуся. Если контекст не был создан, создаст. |
| `SetErrorType(et IErrType) Options` | Установить тип ошибки. Если не указано, то устанавливается тип `Unknown`. |

#### Логгирование

Существует возможность логгирования ошибки.

Возможные варианты вызова:

- из конструтороа `NewLog` (аналогично конструтору `New`, но с логгированием) или `NewWithLog` (аналогично `NewWith`, но с логгированием);
- вызов метода `*Error.Log(...Logger)`;
- хелпер `Log(error, ...Logger)`.

#### Типизированная ошибка

Ошибки могут иметь тип. Например, `Validation`, `InputBody`.
Типизация может понадобится для удобной обработки обработки в хендлерах.

Тип реализует интерфейс:

```golang
type IErrType interface {
    HTTPStatusCode() int
    GRPCStatusCode() codes.Code
    Number() int
    String() string
}
```

В пакете присутсвуют конструкторы для типизированных ошибок. Например, для создания ошибки валидации можно использовать разные варианты, например:

```golang
errors.NewWith(
    errors.SetMsg("some validation error"),
    errors.SetErrorType(errors.Validation),
)
// the same
errors.ValidationErr("some validation error")
```

Методы конструкторы именуются по следующей схеме:

1. `{{type}}Err(string) *Error` - создание ошибки типа `type` и сообщения;

2. `{{type}}ErrWith(...Options) *Error` - создание ошибки `type` и с опциями;

Для второго пособа важно отметить, что тип устанавливается в соответсвии с именем функции конструктора, даже если использовалась опция `SetErrorType` с другим типом.

Типы:

```golang
const (
    _ errType = iota

    // Неизвестный тип ошибки. Дефолтное значение.
    Unknown

    // Internal внутренняя системная ошибка. Например, отказ базы данных.
    Internal

    // Validation ошибка валидации. Например, не корректный email-адрес.
    Validation

    // InputBody ошибка обработки входных данных. Например, ошибка сериализации JSON.
    InputBody

    // Duplicate дубликат данных, нарушения уникальности.
    Duplicate

    // Unauthenticated для выполнения запроса требуется аутентфиикация.
    Unauthenticated

    // Unauthorized доступ запрещен, запрос не авторизован.
    Unauthorized

    // Empty запрос или не ответ не должен быть пустым.
    Empty

    // NotFound запрашиваемые данные не найдены. Например, пользователь с заданным ID не найден.
    NotFound

    // MaximumAttempts превышение числе разрешенных попуток выполнения одного и того же действия.
    MaximumAttempts

    // SubscriptionExpired срок действия "оплаченой" подписки истек.
    SubscriptionExpired

    // DownstreamDependencyTimedout время ожидания выполнения запрос к нижестоящему сервису истек.
    DownstreamDependencyTimedout

    // Unavailable сервис не доступен.
    Unavailable
)
```

### Перевод сообщения ошибки

Сообщение об ошибке можно перевести.
Для корректного выполнения перевода в `*Error` должен быть установлен идентификатор,
который должен быть идентичным с объектом сообщения (`i18n.Message`).

Возможные варианты вызова:

- вызов метода `*Error.Translate(...Translater) (string, error)`;
- хелпер `Translate(error, ...Translater) (string, error)`;
- форматированныый вывод `Printf` с руной `#s` (используется дефолтный контекст перевода).

В случае ошибки перевода все эти методы вернут оригинальное сообщение.

Подробнее описано [тут](#настройка-перевода-сообщения-ошибки).

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
    w.Write(s2b(msg))
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

### Финальная обработка ошибок

В golang, к сожалению, нет удобного механизма обработки исключений, как в python. Принято "поднимать" ошибку на более высокий уровень вызова по цепочке.

Например:

```golang

import "errors"

var Err1 = errors.New("some error")

func fn1() error {
    return Err1
}

func fn2() error {
    err := fn1()
    if err != nil {
        return err
    }
    return nil
}
```

Тем не мене подход в python try-except выглядит интересным.

```python
class Err1(Exception):
    pass

class Err2(Exception):
    pass

class Err3(Exception):
    pass

def fn():
    raise Err1()

def main:
    try:
        fn()
    # except named exception
    except Err1 as err:
        print("Error: {0}".format(err))
    # except named exception
    except Err2 as err:
        print("Error: {0}".format(err))
    # default except
    except:
        raise
```

В пакете `errors` для подобной реализации есть инструменты.
Как пример выше на python можно похоже реализовать на golang:

```golang
package main

import (
    "fmt"
    "github.com/ovsinc/errors"
)

const (
    Err1ID = "one"
    Err2ID = "two"
    Err3ID = "tree"
)

var (
    Err1 = errors.NewWith(
        errors.SetMsg("error one"),
        errors.SetID(Err1ID),
    )
    Err2 = errors.NewWith(
        errors.SetMsg("error two"),
        errors.SetID(Err2ID),
    )
    Err3 = errors.NewWith(
        errors.SetMsg("error tree"),
        errors.SetID(Err3ID),
    )
)

func fn() error {
    return Err1
}

func main() {
    var e error

    err := fn() // try
    switch {
    // except named exception
    case errors.ContainsByID(err, Err1ID):
        e = errors.FindByID(err, Err1ID)

    // except named exception
    case errors.ContainsByID(err, Err2ID):
        e = errors.FindByID(err, Err2ID)

    // default except
    default:
        e = errors.FindByID(err, Err3ID)
    }

    fmt.Printf("%v\n", e)
}

Подробнее можно ознакомится в примере [real_world_example](https://github.com/ovsinc/errors/tree/new_approach/_examples/real_world_example).

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
- хелпер `Log(error, ...Logger)`;
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

- инициализировать локализатор go-i18n (должен соответсвовать интерфейсу `Localizer`);
См. подробности в пакете [i18n](https://github.com/nicksnyder/go-i18n).
- при создании сообщения `*Error`, использовать конструктор `NewWith` с заданием ID.

Может оказаться удобным установить локализатор `DefaultLocalizer` для всего приложения. Тогда, конечно, локализатор должен содержать весь набор переводимых сообщений и настроен на использование требуемых языков.

Получение переведенного сообщения:

- хелпер `TranslateMsg(error, Localizer, *TranslateContext)`;
- форматированный вывода с руной `#s`.

В простых случаях, если задан `DefaultLocalizer` для всего приложения, то можно использовать форматированный вывода из составка `fmt` с руной `#s`. При этом необходимо учитывать, что плюральные форммы требуют задание контекста для каждого переводимого сообщения, что в случае с форматированным выводом сделать нельзя.

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

### Сериализация сообщения ошибки

В пакете предосталены два врианта маршалинга: string, json.

По деволту используется маршалинг в строку (string).

Изменить дефолтный маршалинг можно установив переменную `DefaultMarshaller` в нужное значение.

Кастомный маршалинг должен реализовывать следующий интерфейс:

```golang
type Marshaller interface {
    Marshal(interface{}) ([]byte, error)
    MarshalTo(interface{}, io.Writer) error
}
```

В `*Error` имеется метод маршалинга (`Marshal(fn ...Marshaller) ([]byte, error)`), если не указано, используется дефолтный.

[К оглавлению](#оглавление)

### Хелперы

#### Стандартные хелперы для работы error

Функция `Is(err, target error) bool` сравнивает ошибку со значением.

Функция `As(err error, target interface{}) bool` проверяет, относится ли ошибка к конкретному типу.

Функция `Unwrap(err error) error` вернет, упакованную ошибку.
Например, если ошибка была обернута с помоющь `Combine`, `Wrap`.

#### Для организации errors flow

Функция `ContainsByID(err error, id string) bool` проверит содержится ли ошибка с указанным ID в цепочке.

Функция `FindByID(err error, id string) error` вернет ошибку с указанным ID, если он есть в цепочке.
Если нет вернет `nil`. Рекомендуется применять в комбинации с `ContainsByID`.

Функция `GetID(err error) (id string)` вернет ID ошибки. Если ID нет, то вернет пустую строку.

[К оглавлению](#оглавление)

## Список задач

- [ ] Повысить покрытие тестами;
- [ ] Более подробные комментарии для описания методов и функций;
- [ ] Перевод README на en.
- [+] Проработать сценарии использования в handler (HTTP, GRPC,..)

[К оглавлению](#оглавление)

## Лицензия

Код пакета распространяется под лицензией [Apache 2.0](http://directory.fsf.org/wiki/License:Apache2.0).
