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
	logLock = &sync.Mutex{}
)

type Logger struct {
	trace  *log.Logger
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	error  *log.Logger
	panic  *log.Logger
	closer []io.Closer
	level  Level
}

func New(writer io.Writer, flags int, level Level) *Logger {
	if flags <= 0 {
		flags = defaultLogFlags
	}
	s := &Logger{
		trace: log.New(writer, tagTrace, flags),
		debug: log.New(writer, tagDebug, flags),
		info:  log.New(writer, tagInfo, flags),
		warn:  log.New(writer, tagWarn, flags),
		error: log.New(writer, tagError, flags),
		panic: log.New(writer, tagPanic, flags),
		level: level,
	}
	if c, ok := writer.(io.Closer); ok {
		s.closer = append(s.closer, c)
	}
	return s
}

func (s *Logger) output(level Level, depth int, txt string) {
	if level < s.level {
		return
	}
	logLock.Lock()
	defer logLock.Unlock()
	switch level {
	case TRACE:
		s.trace.Output(3+depth, txt)
	case DEBUG:
		s.debug.Output(3+depth, txt)
	case INFO:
		s.info.Output(3+depth, txt)
	case WARN:
		s.warn.Output(3+depth, txt)
	case ERROR:
		s.error.Output(3+depth, txt)
	case PANIC:
		s.panic.Output(3+depth, txt)
	}
}

func (s *Logger) Close() {
	logLock.Lock()
	defer logLock.Unlock()
	for _, c := range s.closer {
		if err := c.Close(); err != nil {
			log.Printf("close log error: %v %s", c, err.Error())
		}
	}
}

func (s *Logger) Trace(v ...interface{}) {
	s.output(TRACE, 0, fmt.Sprint(v...))
}

func (s *Logger) TraceDepth(depth int, v ...interface{}) {
	s.output(TRACE, depth, fmt.Sprint(v...))
}

func (s *Logger) Traceln(v ...interface{}) {
	s.output(TRACE, 0, fmt.Sprintln(v...))
}

func (s *Logger) Tracef(format string, v ...interface{}) {
	s.output(TRACE, 0, fmt.Sprintf(format, v...))
}

func (s *Logger) Debug(v ...interface{}) {
	s.output(DEBUG, 0, fmt.Sprint(v...))
}

func (s *Logger) DebugDepth(depth int, v ...interface{}) {
	s.output(DEBUG, depth, fmt.Sprint(v...))
}

func (s *Logger) Debugln(v ...interface{}) {
	s.output(DEBUG, 0, fmt.Sprintln(v...))
}

func (s *Logger) Debugf(format string, v ...interface{}) {
	s.output(DEBUG, 0, fmt.Sprintf(format, v...))
}

func (s *Logger) Info(v ...interface{}) {
	s.output(INFO, 0, fmt.Sprint(v...))
}

func (s *Logger) InfoDepth(depth int, v ...interface{}) {
	s.output(INFO, depth, fmt.Sprint(v...))
}

func (s *Logger) Infoln(v ...interface{}) {
	s.output(INFO, 0, fmt.Sprintln(v...))
}

func (s *Logger) Infof(format string, v ...interface{}) {
	s.output(INFO, 0, fmt.Sprintf(format, v...))
}

func (s *Logger) Warn(v ...interface{}) {
	s.output(WARN, 0, fmt.Sprint(v...))
}

func (s *Logger) WarnDepth(depth int, v ...interface{}) {
	s.output(WARN, depth, fmt.Sprint(v...))
}

func (s *Logger) Warnln(v ...interface{}) {
	s.output(WARN, 0, fmt.Sprintln(v...))
}

func (s *Logger) Warnf(format string, v ...interface{}) {
	s.output(WARN, 0, fmt.Sprintf(format, v...))
}

func (s *Logger) Error(v ...interface{}) {
	s.output(ERROR, 0, fmt.Sprint(v...))
}

func (s *Logger) ErrorDepth(depth int, v ...interface{}) {
	s.output(ERROR, depth, fmt.Sprint(v...))
}

func (s *Logger) Errorln(v ...interface{}) {
	s.output(ERROR, 0, fmt.Sprintln(v...))
}

func (s *Logger) Errorf(format string, v ...interface{}) {
	s.output(ERROR, 0, fmt.Sprintf(format, v...))
}

func (s *Logger) Panic(v ...interface{}) {
	txt := fmt.Sprint(v...)
	s.output(PANIC, 0, txt)
	panic(txt)
}

func (s *Logger) PanicDepth(depth int, v ...interface{}) {
	txt := fmt.Sprint(v...)
	s.output(PANIC, depth, txt)
	panic(txt)
}

func (s *Logger) Panicln(v ...interface{}) {
	txt := fmt.Sprintln(v...)
	s.output(PANIC, 0, txt)
	panic(txt)
}

func (s *Logger) Panicf(format string, v ...interface{}) {
	txt := fmt.Sprintf(format, v...)
	s.output(PANIC, 0, txt)
	panic(txt)
}

func (s *Logger) SetLevel(level Level) *Logger {
	s.level = level
	return s
}

func (s *Logger) Level() Level {
	return s.level
}

func (s *Logger) SetFlags(flags int) *Logger {
	s.debug.SetFlags(flags)
	s.info.SetFlags(flags)
	s.warn.SetFlags(flags)
	s.error.SetFlags(flags)
	s.panic.SetFlags(flags)
	return s
}
