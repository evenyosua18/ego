package logger

var globalLogger Logger = &defaultLogger{}

func SetLogger(l Logger) {
	if l != nil {
		globalLogger = l
	}
}

func L() Logger {
	return globalLogger
}

func Debug(msg string, fields ...Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(err error, fields ...Field) {
	globalLogger.Error(err, fields...)
}

func Fatal(err error, fields ...Field) {
	globalLogger.Fatal(err, fields...)
}

// helper for redundant debug call

func DebugQuery(query, args any) {
	Debug("query statement", Field{Key: "query", Value: query}, Field{Key: "args", Value: args})
}
