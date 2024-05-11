package logger

import (

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log   *zap.Logger
	sugar *zap.SugaredLogger
}

// func (z *ZapLogger) New() {
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync() // flushes buffer, if any
// 	sugar := logger.Sugar()
// 	sugar.Infow("failed to fetch URL",
// 		// Structured context as loosely typed key-value pairs.
// 		"url", "url",
// 		"attempt", 3,
// 		"backoff", time.Second,
// 	)
// 	sugar.Infof("Failed to fetch URL: %s", "url")
// }


func (z *ZapLogger) New() (Logger, error) {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	lvl, _ := zapcore.ParseLevel(string(DebugLevel))
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, lvl)
	l := zap.New(core, zap.AddCallerSkip(1))
	s := l.Sugar()
	ll := &ZapLogger{
		log:   l,
		sugar: s,
	}
	return ll, nil
}


func (z *ZapLogger) Info(args ...interface{}) {
	z.sugar.Info(args...)
}

func (z *ZapLogger) Infof(msg string, args ...interface{}) {
	z.sugar.Infof(msg, args...)
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.sugar.Debug(args...)
}

func (z *ZapLogger) Debugf(msg string, args ...interface{}) {
	z.sugar.Debugf(msg, args...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.sugar.Warn(args...)
}

func (z *ZapLogger) Warnf(msg string, args ...interface{}) {
	z.sugar.Warnf(msg, args...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.sugar.Error(args...)
}

func (z *ZapLogger) Errorf(msg string, args ...interface{}) {
	z.sugar.Errorf(msg, args...)
}

func (z *ZapLogger) Panic(args ...interface{}) {
	z.sugar.Panic(args...)
}

func (z *ZapLogger) Panicf(msg string, args ...interface{}) {
	z.sugar.Panicf(msg, args...)
}

func (z *ZapLogger) Fatal(args ...interface{}) {
	z.sugar.Fatal(args...)
}

func (z *ZapLogger) Fatalf(msg string, args ...interface{}) {
	z.sugar.Fatalf(msg, args...)
}

func (z *ZapLogger) Sync() error {
	return z.sugar.Sync()
}


func (z *ZapLogger) Close() error {
	return z.sugar.Sync()
}
