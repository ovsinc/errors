package log

import (
	pkglog "log"
	"os"

	"gitlab.com/ovsinc/errors/log/common"
	"gitlab.com/ovsinc/errors/log/golog"
)

// DefaultLogger логгер, используемый по умолчанию.
// Можно переопределить.
var DefaultLogger common.Logger = golog.New(pkglog.New(os.Stderr, "ovsinc/errors ", pkglog.LstdFlags))
