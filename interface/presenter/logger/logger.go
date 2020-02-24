package presenter

import (
	"fmt"
	"io"
)

type Logger struct {
	Output io.Writer
}

func (l *Logger) Debug(msg string) {
	fmt.Fprintf(l.Output, "[DEBUG] %s", msg)
}

func (l *Logger) Info(msg string) {
	fmt.Fprintf(l.Output, "[INFO] %s", msg)
}

func (l *Logger) Error(msg string) {
	fmt.Fprintf(l.Output, "[ERROR] %s", msg)
}
