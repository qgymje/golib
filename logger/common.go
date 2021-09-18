package logger

type ILogger interface {
	Fatal(format string, args ...interface{})
	Error(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
}
