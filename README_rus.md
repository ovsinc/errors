# Errors

В процессе работы приложения часто приходится кидать и обрабатывать ошибки. Стандартный пакет ошибок `errors` достаточно беден в плане возможностей. Пакет `github.com/pkg/errors` более интересен в плане возможностей, но тоже не лишен недостатоков.

Этот пакет призван добавить возможностей к обработке ошибок. Для удобства использования используется стратегия, принятая в `github.com/pkg/errors` и в целом в [golang](https://golang.org/).

## Оглавление

1. [Установка](#Установка)
2. [Тестирование](#Тестирование)
3. [Фичи](#Фичи)
4. [Использование](#Использование)
5. [Список задач](#Список-задач)

____

## Установка

```text
go get gitlab.com/ovsinc/errors
```

Для простого использования достаточно будет импортировать пакет в своём приложении:

```golang
package main

import (
    "fmt"
    "gitlab.com/ovsinc/errors"
)

func main() {
    fmt.Printf("%v\n", errors.New("hello error"))
}

```

[К оглавлению](#Оглавление)

## Тестирование

Склонируйте репозиторий:

```text
git clone https://gitlab.com/ovsinc/errors
cd errors
```

Для запуска юнит-тестов перейдите каталог репозитория в выполните:

```text
go test -short -mod=readonly ./...
```

или

```text
make test
```

Для запуска теста производительностив перейдите каталог репозитория в выполните:

```text
go test -benchmem -run=^# -bench=. .
```

или

```text
make bench
```

[К оглавлению](#Оглавление)

## Фичи

- Стандарный интерфейс ошибки (`error`);
- Дополнительные поля описания ошибки: тип (ErrorType), контекст (map[string]interface{}), severity (Enum: Error, Warn), operations([]Operation);
- Дополнительные методы: [перевода](###Перевод-ошибки), [логгирования](###Логгирование-ошибки);
- Хелперы для поддержки наиболее распространенных логгерров: [journald](https://github.com/coreos/go-systemd/), [log](https://golang.org/pkg/log/), [log15](https://github.com/inconshreveable/log15/), [logrus](https://github.com/sirupsen/logrus), [syslog](https://golang.org/pkg/log/syslog/), [zap](https://go.uber.org/zap);
- Логгирование в цепочке.

[К оглавлению](#Оглавление)

## Использование

Пакет имеет обратную совместимость по методам со стандартной библиотекой `errors` и `github.com/pkg/errors`. Поэтому может быть использован в стандартных сценариях, а также с дополнитеьными возможностями.

Внимание! Тип `*Error` не является потокобезопасным. Смотрите пример с использованием мьютекса от [go-multierror](https://github.com/hashicorp/go-multierror/blob/master/group.go) для обработки ошибок в горутинах.

Дополнительные к стандартным кейсы:

- оборачивание цепочки ошибок;
- логгирование ошибки в момент ее формирования;
- формирование ошибок в процессе выполнения цепочки методов и проверка в верхнем методе (с возвожным логгированием);
- выдача ошибок (ошибки) клиентскому приложению с переводом сообщения на язык клиента.

Для работы с цепочкой ошибок используется пакет `github.com/hashicorp/go-multierror`.

Для переводов сообщения используется библиотека `github.com/nicksnyder/go-i18n/v2/i18n`. Ознакомится с особенностями работы i18n можно [тут](https://github.com/nicksnyder/go-i18n).

Можно также ознакомится с [примерами](https://gitlab.com/ovsinc/errors/-/blob/master/example_test.go) использования `gitlab.com/ovsinc/errors`.

### Методы *Error

| Метод | Описание |
| ----- | -------- |
| `func New(msg string, ops ...Options) *Error` | Конструктор `*Error`. Обязательно нужно указать сообщение ошибки. Дополнительные свойства устанваливаются с момощью [функций-параметров](###Функции-параметры). |
| `func NewWithLog(msg string, ops ...Options) *Error` | Конструктор `*Error`, как и `New`. Перед возвратом `*Error` производит логгирование на дефолтном логгере. |
| `func (e *Error) Error() string` | Метод, возвращающий строку сообщения ошибки. |
| `func (e *Error) WithOptions(ops ...Options) *Error` | Изменение свойств `*Error` производится с помощью [функций-пааметров](###Функции-параметры). |
| `func (e *Error) Severity() log.Severity` | Геттер получения важности ошибки. |
| `func (e *Error) Msg() string` | Геттер получения сообщения ошибки. |
| `func (e *Error) Sdump() string` | Геттер получения текстового дампа `*Error`. Может использоваться для отладки. |
| `func (e *Error) ErrorOrNil() error` | Геттер получения ошибки или `nil`. `*Error` с типом `log.SeverityWarn` не является ошибкой; метод `ErrorOrNil` с таким типом вернет `nil`. |
| `func (e *Error) Operations() []Operation` | Геттер получения списка операций. |
| `func (e *Error) ErrorType() ErrorType` | Геттер получения типа ошибки. |
| `func (e *Error) Format(s fmt.State, verb rune)` | Функция форматирования для обработки строк с возможностью задания формата, например: `fmt.Printf`, `fmt.Sprintf`, `fmt.Fprintf`,.. |

### Функции-параметры

Параметризация `*Error` производится с помощью функций-параметров типа `type Options func(e *Error)`.

| Метод | Описание |
| ----- | -------- |
| `func SetFormatFn(fn FormatFn) Options` | Устанавливает кастомную функцию форматирования. Если значение `nil`, будет использоваться дефолтная функция форматирования. |
| `func SetMsg(msg string) Options` | Установить сообщение. |
| `func SetSeverity(severity log.Severity) Options` | Установить уровень важности сообщения. Доступные значения: `log.SeverityWarn`, `log.SeverityError`. |
| `func SetLang(lang string) Options` | Установить язык перевода. |
| `func AddTranslatedMessage(lang string, message *TranslateMessage) Options` | Установить `*TranslateMessage` для указанного языка. |
| `func SetErrorType(etype ErrorType) Options` | Установить тип ошибки. Тип `ErrorType` является производным типа `string`, так что создание собственных типов ошибки легко выполнить с помощью конструктора `func NewErrorType(s string) ErrorType`. |
| `func SetOperations(ops ...Operation) Options` | Установить список выполненных операций. Тип `Operation` является производным типа `string`. Задать операцию можно с помощью конструктора `func NewOperation(s string) Operation`. |
| `func AppendOperations(ops ...Operation) Options` | Добавить операции к уже имеющемуся списку. Если список операций не существует, он будет создан. |
| `func SetContextInfo(ctxinf CtxMap) Options` | Задать контекст ошибки. |
| `func AppendContextInfo(key string, value interface{}) Options` | Добавить значения к уже имеющемуся контексту ошибки. Если контекст ошибки не существует, он будет создан. |

### Хелперы

Все хелперы работают с типом `error`.

| Хелпер | Описание |
| ------ | -------- |
| `func Is(err, target error) bool` | Обёртка над методом стандартной библиотеки `errors.Is`. |
| `func As(err error, target interface{}) bool` | Обёртка над методом стандартной библиотеки `errors.As`. |
| `func GetErrorType(err error) ErrorType` | Получить тип ошибки. Для НЕ `*Error` всегда будет `UnknownErrorType`. |
| `func ErrorOrNil(err error) error` | Возможна обработка цепочки или одиночной ошибки. Если хотя бы одна ошибка в цепочке является ошибкой, то она будет возвращена в качестве результата. Важно: `*Error` c Severity `Warn` не является ошибкой. |
| `func Cast(err error) *Error` | Преобразование типа `error` в `*Error`. |
| `func Append(err error, errs ...error) *multierror.Error` | Добавить в цепочку ошибок еще ошибки. Допускается использование `nil` в обоих аргументах. |
| `func Wrap(olderr error, err error) *multierror.Error` | Обернуть ошибку `olderr` ошибкой `err`, получив цепочку. Допускается использование `nil` в обоих аргументах. |
| `func Unwrap(err error) error` | Развернуть цепочку ошибок, получив первую. |

### Логгирование ошибки

Логгирование в пакете представлено интерфейсом:

```golang
type Logger interface {
    Warn(err error)
    Error(err error)
}
```

В пакете пристутвует логгер по-умолчанию `gitlab.com/ovsinc/errors/log.DefaultLogger`.
Он установлен в значение

```golang
var DefaultLogger = golog.New(pkglog.New(os.Stderr, "ovsinc/errors ", pkglog.LstdFlags))
```

При необходимости его можно легко переопределить на более подходящее значение.

Для логгирования в `*Error` имеется метод `Log(l ...logcommon.Logger)`. Однако, приводить `error` к `*Error` каждый раз не требуется.

Для логгирования в пакете есть несколько хелперов.

| Хелпер | Описание |
| ------ | -------- |
| NewWithLog(msg string, ops ...Options) *Error | Функция произведет логгирование ошибки дефолтным логгером. |
|Log(err error, l ...logcommon.Logger) | Функция произведет логгирование ошибки дефолтным логгером или логгером указанным в l (будет использоваться только первое значение). |
| AppendWithLog(err error, errs ...error) *multierror.Error | Обертка над  `github.com/hashicorp/go-multierror.Append`, которая выполнит логгирование дефолтным логгером и вернет цепочку `*multierror.Error`. |
| WrapWithLog(olderr error, err error) *multierror.Error | Хелпер обернет `olderr` ошибкой `err`, выполнит логгирование дефолтным логгером и вернет цепочку `*multierror.Error`. |

Для удобства поддерживаются несколько оберток над наиболее популярными логгерами.

Ниже приведен пример использования `gitlab.com/ovsinc/errors` c логгированием:

```golang
package main

import (
    "time"

    origlogrus "github.com/sirupsen/logrus"
    "gitlab.com/ovsinc/errors"
    "gitlab.com/ovsinc/errors/log"
    "gitlab.com/ovsinc/errors/log/chain"
    "gitlab.com/ovsinc/errors/log/journald"
    "gitlab.com/ovsinc/errors/log/logrus"
)

func main() {
    now := time.Now()

    logrusLogger := logrus.New(origlogrus.New())

    log.DefaultLogger = logrusLogger

    err := errors.NewWithLog(
        "hello error",
        errors.SetSeverity(log.SeverityWarn),
        errors.SetContextInfo(
            errors.CtxMap{
                "time": now,
            },
        ),
    )

    err = err.WithOptions(
        errors.SetSeverity(log.SeverityError),
        errors.AppendContextInfo("duration", time.Since(now)),
    )

    journalLogger := journald.New()

    chainLogger := chain.New(logrusLogger, journalLogger)

    err.Log(chainLogger)
}
```

### Перевод ошибки

Для переводов сообщения ошибки используется библиотека `github.com/nicksnyder/go-i18n/v2/i18n`.

В структуре `*Error` за перевод отвечают несколько свойств.

| Свойство | Тип |Назначение | Значение по-умолчанию |
| -------- | --- | --------- | --------------------- |
| translateMap translateMap | `type translateMap map[string]*TranslateMessage`, `type TranslateMessage struct { ID string DefaultMessage *i18n.Message TemplateData map[string]interface{} PluralCount interface{} Localizer *i18n.Localizer }` | Сообщение ошибки для установленного языка. | `nil` |
| lang string | `string` | Выбранный язык для перевода сообщения ошибки. | `en` |

Для каждого языка перводов необходимо установить собственное сообщение `translateMap` с обязательным указанием локализатора `Localizer *i18n.Localizer` и `ID string` сообщения. Можно установить дефолтное локализованное сообщение `DefaultMessage *i18n.Message`. Также можно установить шаблон подстановки `TemplateData map[string]interface{}` и множественность `PluralCount interface{}`.

Тогда при вызове метода `gitlab.com/ovsinc/errors.*Errors.Error()` будет выдана строка с переведенным сообщением.

В случае возникновения ошибки при переводе сообщения `*Error` будет выдана строка с оригинальным сообщением, без перевода.

Пример ошибки с переводом:

```golang
package main

import (
    "fmt"
    "gitlab.com/ovsinc/errors"
    "github.com/BurntSushi/toml"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
)

func main() {
    bundle := i18n.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
    bundle.MustLoadMessageFile("./testdata/active.ru.toml")

    err :=  errors.New(
        "fallback message",
        errors.SetErrorType(errors.NewErrorType("not found")),
        errors.AddTranslatedMessage("ru", &errors.TranslateMessage{
            TemplateData: map[string]interface{}{
                "Name":        "John Snow",
                "PluralCount": 5,
            },
            PluralCount: 5,
            Localizer: i18n.NewLocalizer(bundle, "ru"),
            ID: "ErrEmailsUnreadMsg",
        }),
        errors.SetLang("ru"),
    )

    fmt.Printf("%v\n", err)
}
```

### Функции форматирования сообщения ошибки

В пакете представлены по паре (JSON, string) функций форматирования для единичного сообщения и multierror-сообщения.

Для multierror-сообщений изменние функции-форматера осуществляется через переменную `DefaultMultierrFormatFunc`. Для неё определено значение по-умолчанию `var DefaultMultierrFormatFunc = StringMultierrFormatFunc`.

Multierror-сообщения форматируются в пакете следующими функциями:

- для вывода в формате JSON - `func JsonMultierrFuncFormat(es []error) string`;
- для строкового вывода - `func StringMultierrFormatFunc(es []error) string`.

Для сообщений с типом `*Error` используются функции-форматеры типа `type FormatFn func(e *Error) string`. Задать требуемую фукнцию форматирования можно задать с помощью функции-параметра `SetFormatFn` в констуркторе или изменить это значение с помощью метода `WithOptions`.

В пакете представлены следующии функции-форматеры:

- для вывода в формате JSON - `func JSONFormat(e *Error) string`;
- для строкового вывода - `func StringFormat(e *Error) string`.

Внимание! При использовании форматирования multierror-сообщения `JsonMultierrFuncFormat` функция форматирование `*Error` по-умолчанию переключается на `JSONFormat`.

Все функций форматирования используют `github.com/valyala/bytebufferpool`, что хорошо сказывается на общей производительности и уменьшает потребление памяти.

[К оглавлению](#Оглавление)

## Список задач

- [ ] Повысить покрытие тестами;
- [ ] Более подробные коментарии для описаниея методов и фукций;
- [ ] Перевод типа ошибки, операций, уровня опасности;
- [ ] Перевод README на en;
- [ ] Выпуск на godoc.

[К оглавлению](#Оглавление)
