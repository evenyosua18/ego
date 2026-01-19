package logger

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestDefaultLogger_Level(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil) // Reset to default

	tests := []struct {
		name          string
		level         Level
		action        func(l Logger)
		expectedLevel string
		shouldLog     bool
	}{
		{
			name:  "Debug level - log debug",
			level: LevelDebug,
			action: func(l Logger) {
				l.Debug("test debug")
			},
			expectedLevel: "DEBUG",
			shouldLog:     true,
		},
		{
			name:  "Info level - skip debug",
			level: LevelInfo,
			action: func(l Logger) {
				l.Debug("test debug")
			},
			shouldLog: false,
		},
		{
			name:  "Info level - log info",
			level: LevelInfo,
			action: func(l Logger) {
				l.Info("test info")
			},
			expectedLevel: "INFO",
			shouldLog:     true,
		},
		{
			name:  "Warn level - skip info",
			level: LevelWarn,
			action: func(l Logger) {
				l.Info("test info")
			},
			shouldLog: false,
		},
		{
			name:  "Error level - log error",
			level: LevelError,
			action: func(l Logger) {
				l.Error(fmt.Errorf("test error"))
			},
			expectedLevel: "ERROR",
			shouldLog:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			l := NewDefaultLogger(tt.level)
			
			tt.action(l)

			output := buf.String()
			if tt.shouldLog {
				if output == "" {
					t.Errorf("Expected output, got empty")
				}
				if !strings.Contains(output, "["+tt.expectedLevel+"]") {
					t.Errorf("Expected level %s, got %s", tt.expectedLevel, output)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no output, got %s", output)
				}
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"error", LevelError},
		{"fatal", LevelFatal},
		{"unknown", LevelInfo}, // Default
	}

	for _, tt := range tests {
		if got := ParseLevel(tt.input); got != tt.expected {
			t.Errorf("ParseLevel(%s) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}
