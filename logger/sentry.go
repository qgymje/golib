package logger

import (
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type SentryLogger struct {
	client  *raven.Client
	appName string
}

func NewSentryLogger(appName string, dsn string) *SentryLogger {
	l := new(SentryLogger)
	l.client = raven.DefaultClient
	l.client.SetDSN(dsn)

	l.appName = appName
	return l
}

func (l *SentryLogger) Error(format string, args ...interface{}) {
	tag := map[string]string{
		"err_level": "error",
		"app":       l.appName,
	}
	err := fmt.Errorf(format, args...)
	err = errors.Wrap(err, "")
	l.client.CaptureMessage(fmt.Sprintf("%+v", err), tag)
}

func (l *SentryLogger) Fatal(format string, args ...interface{}) {
	tag := map[string]string{
		"err_level": "fatal",
		"app":       l.appName,
	}
	l.client.CaptureMessage(fmt.Sprintf(format, args...), tag)
}

func (l *SentryLogger) Info(format string, args ...interface{}) {
}

func (l *SentryLogger) Debug(format string, args ...interface{}) {
}
