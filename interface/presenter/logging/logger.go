package logging

import (
	"fmt"
	"io"
)

// Logger includes output interface.
type Logger struct {
	Output io.Writer
}

// Debug outputs debug log.
func (l *Logger) Debug(msg string) {
	fmt.Fprintf(l.Output, "[DEBUG] %s", msg)
}

// Info outputs info log.
func (l *Logger) Info(msg string) {
	fmt.Fprintf(l.Output, "[INFO] %s", msg)
}

// Error outputs error log.
func (l *Logger) Error(msg string) {
	fmt.Fprintf(l.Output, "[ERROR] %s", msg)
}
