package logger

type GLBLogger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)

	SetLevel(level Level)
	WithParam(key string, value any) GLBLogger
}

type Level string

const (
	GlbDebug   Level = "debug"
	GlbInfo    Level = "info"
	GlbWarning Level = "warning"
	GlbError   Level = "error"
	GlbFatal   Level = "fatal"
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
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}
func WithParam(key string, value any) GLBLogger {
	return defaultLogger.WithParam(key, value)
}

func SetDefaultLogger(logger GLBLogger) {
	defaultLogger = logger
}

func GetDefaultLogger() GLBLogger {
	return defaultLogger
}

func ParseLogLevel(level string) Level {
	switch level {
	case "debug":
		return GlbDebug
	case "info":
		return GlbInfo
	case "warn":
		return GlbWarning
	case "error":
		return GlbError
	case "fatal":
		return GlbFatal
	default:
		return GlbInfo
	}
}
