package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init initializes the logger
func Init(env string) {
	once.Do(func() {
		var config zap.Config

		if env == "production" {
			config = zap.NewProductionConfig()
			config.EncoderConfig.TimeKey = "timestamp"
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		var err error
		log, err = config.Build()
		if err != nil {
			os.Exit(1)
		}
	})
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		Init("development")
	}
	return log
}

// Info logs info level message
func Info(msg string, fields ...zapcore.Field) {
	GetLogger().Info(msg, fields...)
}

// Error logs error level message
func Error(msg string, fields ...zapcore.Field) {
	GetLogger().Error(msg, fields...)
}

// Debug logs debug level message
func Debug(msg string, fields ...zapcore.Field) {
	GetLogger().Debug(msg, fields...)
}

// Warn logs warn level message
func Warn(msg string, fields ...zapcore.Field) {
	GetLogger().Warn(msg, fields...)
}

// Fatal logs fatal level message and exits
func Fatal(msg string, fields ...zapcore.Field) {
	GetLogger().Fatal(msg, fields...)
}

// With creates a child logger with additional fields
func With(fields ...zapcore.Field) *zap.Logger {
	return GetLogger().With(fields...)
}
