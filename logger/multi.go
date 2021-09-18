package logger

type MultiLogger struct {
	loggers []ILogger
}

func NewMultiLogger(loggers ...ILogger) ILogger {
	ml := new(MultiLogger)
	ml.loggers = loggers
	return ml
}

func (ml *MultiLogger) Debug(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Debug(format, args...)
	}
}

func (ml *MultiLogger) Info(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Info(format, args...)
	}
}

func (ml *MultiLogger) Fatal(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Fatal(format, args...)
	}
}

func (ml *MultiLogger) Error(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Error(format, args...)
	}
}
