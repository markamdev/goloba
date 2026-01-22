package logger

import (
	"log/slog"
	"os"
)

func NewBasicLogger() GLBLogger {

	lvl := &slog.LevelVar{}
	lvl.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{
		AddSource: false, // TODO consider making this configurable
		Level:     lvl,
	}
	return &basicLogger{
		logger:   slog.New(slog.NewTextHandler(os.Stdout, opts)),
		logLevel: lvl,
	}
}

type basicLogger struct {
	logger   *slog.Logger
	logLevel *slog.LevelVar
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
	os.Exit(1)
}

func (b *basicLogger) SetLevel(level Level) {
	b.logLevel.Set(glbLevelToSlogLevel(level))
}

func (b *basicLogger) WithParam(key string, value any) GLBLogger {
	// Implement the logic to add a parameter to the logger
	return &basicLogger{logger: b.logger.With(key, value)}
}

func glbLevelToSlogLevel(level Level) slog.Level {
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
		return slog.LevelDebug // Default to INFO if unknown level
	}
}
