# multilog

[![codecov](https://codecov.io/gh/ovsinc/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/ovsinc/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/ovsinc/multilog)](https://goreportcard.com/report/github.com/ovsinc/multilog)
[![Go Reference](https://pkg.go.dev/badge/github.com/ovsinc/multilog.svg)](https://pkg.go.dev/github.com/ovsinc/multilog)
[![license](https://img.shields.io/badge/license-MIT-green)](https://github.com/ovsinc/multilog/blob/main/LICENSE)

[circleci]: https://app.circleci.com/pipelines/github/ovsinc/multilog
[godocs]: https://pkg.go.dev/github.com/ovsinc/multilog

Package multilog is a simple logging wrapper for common logging applications.
The following loggers are supported:

* [logrus](https://pkg.go.dev/github.com/sirupsen/logrus)
* [golog](https://pkg.go.dev/log)
* [log15](https://pkg.go.dev/github.com/inconshreveable/log15)
* [journald](https://pkg.go.dev/github.com/coreos/go-systemd/journal)
* [syslog](https://pkg.go.dev/log/syslog)
* [zap](https://pkg.go.dev/go.uber.org/zap)

It is possible to combine supported loggers in a chain.
