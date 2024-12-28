package actor

import (
	"fmt"
	"os"
)

type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

var logger Logger = defaultLogger{}

type defaultLogger struct{}

func (defaultLogger) Info(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Printf("\n")
	}
}

func (defaultLogger) Error(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func (defaultLogger) Fatal(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	if format[len(format)-1] != '\n' {
		fmt.Fprintf(os.Stderr, "\n")
	}
}
