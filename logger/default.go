package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type defaultLogger struct {
	fields []Field
}

func (d *defaultLogger) Debug(msg string, fields ...Field) {
	d.print("DEBUG", msg, fields...)
}

func (d *defaultLogger) Info(msg string, fields ...Field) {
	d.print("INFO", msg, fields...)
}

func (d *defaultLogger) Warn(msg string, fields ...Field) {
	d.print("WARN", msg, fields...)
}

func (d *defaultLogger) Error(err error, fields ...Field) {
	d.print("ERROR", err.Error(), fields...)
}

func (d *defaultLogger) Fatal(err error, fields ...Field) {
	d.print("FATAL", err.Error(), fields...)
	os.Exit(1)
}

func (d *defaultLogger) With(fields ...Field) Logger {
	merged := make([]Field, 0, len(d.fields)+len(fields))
	merged = append(merged, d.fields...)
	merged = append(merged, fields...)
	return &defaultLogger{fields: merged}
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
