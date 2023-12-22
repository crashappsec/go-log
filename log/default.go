package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/crashappsec/go-log"
)

var defaultLogger *log.Logger

func init() {
	defaultLogger = log.NewLogger().WithOptions(
		zap.AddCallerSkip(1), // default logger needs to skip one more caller
	)
}

func With(fields ...zap.Field) *log.Logger {
	return defaultLogger.With(fields...)
}

func Log(level zapcore.Level, msg string, i ...zap.Field) {
	defaultLogger.Log(level, msg, i...)
}

func Print(msg string, i ...zap.Field) {
	defaultLogger.Info(msg, i...)
}

func Debug(msg string, i ...zap.Field) {
	defaultLogger.Debug(msg, i...)
}

func Info(msg string, i ...zap.Field) {
	defaultLogger.Info(msg, i...)
}

func Warn(msg string, i ...zap.Field) {
	defaultLogger.Warn(msg, i...)
}

func Error(msg string, i ...zap.Field) {
	defaultLogger.Error(msg, i...)
}

func Fatal(msg string, i ...zap.Field) {
	defaultLogger.Fatal(msg, i...)
}
