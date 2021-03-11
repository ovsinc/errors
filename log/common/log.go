package common

// Logger интерфейс логгера.
type Logger interface {
	Warn(err error)
	Error(err error)
}
