package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"runtime"
	"time"
)

type Level int8

const (
	TraceLevel = Level(zerolog.TraceLevel)
	DebugLevel = Level(zerolog.DebugLevel)
	InfoLevel  = Level(zerolog.InfoLevel)
	WarnLevel  = Level(zerolog.WarnLevel)
	ErrorLevel = Level(zerolog.ErrorLevel)
	FatalLevel = Level(zerolog.FatalLevel)
	PanicLevel = Level(zerolog.PanicLevel)
	NoLevel    = Level(zerolog.NoLevel)
	Disabled   = Level(zerolog.Disabled)
)

func (s Level) String() string {
	return zerolog.Level(s).String()
}

type Event struct {
	event *zerolog.Event
	level Level
}

func newEvent(event *zerolog.Event, level Level) *Event {
	return &Event{
		event: event,
		level: level,
	}
}

func (s *Event) Level() Level {
	return s.level
}

func (s *Event) CallerSkipFrame(skip int) *Event {
	s.event.CallerSkipFrame(skip)
	return s
}

func (s *Event) Discard() *Event {
	s.event.Discard()
	return s
}

func (s *Event) Func(f func(e *Event)) *Event {
	if f != nil {
		f(s)
	}
	return s
}

func (s *Event) Any(key string, v any) *Event {
	s.event.Any(key, v)
	return s
}

func (s *Event) Str(key, v string) *Event {
	s.event.Str(key, v)
	return s
}

func (s *Event) Bytes(key string, v []byte) *Event {
	s.event.Bytes(key, v)
	return s
}

func (s *Event) Bool(key string, v bool) *Event {
	s.event.Bool(key, v)
	return s
}

func (s *Event) Int(key string, v int) *Event {
	s.event.Int(key, v)
	return s
}

func (s *Event) Int8(key string, v int8) *Event {
	s.event.Int8(key, v)
	return s
}

func (s *Event) Int16(key string, v int16) *Event {
	s.event.Int16(key, v)
	return s
}

func (s *Event) Int32(key string, v int32) *Event {
	s.event.Int32(key, v)
	return s
}

func (s *Event) Int64(key string, v int64) *Event {
	s.event.Int64(key, v)
	return s
}

func (s *Event) Uint(key string, v uint) *Event {
	s.event.Uint(key, v)
	return s
}

func (s *Event) Uint8(key string, v uint8) *Event {
	s.event.Uint8(key, v)
	return s
}

func (s *Event) Uint16(key string, v uint16) *Event {
	s.event.Uint16(key, v)
	return s
}

func (s *Event) Uint32(key string, v uint32) *Event {
	s.event.Uint32(key, v)
	return s
}

func (s *Event) Uint64(key string, v uint64) *Event {
	s.event.Uint64(key, v)
	return s
}

func (s *Event) Float32(key string, v float32) *Event {
	s.event.Float32(key, v)
	return s
}

func (s *Event) Float64(key string, v float64) *Event {
	s.event.Float64(key, v)
	return s
}

func (s *Event) Ints(key string, v []int) *Event {
	s.event.Ints(key, v)
	return s
}

func (s *Event) Ints8(key string, v []int8) *Event {
	s.event.Ints8(key, v)
	return s
}

func (s *Event) Ints16(key string, v []int16) *Event {
	s.event.Ints16(key, v)
	return s
}

func (s *Event) Ints32(key string, v []int32) *Event {
	s.event.Ints32(key, v)
	return s
}

func (s *Event) Ints64(key string, v []int64) *Event {
	s.event.Ints64(key, v)
	return s
}

func (s *Event) Uints(key string, v []uint) *Event {
	s.event.Uints(key, v)
	return s
}

func (s *Event) Uints8(key string, v []uint8) *Event {
	s.event.Uints8(key, v)
	return s
}

func (s *Event) Uints16(key string, v []uint16) *Event {
	s.event.Uints16(key, v)
	return s
}

func (s *Event) Uints32(key string, v []uint32) *Event {
	s.event.Uints32(key, v)
	return s
}

func (s *Event) Uints64(key string, v []uint64) *Event {
	s.event.Uints64(key, v)
	return s
}

func (s *Event) Floats32(key string, v []float32) *Event {
	s.event.Floats32(key, v)
	return s
}

func (s *Event) Floats64(key string, v []float64) *Event {
	s.event.Floats64(key, v)
	return s
}

func (s *Event) Time(key string, v time.Time) *Event {
	s.event.Time(key, v)
	return s
}

func (s *Event) Times(key string, v []time.Time) *Event {
	s.event.Times(key, v)
	return s
}

func (s *Event) Duration(key string, v time.Duration) *Event {
	s.event.Dur(key, v)
	return s
}

func (s *Event) Durations(key string, v []time.Duration) *Event {
	s.event.Durs(key, v)
	return s
}

func (s *Event) Title(title string) *Event {
	return s.Str("title", title)
}

func (s *Event) Data(data any) *Event {
	return s.Any("data", data)
}

func (s *Event) Err(err error) {
	if err != nil {
		s.event.Msg(err.Error())
	}
}

func (s *Event) Msg(msg string) {
	s.event.Msg(msg)
}

