package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
    Info(args ...interface{})
    Warn(args ...interface{})
    Debug(args ...interface{})
    Error(args ...interface{})
    Fatal(args ...interface{})
}

func NewLogger() *logrus.Logger {
    log := logrus.New()

    defer log.Info("Logger has been init")
    return log
}
