package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger
func Init(debug bool) error {
	var err error
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// Get returns the logger instance
func Get() *zap.Logger {
	if log == nil {
		// Initialize with production config if not initialized
		_ = Init(false)
	}
	return log
}

// Sync flushes any buffered log entries
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

// LogError logs an error with context
func LogError(err error, msg string, fields map[string]interface{}) {
	if log == nil {
		return
	}
	zapFields := make([]zapcore.Field, 0, len(fields)+1)
	zapFields = append(zapFields, zap.Error(err))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	log.Error(msg, zapFields...)
}

// LogInfo logs an info message with context
func LogInfo(msg string, fields map[string]interface{}) {
	if log == nil {
		return
	}
	zapFields := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	log.Info(msg, zapFields...)
}

// LogDebug logs a debug message with context
func LogDebug(msg string, fields map[string]interface{}) {
	if log == nil {
		return
	}
	zapFields := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	log.Debug(msg, zapFields...)
} 