type Logger struct {
	logger *zerolog.Logger

	events []func(event *Event)
}

func NewLogger(writer io.Writer) *Logger {
	if writer == nil {
		writer = os.Stdout
	}
	logger := zerolog.New(writer).With().Caller().Timestamp().Logger()
	logger.Level(zerolog.TraceLevel)
	result := &Logger{
		logger: &logger,
	}
	return result.AddEvent(func(event *Event) { event.CallerSkipFrame(1) })
}

func (s *Logger) GetLevel() Level {
	return Level(s.logger.GetLevel())
}

func (s *Logger) SetLevel(level Level) *Logger {
	logger := s.logger.Level(zerolog.Level(level))
	s.logger = &logger
	return s
}

// WithContext Set common log properties.
func (s *Logger) WithContext(custom func(ctx zerolog.Context) zerolog.Logger) *Logger {
	if custom != nil {
		ctx := s.logger.With()
		logger := custom(ctx)
		s.logger = &logger
	}
	return s
}

// AddEvent Set custom properties before calling output log.
func (s *Logger) AddEvent(event func(event *Event)) *Logger {
	if event != nil {
		s.events = append(s.events, event)
	}
	return s
}

func (s *Logger) GetLogger() *zerolog.Logger {
	return s.logger
}

func (s *Logger) SetLogger(logger *zerolog.Logger) *Logger {
	s.logger = logger
	return s
}

// SetOutput Duplicates the current logger and sets writer as its output.
func (s *Logger) SetOutput(writer io.Writer) *Logger {
	logger := s.logger.Output(writer)
	s.logger = &logger
	return s
}

func (s *Logger) newEvent(level Level) *Event {
	switch level {
	case DebugLevel:
		return newEvent(s.logger.Debug(), level)
	case InfoLevel:
		return newEvent(s.logger.Info(), level)
	case WarnLevel:
		return newEvent(s.logger.Warn(), level)
	case ErrorLevel:
		return newEvent(s.logger.Error(), level)
	case FatalLevel:
		return newEvent(s.logger.Fatal(), level)
	case PanicLevel:
		return newEvent(s.logger.Panic(), level)
	default:
		return newEvent(s.logger.Trace(), level)
	}
}

// Trace Logger level TRACE.
func (s *Logger) Trace() *Event {
	level := TraceLevel
	event := s.newEvent(level)
	for _, f := range s.events {
		f(event)
	}
	return event
}

// Debug Logger level DEBUG.
func (s *Logger) Debug() *Event {
	level := DebugLevel
	event := s.newEvent(level)
	for _, f := range s.events {
		f(event)
	}
	return event
}

// Info Logger level INFO.
func (s *Logger) Info() *Event {
	level := InfoLevel
	event := s.newEvent(level)
	for _, f := range s.events {
		f(event)
	}
	return event
}

// Warn Logger level WARN.
func (s *Logger) Warn() *Event {
	level := WarnLevel
	event := s.newEvent(level)
	for _, f := range s.events {
		f(event)
	}
	return event
}

// Error Logger level ERROR.
func (s *Logger) Error() *Event {
	level := ErrorLevel
	event := s.newEvent(level)
	for _, f := range s.events {
		f(event)
	}
	return event
}

// Fatal Logger level FATAL.
func (s *Logger) Fatal() *Event {
	level := FatalLevel
	event := s.newEvent(level)
	for _, f := range s.events {
		f(event)
	}
	return event
}

// Panic Logger level PANIC.
func (s *Logger) Panic() *Event {
	level := PanicLevel
	event := s.newEvent(level)
	for _, f := range s.events {
		f(event)
	}
	return event
}

// defaultLogger Default logger.
var defaultLogger = NewLogger(nil)

func Default() *Logger {
	return defaultLogger
}

func Trace() *Event {
	return defaultLogger.Trace()
}

func Debug() *Event {
	return defaultLogger.Debug()
}

func Info() *Event {
	return defaultLogger.Info()
}

func Warn() *Event {
	return defaultLogger.Warn()
}

func Error() *Event {
	return defaultLogger.Error()
}

func Fatal() *Event {
	return defaultLogger.Fatal()
}

func Panic() *Event {
	return defaultLogger.Panic()
}

// Callers Get called lists.
func Callers(skip int) *runtime.Frames {
	pc := make([]uintptr, 1<<5)
	for {
		n := runtime.Callers(skip, pc)
		if n < len(pc) {
			pc = pc[:n]
			break
		}
		pc = make([]uintptr, 2*len(pc))
	}
	return runtime.CallersFrames(pc)
}

// ParseLevel Parse logger level.
func ParseLevel(level string) (Level, error) {
	switch level {
	case "TRACE", "trace":
		return TraceLevel, nil
	case "DEBUG", "debug":
		return DebugLevel, nil
	case "INFO", "info":
		return InfoLevel, nil
	case "WARN", "warn":
		return WarnLevel, nil
	case "ERROR", "error":
		return ErrorLevel, nil
	case "FATAL", "fatal":
		return FatalLevel, nil
	case "PANIC", "panic":
		return PanicLevel, nil
	default:
		return NoLevel, fmt.Errorf("invalid level: %s", level)
	}
}
