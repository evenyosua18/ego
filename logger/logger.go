package logger

type (
	Logger interface {
		Debug(msg string, fields ...Field)
		Info(msg string, fields ...Field)
		Warn(msg string, fields ...Field)
		Error(err error, fields ...Field)
		Fatal(err error, fields ...Field)
		With(fields ...Field) Logger
	}

	Field struct {
		Key   string
		Value any
	}

	Level int
)

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)
