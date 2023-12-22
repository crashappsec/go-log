// log defines a logging intefrace and wraps uber-go/zap logger
package log

import (
	"os"
	"strings"

	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// PanicLevel logs message and then calls `panic``
	PanicLevel = zap.PanicLevel
	// FatalLevel logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = zap.FatalLevel
	// ErrorLevel should be used for any error
	ErrorLevel = zap.ErrorLevel
	// WarnLevel to be used on non-erroneous cases that could however be undesired
	WarnLevel = zap.WarnLevel
	// InfoLevel to be used for output directed at the user
	InfoLevel = zap.InfoLevel
	// DebugLevel to only be enabled when debugging
	DebugLevel = zap.DebugLevel
)

var (
	logFormat = strings.ToLower(os.Getenv("LOG_FORMAT"))
)

type Level = zapcore.Level
type Field = zap.Field

type Logger struct {
	zap  *zap.Logger
	atom *zap.AtomicLevel
	// since zap repeats key/values for each With statement, allocate this
	// per-logger global map to only ever register a field via With() the first
	// time. https://github.com/uber-go/zap/issues/622
	context map[string]zap.Field
}

func envLogLevel() zapcore.Level {
	var level zapcore.Level
	envLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch envLevel {
	case "FATAL":
		level = FatalLevel
	case "ERROR":
		level = ErrorLevel
	case "WARN":
		level = WarnLevel
	case "DEBUG":
		level = DebugLevel
	case "INFO":
		fallthrough
	default:
		level = InfoLevel
	}
	return level
}

func IsConsole() bool {
	return logFormat == "console"
}

func IsJson() bool {
	return logFormat != "console"
}

func NewLogger() *Logger {
	var encoder zapcore.Encoder

	if IsConsole() {
		config := prettyconsole.NewEncoderConfig()
		config.CallerKey = "logger"
		encoder = prettyconsole.NewEncoder(config)
	} else {
		config := zap.NewProductionEncoderConfig()
		config.TimeKey = "timestamp"
		config.CallerKey = "logger"
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(config)
	}

	opts := []zap.Option{
		zap.AddCallerSkip(1), // traverse call depth for more useful log lines
		zap.AddCaller(),
	}
	atom := zap.NewAtomicLevelAt(envLogLevel())
	logger := zap.New(zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		atom,
	), opts...)

	return &Logger{
		zap:     logger,
		atom:    &atom,
		context: map[string]zap.Field{},
	}
}

func (l *Logger) SetLevel(level zapcore.Level) {
	if l.atom == nil {
		return
	}
	l.atom.SetLevel(level)
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		zap:     l.zap,
		atom:    l.atom,
		context: l.mergeFields(fields),
	}
}

func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	logger := l.zap.WithOptions(opts...)
	return &Logger{
		zap:     logger,
		atom:    l.atom,
		context: l.context,
	}
}

func (l *Logger) Log(level zapcore.Level, msg string, i ...zap.Field) {
	l.zap.Log(level, msg, l.fields(i)...)
}

func (l *Logger) Print(msg string, i ...zap.Field) {
	l.zap.Info(msg, l.fields(i)...)
}

func (l *Logger) Debug(msg string, i ...zap.Field) {
	l.zap.Debug(msg, l.fields(i)...)
}

func (l *Logger) Info(msg string, i ...zap.Field) {
	l.zap.Info(msg, l.fields(i)...)
}

func (l *Logger) Warn(msg string, i ...zap.Field) {
	l.zap.Warn(msg, l.fields(i)...)
}

func (l *Logger) Error(msg string, i ...zap.Field) {
	l.zap.Error(msg, l.fields(i)...)
}

func (l *Logger) Fatal(msg string, i ...zap.Field) {
	l.zap.Fatal(msg, l.fields(i)...)
}

func (l *Logger) mergeFields(fields []zap.Field) map[string]zap.Field {
	merged := make(map[string]zap.Field, len(l.context))
	// copy fields from logger current logger
	for k, v := range l.context {
		merged[k] = v
	}
	// overwrite context fields with any provided fields
	for _, v := range fields {
		merged[v.Key] = v
	}
	return merged
}

func (l *Logger) fields(fields []zap.Field) []zap.Field {
	merged := l.mergeFields(fields)
	values := make([]zap.Field, len(merged))
	i := 0
	for _, v := range merged {
		values[i] = v
		i += 1
	}
	return values
}
