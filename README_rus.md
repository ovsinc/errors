# Errors

В процессе работы приложения часто приходится кидать и обрабатывать ошибки. Стандартный пакет ошибок `errors` достаточно беден в плане возможностей. Пакет `github.com/pkg/errors` более интересен в плане возможностей, но тоже не лишен недостатков.

Этот пакет призван добавить возможностей к обработке ошибок. Для удобства использования используется стратегия, принятая в `github.com/pkg/errors` и в целом в [golang](https://golang.org/).

## Оглавление

1. [Установка](#Установка)
2. [Тестирование](#Тестирование)
3. [Фичи](#Фичи)
4. [Использование](#Использование)
   - [Методы *Error](#Методы-*Error)
   - [Функции-параметры](#Функции-параметры)
   - [Основные хелперы](#Основные-хелперы)
   - [Логгирование ошибки](#Логгирование-ошибки)
   - [Перевод ошибки](#Перевод-ошибки)
   - [Функции форматирования сообщения ошибки](#Функции-форматирования-сообщения-ошибки)
5. [Список задач](#Список-задач)
6. [Лицензия](#Лицензия)

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

Для запуска теста производительности перейдите каталог репозитория в выполните:

```text
go test -benchmem -run=^# -bench=. .
```

или

```text
make bench
```

[К оглавлению](#Оглавление)

## Фичи

- Стандартный интерфейс ошибки (`error`);
- Дополнительные поля описания ошибки: тип (`string`), контекст (`map[string]interface{}`), severity (Enum: `log.SeveritError`, `log.SeverityWarn`), операции (`[]string`);
- Дополнительный метод: [логгирование](#Логгирование-ошибки);
- Сообщение ошибки может быть [локализовано](#Перевод-ошибки);
- Хелперы для поддержки наиболее распространенных логгерров: [journald](https://github.com/coreos/go-systemd/), [log](https://golang.org/pkg/log/), [log15](https://github.com/inconshreveable/log15/), [logrus](https://github.com/sirupsen/logrus), [syslog](https://golang.org/pkg/log/syslog/), [zap](https://go.uber.org/zap);
- Логгирование в цепочке.

[К оглавлению](#Оглавление)

## Использование

Тип собственной ошибки (`*Error`) представлен интерфейсом `Errorer`:

```golang
type Errorer interface {
    WithOptions(ops ...Options) Errorer
    ID() string
    Severity() log.Severity
    Msg() string
    Error() string
    Sdump() string
    ErrorOrNil() error

    ErrorType() string

    Operations() []string
    Format(s fmt.State, verb rune)

    TranslateContext() *TranslateContext
    Localizer() *i18n.Localizer

    Log(l ...logcommon.Logger)
}
```

Пакет имеет обратную совместимость по методам со стандартной библиотекой `errors` и `github.com/pkg/errors`. Поэтому может быть использован в стандартных сценариях, а также с дополнительными возможностями.

Внимание! Тип `*Error`, реализующий интерфейс `Errorer` не является потокобезопасным.

Дополнительные к стандартным кейсы:

- оборачивание цепочки ошибок;
- логгирование ошибки в момент ее формирования;
- формирование ошибок в процессе выполнения цепочки методов и проверка в верхнем методе (с возможным логгированием);
- выдача ошибок (ошибки) клиентскому приложению с переводом сообщения на язык клиента (при установки локализатора).

Для переводов сообщения используется библиотека `github.com/nicksnyder/go-i18n/v2/i18n`. Ознакомится с особенностями работы i18n можно [тут](https://github.com/nicksnyder/go-i18n).

Можно также ознакомится с [примерами](https://gitlab.com/ovsinc/errors/-/blob/master/example_test.go) использования `gitlab.com/ovsinc/errors`.

### Методы *Error

| Метод | Описание |
| ----- | -------- |
| `func New(msg string, ops ...Options) Errorer` | Конструктор `*Error`. Обязательно нужно указать сообщение ошибки. Для ошибки будет установлен `severity = log.SeverityError`, `errorType = UnknownErrorType`. Свойства `*Error` можно установить или переопределить с помощью [функций-параметров](#Функции-параметры). |
| `func NewWithLog(msg string, ops ...Options) Errorer` | Конструктор `*Error`, как и `New`. Перед возвратом `*Error` производит логгирование на дефолтном логгере. |
| `func (e *Error) Error() string` | Метод, возвращающий строку сообщения ошибки. |
| `func (e *Error) WithOptions(ops ...Options) Errorer` | Изменение свойств `*Error` производится с помощью [функций-параметров](#Функции-параметры). |
| `func (e *Error) Severity() log.Severity` | Геттер получения важности ошибки. |
| `func (e *Error) Msg() string` | Геттер получения сообщения ошибки. |
| `func (e *Error) Sdump() string` | Геттер получения текстового дампа `*Error`. Может использоваться для отладки. Или в выводе форматированного сообщения, например, `fmt.Printf("%+v", err)`. |
| `func (e *Error) ErrorOrNil() error` | Геттер получения ошибки или `nil`. `Errorer` с типом `log.SeverityWarn` не является ошибкой; метод `ErrorOrNil` с таким типом вернет `nil`. |
| `func (e *Error) Operations() []string` | Геттер получения списка операций. |
| `func (e *Error) ErrorType() string` | Геттер получения типа ошибки. |
| `func (e *Error) Format(s fmt.State, verb rune)` | Функция форматирования для обработки строк с возможностью задания формата, например: `fmt.Printf`, `fmt.Sprintf`, `fmt.Fprintf`,.. |
| `func (e *Error) ID() string` | Геттер получения ID. |
| `func (e *Error) TranslateContext() *TranslateContext` | Геттер получения контекста перевода. |
| `func (e *Error) Localizer() *i18n.Localizer` | Геттер получения локализатора. |
| `func (e *Error) Log(l ...logcommon.Logger) ` | Метод логгирования. Выполнит логгирование ошибки с использованием логгера `l[0]`. |

### Функции-параметры

Параметризация `*Error` производится с помощью функций-параметров типа `type Options func(e *Error)`.

| Метод | Описание |
| ----- | -------- |
| `func SetFormatFn(fn FormatFn) Options` | Устанавливает функцию форматирования. Если значение `nil`, будет использоваться функция форматирования по-умолчанию. |
| `func SetMsg(msg string) Options` | Установить сообщение. |
| `func SetSeverity(severity log.Severity) Options` | Установить уровень важности сообщения. Доступные значения: `log.SeverityWarn`, `log.SeverityError`. |
| `func SetLocalizer(localizer *i18n.Localizer) Options ` | Установить локализатор для перевода. |
| `func SetTranslateContext(tctx *TranslateContext) Options` | Установить `*TranslateContext` для указанного языка. Используется для настройки дополнительных параметров, требуемых для корректного перевода. Например, `TranslateContext.PluralCount` позволяет установить множественное значение используемых в переводе объектов. |
| `func SetErrorType(etype string) Options` | Установить тип ошибки. Тип ошибки - `string`. |
| `func SetOperations(ops ...string) Options` | Установить список выполненных операций. |
| `func AppendOperations(ops ...string) Options` | Добавить операции к уже имеющемуся списку. Если список операций не существует, он будет создан. |
| `func SetContextInfo(ctxinf CtxMap) Options` | Задать контекст ошибки. |
| `func AppendContextInfo(key string, value interface{}) Options` | Добавить значения к уже имеющемуся контексту ошибки. Если контекст ошибки не существует, он будет создан. |
| `func SetID(id string) Options` | Установить ID ошибки. |

### Основные хелперы

Все хелперы работают с типом `error`.

| Хелпер | Описание |
| ------ | -------- |
| `func Is(err, target error) bool` | Обёртка над методом стандартной библиотеки `errors.Is`. |
| `func As(err error, target interface{}) bool` | Обёртка над методом стандартной библиотеки `errors.As`. |
| `func GetErrorType(err error) string` | Получить тип ошибки. Для НЕ `*Error` всегда будет "". |
| `func ErrorOrNil(err error) error` | Возможна обработка цепочки или одиночной ошибки. Если хотя бы одна ошибка в цепочке является ошибкой, то она будет возвращена в качестве результата. Важно: `*Error` c Severity `Warn` не является ошибкой. |
| `func Cast(err error) Errorer` | Преобразование типа `error` в `*Error`. |
| `func Append(errs ...error) error` | Создать цепочку ошибок. Допускается использование `nil` в аргументах. |
| `func Wrap(left error, right error) error` | Обернуть ошибку `left` ошибкой `right`, получив цепочку. Допускается использование `nil` в обоих аргументах. |
| `func Errors(err error) []error` | Получить список ошибок из цепочки. Вернет `nil`, при пустой цепочке. |
| `func Unwrap(err error) error` | Развернуть цепочку ошибок, получив первую. |
| `func UnwrapByID(err error, id string) Errorer` | Получить ошибку (`Errorer`) по ID. Вернет `nil`, если в случае провала поиска. |
| `func GetID(err error) (id string)` | Получить ID ошибки. Для НЕ `*Error` всегда будет "". |
| `func Contains(err error, id string) bool` | Проверить присутствует ли в цепочке ошибка с указанным ID. |

### Логгирование ошибки

Логгирование в пакете представлено интерфейсом:

```golang
type Logger interface {
    Warn(err error)
    Error(err error)
}
```

В пакете присутствует логгер по-умолчанию `gitlab.com/ovsinc/errors/log.DefaultLogger`.
Он установлен в значение

```golang
var DefaultLogger = golog.New(pkglog.New(os.Stderr, "ovsinc/errors ", pkglog.LstdFlags))
```

При необходимости его можно легко переопределить на более подходящее значение.

Для логгирования в `*Error` имеется метод `Log(l ...logcommon.Logger)`. Однако, приводить `error` к `*Error` каждый раз не требуется.

Для логгирования в пакете есть несколько хелперов.

| Хелпер | Описание |
| ------ | -------- |
| `func NewWithLog(msg string, ops ...Options) Errorer` | Функция произведет логгирование ошибки дефолтным логгером. |
| `func Log(err error, l ...logcommon.Logger)` | Функция произведет логгирование ошибки дефолтным логгером или логгером указанным в l (будет использоваться только первое значение). |
| `func AppendWithLog(errs ...error) error` | Хелпер создать цепочку ошибок., выполнит логгирование дефолтным логгером и вернет цепочку. |
| `func WrapWithLog(olderr error, err error) error` | Хелпер обернет `olderr` ошибкой `err`, выполнит логгирование дефолтным логгером и вернет цепочку. |

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

Для работы переводов нужно установить:

- `DefaultLocalizer`, тогда он будет использоваться для перевода всех ошибок;
- или локализатор для каждой отдельно взятой ошибки, используя функцию-параметр `*Error.SetLocalizer` при её создании.

Может оказаться удобным установить локализатор `DefaultLocalizer` для всего вашего приложения. Тогда, конечно, ваш локализатор должен содержать весь набор переводимых сообщений и настроен на использование требуемых языков.

В структуре `*Error` за перевод отвечают несколько свойств.

| Свойство | Тип |Назначение | Значение по-умолчанию |
| -------- | --- | --------- | --------------------- |
| translateContext | `*TranslateContext` | Дополнительная информация (контекст) для перевода. | `nil` |
| localizer  | `*i18n.Localizer` | Локализатор. Используется для выполнения переводов сообщения ошибки. | `nil` |

Для выполнения перевода ошибки требуется установить локализатор (если значение `DefaultLocalizer` не было установлено), используя функцию-параметр `SetLocalizer`.
Тогда при вызове метода `*Errors.Error()` будет выдана строка с переведенным сообщением.

В случае возникновения ошибки при переводе сообщения `*Error` будет выдана строка с оригинальным сообщением, без перевода.

Для ошибки `*Error` можно установить контекст перевода. Обычно это требуется для сложных сообщений, например, содержащих имена собственные или количественные значения. Для таких сообщений в составе контекста перевода необходимо установить шаблон `TemplateData map[string]interface{}`.
При использовании множественных форм в сообщении ошибки необходимо установить число в `PluralCount interface{}`.
Можно указать `DefaultMessage *i18n.Message`, если требуется указать значения перевода в случае ошибки перевода из файла.

См. подробности в пакете [i18n](https://github.com/nicksnyder/go-i18n).

Пример использование перевода в сообщении ошибки:

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
        errors.SetID("ErrEmailsUnreadMsg"),
        errors.SetErrorType("not found"),
        errors.SetLocalizer(i18n.NewLocalizer(bundle, "ru")),
        errors.SetTranslateContext(&errors.TranslateContext{
            TemplateData: map[string]interface{}{
                "Name":        "John Snow",
                "PluralCount": 5,
            },
            PluralCount: 5,
        }),
    )

    fmt.Printf("%v\n", err)
}
```

### Функции форматирования сообщения ошибки

В пакете представлены по паре (JSON, string) функций форматирования для единичного сообщения и цепочке сообщений ошибки.

Для цепочки сообщений изменение функции-форматера осуществляется через переменную `DefaultMultierrFormatFunc`. Для неё определено значение по-умолчанию `var DefaultMultierrFormatFunc = StringMultierrFormatFunc`.

Multierror-сообщения форматируются в пакете следующими функциями:

- для вывода в формате JSON - `func JSONMultierrFuncFormat(es []error) string`;
- для строкового вывода - `func StringMultierrFormatFunc(es []error) string`.

Для сообщений с типом `*Error` используются функции-форматеры типа `type FormatFn func(e *Error) string`. Задать требуемую функцию форматирования можно с помощью функции-параметра `SetFormatFn` в конструкторе или изменить это значение с помощью метода `WithOptions`.

В пакете представлены следующие функции-форматеры:

- для вывода в формате JSON - `func JSONFormat(e *Error) string`;
- для строкового вывода - `func StringFormat(e *Error) string`.

Внимание! При использовании форматирования цепочки сообщения `JSONMultierrFuncFormat` функция форматирование `*Error` по-умолчанию переключается на `JSONFormat`.

Все функций форматирования используют `github.com/valyala/bytebufferpool`, что хорошо сказывается на общей производительности и уменьшает потребление памяти.

[К оглавлению](#Оглавление)

## Список задач

- [Х] Повысить производительность функции форматирования для multierror;
- [ ] Повысить покрытие тестами;
- [ ] Более подробные комментарии для описания методов и функций;
- [ ] Перевод типа ошибки, операций, уровня опасности;
- [ ] Перевод README на en;
- [ ] Потокобезопасное управления ошибкой/ошибками;
- [ ] Выпуск на godoc.

[К оглавлению](#Оглавление)

## Лицензия

Код пакета распространяется под лицензией [Apache 2.0](http://directory.fsf.org/wiki/License:Apache2.0).
