package logging

import (
	"errors"
	"fmt"
	"io"
	"time"
)

// Logger includes output interface.
type Logger struct {
	output   io.Writer
	location *time.Location
}

// NewLogger returns new Logger.
func NewLogger(output io.Writer, location *time.Location) (Logger, error) {
	if location == nil {
		return Logger{}, errors.New("Time location is nil")
	}
	return Logger{
		output:   output,
		location: location,
	}, nil
}

// Debug outputs debug log.
func (l *Logger) Debug(msg string) {
	fmt.Fprintf(l.output, "[DEBUG] %v %s\n", time.Now().In(l.location), msg)
}

// Info outputs info log.
func (l *Logger) Info(msg string) {
	fmt.Fprintf(l.output, "[INFO] %v %s\n", time.Now().In(l.location), msg)
}

// Error outputs error log.
func (l *Logger) Error(msg string) {
	fmt.Fprintf(l.output, "[ERROR] %v %s\n", time.Now().In(l.location), msg)
}
