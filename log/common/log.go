package common

// Logger упрощенный интерфейс логгирования
type Logger interface {
	Warn(err error)
	Error(err error)
}
