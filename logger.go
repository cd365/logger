package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type Level = slog.Level

const (
	LevelAll   Level = -6
	LevelTrace Level = -5
	LevelDebug       = slog.LevelDebug
	LevelInfo        = slog.LevelInfo
	LevelWarn        = slog.LevelWarn
	LevelError       = slog.LevelError
	LevelFatal Level = 9
	LevelOff   Level = 10
)

var LevelMap = map[Level]string{
	LevelAll:   "ALL",
	LevelTrace: "TRACE",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelOff:   "OFF",
}

const (
	LogDefaultSkip = 3
)

type Logger struct {
	// HandlerOptions Log handler options.
	HandlerOptions *slog.HandlerOptions

	// Logger Log object.
	Logger *slog.Logger

	// LevelVar Dynamically adjust log level.
	LevelVar *slog.LevelVar

	// skip Number of stack frames to skip before recording.
	skip int
}

func New(level Level, skip ...int) *Logger {
	levelVar := slog.LevelVar{}
	levelVar.Set(level)
	handlerOptions := &slog.HandlerOptions{
		// With source.
		AddSource: true,

		// Support dynamic setting of log level.
		Level: &levelVar,

		// Modify the Attr key-value pair in the log (that is, the key/value attached to the log record).
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				levelValue := a.Value.Any().(slog.Level)
				levelLabel := levelValue.String()
				switch levelValue {
				case LevelTrace:
					levelLabel = LevelMap[levelValue]
				case LevelFatal:
					levelLabel = LevelMap[levelValue]
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	skipValue := LogDefaultSkip
	for i := len(skip) - 1; i >= 0; i-- {
		if skip[i] > 0 {
			skipValue = skip[i]
			break
		}
	}
	return &Logger{
		HandlerOptions: handlerOptions,
		Logger:         logger,
		LevelVar:       &levelVar,
		skip:           skipValue,
	}
}

func (s *Logger) log(ctx context.Context, level Level, msg string, args ...any) {
	if !s.Logger.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	if s.HandlerOptions.AddSource {
		var pcs [1]uintptr
		runtime.Callers(s.skip, pcs[:])
		pc = pcs[0]
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = s.Logger.Handler().Handle(ctx, r)
}

func (s *Logger) TraceCtx(ctx context.Context, msg string, args ...any) {
	s.log(ctx, LevelTrace, msg, args...)
}

func (s *Logger) DebugCtx(ctx context.Context, msg string, args ...any) {
	s.log(ctx, LevelDebug, msg, args...)
}

func (s *Logger) InfoCtx(ctx context.Context, msg string, args ...any) {
	s.log(ctx, LevelInfo, msg, args...)
}

func (s *Logger) WarnCtx(ctx context.Context, msg string, args ...any) {
	s.log(ctx, LevelWarn, msg, args...)
}

func (s *Logger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	s.log(ctx, LevelError, msg, args...)
}

func (s *Logger) FatalCtx(ctx context.Context, msg string, args ...any) {
	s.log(ctx, LevelFatal, msg, args...)
}

func (s *Logger) Trace(msg string, args ...any) {
	s.log(nil, LevelTrace, msg, args...)
}

func (s *Logger) Debug(msg string, args ...any) {
	s.log(nil, LevelDebug, msg, args...)
}

func (s *Logger) Info(msg string, args ...any) {
	s.log(nil, LevelInfo, msg, args...)
}

func (s *Logger) Warn(msg string, args ...any) {
	s.log(nil, LevelWarn, msg, args...)
}

func (s *Logger) Error(msg string, args ...any) {
	s.log(nil, LevelError, msg, args...)
}

func (s *Logger) Fatal(msg string, args ...any) {
	s.log(nil, LevelFatal, msg, args...)
}

var (
	DefaultLogger = New(LevelAll)
)

func TraceCtx(ctx context.Context, msg string, args ...any) {
	DefaultLogger.log(ctx, LevelTrace, msg, args...)
}

func DebugCtx(ctx context.Context, msg string, args ...any) {
	DefaultLogger.log(ctx, LevelDebug, msg, args...)
}

func InfoCtx(ctx context.Context, msg string, args ...any) {
	DefaultLogger.log(ctx, LevelInfo, msg, args...)
}

func WarnCtx(ctx context.Context, msg string, args ...any) {
	DefaultLogger.log(ctx, LevelWarn, msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...any) {
	DefaultLogger.log(ctx, LevelError, msg, args...)
}

func FatalCtx(ctx context.Context, msg string, args ...any) {
	DefaultLogger.log(ctx, LevelFatal, msg, args...)
}

func Trace(msg string, args ...any) {
	DefaultLogger.log(nil, LevelTrace, msg, args...)
}

func Debug(msg string, args ...any) {
	DefaultLogger.log(nil, LevelDebug, msg, args...)
}

func Info(msg string, args ...any) {
	DefaultLogger.log(nil, LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	DefaultLogger.log(nil, LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	DefaultLogger.log(nil, LevelError, msg, args...)
}

func Fatal(msg string, args ...any) {
	DefaultLogger.log(nil, LevelFatal, msg, args...)
}
