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

var levelMap = map[Level]string{
	LevelAll:   "ALL",
	LevelTrace: "TRACE",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelOff:   "OFF",
}

type Logger struct {
	// options Logger handler options
	options *slog.HandlerOptions

	// logger Log object
	logger *slog.Logger

	// levelVar Dynamically adjust log level
	levelVar *slog.LevelVar
}

func New(level Level) *Logger {
	levelVar := slog.LevelVar{}
	levelVar.Set(level)
	options := &slog.HandlerOptions{
		// With source
		AddSource: true,
		// Support dynamic setting of log level
		Level: &levelVar,
		// Modify the Attr key-value pair in the log (that is, the key/value attached to the log record)
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				levelValue := a.Value.Any().(slog.Level)
				levelLabel := levelValue.String()
				switch levelValue {
				case LevelTrace:
					levelLabel = levelMap[levelValue]
				case LevelFatal:
					levelLabel = levelMap[levelValue]
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, options))
	return &Logger{
		options:  options,
		logger:   logger,
		levelVar: &levelVar,
	}
}

func (s *Logger) GetOptions() *slog.HandlerOptions {
	return s.options
}

func (s *Logger) SetHandler(handler slog.Handler) *Logger {
	s.logger = slog.New(handler)
	return s
}

func (s *Logger) SetLevel(level Level) *Logger {
	s.levelVar.Set(level)
	return s
}

func (s *Logger) log(ctx context.Context, level Level, msg string, args ...any) {
	if !s.logger.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	if s.options.AddSource {
		var pcs [1]uintptr
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = s.logger.Handler().Handle(ctx, r)
}

func (s *Logger) Trace(msg string, args ...any) {
	s.log(context.Background(), LevelTrace, msg, args...)
}

func (s *Logger) Debug(msg string, args ...any) {
	s.log(context.Background(), LevelDebug, msg, args...)
}

func (s *Logger) Info(msg string, args ...any) {
	s.log(context.Background(), LevelInfo, msg, args...)
}

func (s *Logger) Warn(msg string, args ...any) {
	s.log(context.Background(), LevelWarn, msg, args...)
}

func (s *Logger) Error(msg string, args ...any) {
	s.log(context.Background(), LevelError, msg, args...)
}

func (s *Logger) Fatal(msg string, args ...any) {
	s.log(context.Background(), LevelFatal, msg, args...)
}
