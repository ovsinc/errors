package log

import (
	pkglog "log"
	"os"

	"gitlab.com/ovsinc/errors/log/golog"
)

var DefaultLogger = golog.New(pkglog.New(os.Stderr, "ovsinc/errors ", pkglog.LstdFlags))
