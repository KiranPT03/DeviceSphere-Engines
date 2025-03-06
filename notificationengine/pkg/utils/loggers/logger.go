package logger

import (
	"fmt"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.CallerKey = "caller" // Add caller key

		var err error
		instance, err = config.Build(zap.AddCallerSkip(1), zap.AddCaller()) // Add caller info
		if err != nil {
			panic(err)
		}
	})

	return instance
}

func logWithCaller(logFunc func(string, ...zap.Field), format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)

	pc, _, line, ok := runtime.Caller(2) // Get caller info
	funcName := "unknown"
	if ok {
		f := runtime.FuncForPC(pc)
		if f != nil {
			funcName = f.Name()
		}
	}
	logFunc(msg, zap.String("func", funcName), zap.Int("line", line))
}

func Debug(format string, v ...interface{}) {
	logWithCaller(GetLogger().Debug, format, v...)
}

func Info(format string, v ...interface{}) {
	logWithCaller(GetLogger().Info, format, v...)
}

func Warn(format string, v ...interface{}) {
	logWithCaller(GetLogger().Warn, format, v...)
}

func Error(format string, v ...interface{}) {
	logWithCaller(GetLogger().Error, format, v...)
}

func Fatal(format string, v ...interface{}) {
	logWithCaller(GetLogger().Fatal, format, v...)
}
