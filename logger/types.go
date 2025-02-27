package logger

type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
)

type Logger interface {
	New() (Logger, error)
	Close() error
	Sync() error

	Info(args ...interface{})
	Infof(msg string, args ...interface{})
	Debug(args ...interface{})
	Debugf(msg string, args ...interface{})
	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})
	Error(args ...interface{})
	Errorf(msg string, args ...interface{})
	Panic(args ...interface{})
	Panicf(msg string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type Encode interface {
}

type Option struct {
	ModuleName   string
	OutPath      string
	MaxSize      int64
	MaxRetainDay int
	Code         Encode
}
