package httplib

import (
	"fmt"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type emptyLogger struct {
}

func (emptyLogger) Info(args ...interface{}) {
	if !debug {
		return
	}
	fmt.Println(args...)
}

func (emptyLogger) Infof(format string, args ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(format, args...)
}

func (emptyLogger) Errorf(format string, args ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(format, args...)
}

var defaultLogger Logger = emptyLogger{}
var debug = false

func SetDefaultLogger(l Logger) {
	defaultLogger = l
}

func SetDebug(b bool) {
	debug = b
}
