package loggers

import (
	"fmt"
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

		var err error
		instance, err = config.Build()
		if err != nil {
			panic(err)
		}
	})

	return instance
}

func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Debug(msg)
}

func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Info(msg)
}

func Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Warn(msg)
}

func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Error(msg)
}

func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	GetLogger().Fatal(msg)
}
