package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initialises the global Zap logger.
// Call this once at application startup.
func Init() error {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

// Get returns the global logger instance.
// Panics if Init has not been called first.
func Get() *zap.Logger {
	if log == nil {
		panic("logger: Get() called before Init()")
	}
	return log
}

// Sync flushes any buffered log entries. Call this on application shutdown.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
