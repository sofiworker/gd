package logger

import (
	"sync"
	"sync/atomic"

	"github.com/chuck1024/gd/v2/gerr"
)

var (
	logOnce sync.Once
	v       atomic.Value
)

func init() {
	vv := v.Load()
	if vv == nil {
		z := ZapLogger{}
		l, err := z.New()
		if err != nil {
			panic(gerr.InitDefaultLoggerFailed)
		}
		v.Store(l)
	}
}

func SetDefaultLogger(l Logger) {
	logOnce.Do(func() {
		v.Store(l)
	})
}

func DefaultLogger() Logger {
	vv := v.Load()
	if vv != nil {
		if l, ok := vv.(Logger); ok {
			return l
		}
	}
	panic("not found any default logger")
}

func Info(args ...interface{}) {
	DefaultLogger().Info(args...)
}

func Infof(msg string, args ...interface{}) {
	DefaultLogger().Infof(msg, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger().Debug(args...)
}

func Debugf(msg string, args ...interface{}) {
	DefaultLogger().Debugf(msg, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger().Warn(args...)
}

func Warnf(msg string, args ...interface{}) {
	DefaultLogger().Warnf(msg, args...)
}

func Error(args ...interface{}) {
	DefaultLogger().Error(args...)
}

func Errorf(msg string, args ...interface{}) {
	DefaultLogger().Errorf(msg, args...)
}

func Panic(args ...interface{}) {
	DefaultLogger().Panic(args...)
}

func Panicf(msg string, args ...interface{}) {
	DefaultLogger().Panicf(msg, args...)
}


func Fatal(args ...interface{}) {
	DefaultLogger().Fatal(args...)
}

func Fatalf(msg string, args ...interface{}) {
	DefaultLogger().Fatalf(msg, args...)
}

func Sync() error {
	return DefaultLogger().Sync()
}
