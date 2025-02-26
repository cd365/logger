package logger

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"runtime"
)

type Logger struct {
	logger *zerolog.Logger

	customEvent []func(event *zerolog.Event, level zerolog.Level) *zerolog.Event
}

func NewLogger(writer io.Writer) *Logger {
	if writer == nil {
		writer = os.Stdout
	}
	logger := zerolog.New(writer).With().Caller().Timestamp().Logger()
	logger.Level(zerolog.TraceLevel)
	return &Logger{
		logger: &logger,
	}
}

func (s *Logger) GetLevel() zerolog.Level {
	return s.logger.GetLevel()
}

func (s *Logger) SetLevel(lvl zerolog.Level) *Logger {
	logger := s.logger.Level(lvl)
	s.logger = &logger
	return s
}

// CustomContext Set common log properties.
func (s *Logger) CustomContext(custom func(ctx zerolog.Context) zerolog.Logger) *Logger {
	if custom != nil {
		ctx := s.logger.With()
		logger := custom(ctx)
		s.logger = &logger
	}
	return s
}

// CustomEvent Set custom properties before calling output log.
func (s *Logger) CustomEvent(customEvent func(event *zerolog.Event, level zerolog.Level) *zerolog.Event) *Logger {
	if customEvent != nil {
		s.customEvent = append(s.customEvent, customEvent)
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

func (s *Logger) Trace() *zerolog.Event {
	tmp := s.logger.Trace()
	for _, fc := range s.customEvent {
		fc(tmp, zerolog.TraceLevel)
	}
	return tmp
}

func (s *Logger) Debug() *zerolog.Event {
	tmp := s.logger.Debug()
	for _, fc := range s.customEvent {
		fc(tmp, zerolog.DebugLevel)
	}
	return tmp
}

func (s *Logger) Info() *zerolog.Event {
	tmp := s.logger.Info()
	for _, fc := range s.customEvent {
		fc(tmp, zerolog.InfoLevel)
	}
	return tmp
}

func (s *Logger) Warn() *zerolog.Event {
	tmp := s.logger.Warn()
	for _, fc := range s.customEvent {
		fc(tmp, zerolog.WarnLevel)
	}
	return tmp
}

func (s *Logger) Error() *zerolog.Event {
	tmp := s.logger.Error()
	for _, fc := range s.customEvent {
		fc(tmp, zerolog.ErrorLevel)
	}
	return tmp
}

func (s *Logger) Fatal() *zerolog.Event {
	tmp := s.logger.Fatal()
	for _, fc := range s.customEvent {
		fc(tmp, zerolog.FatalLevel)
	}
	return tmp
}

func (s *Logger) Panic() *zerolog.Event {
	tmp := s.logger.Panic()
	for _, fc := range s.customEvent {
		fc(tmp, zerolog.PanicLevel)
	}
	return tmp
}

var defaultLogger = NewLogger(nil)

func Default() *Logger {
	return defaultLogger
}

func Trace() *zerolog.Event {
	return defaultLogger.Trace()
}

func Debug() *zerolog.Event {
	return defaultLogger.Debug()
}

func Info() *zerolog.Event {
	return defaultLogger.Info()
}

func Warn() *zerolog.Event {
	return defaultLogger.Warn()
}

func Error() *zerolog.Event {
	return defaultLogger.Error()
}

func Fatal() *zerolog.Event {
	return defaultLogger.Fatal()
}

func Panic() *zerolog.Event {
	return defaultLogger.Panic()
}

type Called struct {
	Index int // call index

	Func string // call func

	File string // call file

	FileLine int // call line in the file
}

// Callers Get called lists.
func Callers(skip int) []*Called {
	pcs := make([]uintptr, 1<<8)
	n := runtime.Callers(skip, pcs)
	pcs = pcs[:n]
	result := make([]*Called, 0, n)
	for i, pc := range pcs {
		fc := runtime.FuncForPC(pc)
		if fc == nil {
			continue
		}
		file, line := fc.FileLine(pc)
		caller := &Called{
			Index:    i,
			Func:     fc.Name(),
			File:     file,
			FileLine: line,
		}
		result = append(result, caller)
	}
	return result
}

// CalledLists Called lists.
func CalledLists(calledLists []*Called) []byte {
	buffer := bytes.NewBuffer(nil)
	for _, called := range calledLists {
		if called != nil {
			buffer.WriteString(fmt.Sprintf("%d %s %s:%d\n", called.Index, called.Func, called.File, called.FileLine))
		}
	}
	return buffer.Bytes()
}
