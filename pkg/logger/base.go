package logger

import (
	"context"
	"log/slog"
)

var defaultLogger GLBLogger = createBaseLogger()

type basicLogger struct {
	logger *slog.Logger
}

func createBaseLogger() GLBLogger {
	return &basicLogger{logger: slog.Default()}
}

func (b *basicLogger) Info(msg string, args ...any) {
	b.logger.Info(msg, args...)
}

func (b *basicLogger) Warn(msg string, args ...any) {
	b.logger.Warn(msg, args...)
}

func (b *basicLogger) Error(msg string, args ...any) {
	b.logger.Error(msg, args...)
}

func (b *basicLogger) Debug(msg string, args ...any) {
	b.logger.Debug(msg, args...)
}

func (b *basicLogger) Fatal(msg string, args ...any) {
	b.logger.Error(msg, args...)
	panic(msg)
}

func (b *basicLogger) SetLevel(level GLBLogLevel) {
	b.logger.Enabled(context.Background(), glbLevelToSlogLevel(level))
}

func (b *basicLogger) WithParam(key string, value any) GLBLogger {
	// Implement the logic to add a parameter to the logger
	return &basicLogger{logger: b.logger.With(key, value)}
}

func glbLevelToSlogLevel(level GLBLogLevel) slog.Level {
	switch level {
	case GlbDebug:
		return slog.LevelDebug
	case GlbInfo:
		return slog.LevelInfo
	case GlbWarning:
		return slog.LevelWarn
	case GlbError:
		return slog.LevelError
	case GlbFatal:
		return slog.LevelError // FATAL is treated as ERROR in slog
	default:
		return slog.LevelInfo // Default to INFO if unknown level
	}
}
