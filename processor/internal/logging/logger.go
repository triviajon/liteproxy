package logging

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger   *zap.Logger
	sugared  *zap.SugaredLogger
	initOnce sync.Once
)

// Init initializes the logger
func Init() error {
	var err error
	initOnce.Do(func() {
		config := zap.NewProductionConfig()
		config.Encoding = "console"
		config.Development = true
		config.DisableStacktrace = true
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Colored levels
		config.EncoderConfig.TimeKey = ""                                   // Remove timestamp for cleaner output
		config.EncoderConfig.CallerKey = "caller"                           // Show caller (file:line)
		config.OutputPaths = []string{"stdout"}

		logger, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
		if err != nil {
			return
		}
		sugared = logger.Sugar()
	})
	return err
}

// get ensures the logger is initialized and returns it.
func get() *zap.Logger {
	if logger == nil {
		_ = Init()
	}
	return logger
}

// getSugared ensures the logger is initialized and returns the sugared logger.
func getSugared() *zap.SugaredLogger {
	if sugared == nil {
		_ = Init()
	}
	return sugared
}

// SetLevel sets the log level globally.
func SetLevel(level zapcore.Level) {
	if logger == nil {
		_ = Init()
	}
	logger = logger.WithOptions(zap.IncreaseLevel(level))
	sugared = logger.Sugar()
}

// ConfigureFromEnv reads the LP_LOG_LEVEL environment variable and sets the log level.
// Valid values are: DEBUG, INFO, WARN, ERROR. Defaults to INFO if not set or invalid.
func ConfigureFromEnv() {
	logLevel := os.Getenv("LP_LOG_LEVEL")
	if logLevel == "" {
		return // Use default (INFO)
	}

	logLevel = strings.ToUpper(strings.TrimSpace(logLevel))
	var level zapcore.Level
	switch logLevel {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	default:
		Warnf("Invalid LP_LOG_LEVEL: %s (use DEBUG, INFO, WARN, or ERROR)", logLevel)
		return
	}
	SetLevel(level)
}

// Debugf logs a debug message.
func Debugf(format string, args ...interface{}) {
	getSugared().Debugf(format, args...)
}

// Debug logs a debug message with structured fields.
func Debug(msg string, fields ...zap.Field) {
	get().Debug(msg, fields...)
}

// Infof logs an info message.
func Infof(format string, args ...interface{}) {
	getSugared().Infof(format, args...)
}

// Info logs an info message with structured fields.
func Info(msg string, fields ...zap.Field) {
	get().Info(msg, fields...)
}

// Warnf logs a warning message.
func Warnf(format string, args ...interface{}) {
	getSugared().Warnf(format, args...)
}

// Warn logs a warning message with structured fields.
func Warn(msg string, fields ...zap.Field) {
	get().Warn(msg, fields...)
}

// Errorf logs an error message.
func Errorf(format string, args ...interface{}) {
	getSugared().Errorf(format, args...)
}

// Error logs an error message with structured fields.
func Error(msg string, fields ...zap.Field) {
	get().Error(msg, fields...)
}

// Fatalf logs an error message and exits.
func Fatalf(format string, args ...interface{}) {
	getSugared().Errorf(format, args...)
	os.Exit(1)
}
