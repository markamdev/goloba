package logger

type GLBLogger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)

	SetLevel(level GLBLogLevel)
	WithParam(key string, value any) GLBLogger
}

type GLBLogLevel int

const (
	GlbDebug GLBLogLevel = iota
	GlbInfo
	GlbWarning
	GlbError
	GlbFatal
)

func Debug(message string, args ...any) {
	defaultLogger.Debug(message, args...)
}
func Info(message string, args ...any) {
	defaultLogger.Info(message, args...)
}
func Warn(message string, args ...any) {
	defaultLogger.Warn(message, args...)
}
func Error(message string, args ...any) {
	defaultLogger.Error(message, args...)
}
func Fatal(message string, args ...any) {
	defaultLogger.Fatal(message, args...)
}
func SetLevel(level GLBLogLevel) {
	defaultLogger.SetLevel(level)
}
func WithParam(key string, value any) GLBLogger {
	return defaultLogger.WithParam(key, value)
}

func SetDefaultLogger(logger GLBLogger) {
	defaultLogger = logger
}

func NewBaseLogger() GLBLogger {
	return createBaseLogger()
}
