package logging

import (
	"fmt"
	"io"
)

// Logger includes output interface.
type Logger struct {
	output io.Writer
}

// NewLogger returns new Logger.
func NewLogger(output io.Writer) Logger {
	return Logger{
		output: output,
	}
}

// Debug outputs debug log.
func (l *Logger) Debug(msg string) {
	fmt.Fprintf(l.output, "[DEBUG] %s", msg)
}

// Info outputs info log.
func (l *Logger) Info(msg string) {
	fmt.Fprintf(l.output, "[INFO] %s", msg)
}

// Error outputs error log.
func (l *Logger) Error(msg string) {
	fmt.Fprintf(l.output, "[ERROR] %s", msg)
}
