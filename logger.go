package logger

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type Level int

const (
	ALL Level = iota
	TRACE
	DEBUG
	INFO
	WARN
	ERROR
	PANIC
	OFF
)

const (
	tagTrace = "TRACE: "
	tagDebug = "DEBUG: "
	tagInfo  = "INFO : "
	tagWarn  = "WARN : "
	tagError = "ERROR: "
	tagPanic = "PANIC: "
)

const (
	defaultLogFlags = log.Ldate | log.Lmicroseconds | log.Lshortfile
)

var (
	mutex = &sync.Mutex{}
)

type Logger struct {
	trace     *log.Logger
	debug     *log.Logger
	info      *log.Logger
	warn      *log.Logger
	error     *log.Logger
	panic     *log.Logger
	closer    []io.Closer
	Level     Level
	CallDepth int
}

func New(writer io.Writer, flags int, level Level) *Logger {
	if flags <= 0 {
		flags = defaultLogFlags
	}
	s := &Logger{
		trace:     log.New(writer, tagTrace, flags),
		debug:     log.New(writer, tagDebug, flags),
		info:      log.New(writer, tagInfo, flags),
		warn:      log.New(writer, tagWarn, flags),
		error:     log.New(writer, tagError, flags),
		panic:     log.New(writer, tagPanic, flags),
		Level:     level,
		CallDepth: 4,
	}
	if c, ok := writer.(io.Closer); ok {
		s.closer = append(s.closer, c)
	}
	return s
}

func (s *Logger) Close() {
	mutex.Lock()
	defer mutex.Unlock()
	for _, c := range s.closer {
		_ = c.Close()
	}
}

func (s *Logger) Flags(flags int) *Logger {
	if flags <= 0 {
		flags = defaultLogFlags
	}
	s.debug.SetFlags(flags)
	s.info.SetFlags(flags)
	s.warn.SetFlags(flags)
	s.error.SetFlags(flags)
	s.panic.SetFlags(flags)
	return s
}

func (s *Logger) output(level Level, text string) {
	mutex.Lock()
	defer mutex.Unlock()
	switch level {
	case TRACE:
		_ = s.trace.Output(s.CallDepth, text)
	case DEBUG:
		_ = s.debug.Output(s.CallDepth, text)
	case INFO:
		_ = s.info.Output(s.CallDepth, text)
	case WARN:
		_ = s.warn.Output(s.CallDepth, text)
	case ERROR:
		_ = s.error.Output(s.CallDepth, text)
	case PANIC:
		_ = s.panic.Output(s.CallDepth, text)
		panic(text)
	}
}

func (s *Logger) sprint(level Level, v ...interface{}) {
	if level >= s.Level {
		s.output(level, fmt.Sprint(v...))
	}
}

func (s *Logger) sprintf(level Level, format string, v ...interface{}) {
	if level >= s.Level {
		s.output(level, fmt.Sprintf(format, v...))
	}
}

func (s *Logger) Trace(v ...interface{}) {
	s.sprint(TRACE, v...)
}

func (s *Logger) Tracef(format string, v ...interface{}) {
	s.sprintf(TRACE, format, v...)
}

func (s *Logger) Debug(v ...interface{}) {
	s.sprint(DEBUG, v...)
}

func (s *Logger) Debugf(format string, v ...interface{}) {
	s.sprintf(DEBUG, format, v...)
}

func (s *Logger) Info(v ...interface{}) {
	s.sprint(INFO, v...)
}

func (s *Logger) Infof(format string, v ...interface{}) {
	s.sprintf(INFO, format, v...)
}

func (s *Logger) Warn(v ...interface{}) {
	s.sprint(WARN, v...)
}

func (s *Logger) Warnf(format string, v ...interface{}) {
	s.sprintf(WARN, format, v...)
}

func (s *Logger) Error(v ...interface{}) {
	s.sprint(ERROR, v...)
}

func (s *Logger) Errorf(format string, v ...interface{}) {
	s.sprintf(ERROR, format, v...)
}

func (s *Logger) Panic(v ...interface{}) {
	s.sprint(PANIC, v...)
}

func (s *Logger) Panicf(format string, v ...interface{}) {
	s.sprintf(PANIC, format, v...)
}
