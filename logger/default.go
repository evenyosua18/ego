package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type defaultLogger struct {
	fields []Field
	level  Level
}

func (d *defaultLogger) Debug(msg string, fields ...Field) {
	if d.level <= LevelDebug {
		d.print("DEBUG", msg, fields...)
	}
}

func (d *defaultLogger) Info(msg string, fields ...Field) {
	if d.level <= LevelInfo {
		d.print("INFO", msg, fields...)
	}
}

func (d *defaultLogger) Warn(msg string, fields ...Field) {
	if d.level <= LevelWarn {
		d.print("WARN", msg, fields...)
	}
}

func (d *defaultLogger) Error(err error, fields ...Field) {
	if d.level <= LevelError {
		d.print("ERROR", err.Error(), fields...)
	}
}

func (d *defaultLogger) Fatal(err error, fields ...Field) {
	if d.level <= LevelFatal {
		d.print("FATAL", err.Error(), fields...)
		os.Exit(1)
	}
}

func (d *defaultLogger) With(fields ...Field) Logger {
	merged := make([]Field, 0, len(d.fields)+len(fields))
	merged = append(merged, d.fields...)
	merged = append(merged, fields...)
	return &defaultLogger{
		fields: merged,
		level:  d.level,
	}
}

func (d *defaultLogger) print(level, msg string, callFields ...Field) {
	allFields := make([]Field, 0, len(d.fields)+len(callFields))
	allFields = append(allFields, d.fields...)
	allFields = append(allFields, callFields...)

	fieldStrs := make([]string, 0, len(allFields))
	for _, f := range allFields {
		fieldStrs = append(fieldStrs, fmt.Sprintf("%s = %v", f.Key, f.Value))
	}
	fieldText := ""
	if len(fieldStrs) > 0 {
		fieldText = " | " + strings.Join(fieldStrs, " ")
	}

	log.Printf("[%s] %s%s", level, msg, fieldText)
}

func ParseLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}

func NewDefaultLogger(level Level) Logger {
	return &defaultLogger{
		level: level,
	}
}
