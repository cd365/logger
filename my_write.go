package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	ctrlReset = "\x1b[0m"

	ctrlBold      = "\x1b[1m"
	ctrlDim       = "\x1b[2m"
	ctrlItalic    = "\x1b[3m"
	ctrlUnderline = "\x1b[4m"

	ctrlBlack   = "\x1b[30m"
	ctrlRed     = "\x1b[31m"
	ctrlGreen   = "\x1b[32m"
	ctrlYellow  = "\x1b[33m"
	ctrlBlue    = "\x1b[34m"
	ctrlMagenta = "\x1b[35m"
	ctrlCyan    = "\x1b[36m"
	ctrlWhite   = "\x1b[37m"
)

type levelParse struct {
	Level string `json:"level"`
}

// MyWrite Control log output according to custom conditions.
// For example: Control log output by log level and log output source.
type MyWrite struct {
	// allowLevel Minimum log output level.
	allowLevel Level

	// ParseLevel Method for parsing log levels (customizable).
	ParseLevel func(b []byte) (Level, error)

	// priorityWrite Customize the log priority output method (replace the Write method of the log output channel).
	priorityWrite func(b []byte) (int, error)

	// writer Final output channel for logs.
	writer io.Writer
}

// GetAllowLevel Get the minimum log output level.
func (s *MyWrite) GetAllowLevel() Level {
	return s.allowLevel
}

// SetAllowLevel Set the minimum log output level.
func (s *MyWrite) SetAllowLevel(level Level) *MyWrite {
	s.allowLevel = level
	return s
}

// parseLevel Default parsing log level method.
func (s *MyWrite) parseLevel(b []byte) (Level, error) {
	tmp := &levelParse{}
	if err := json.Unmarshal(b, tmp); err != nil {
		return NoLevel, err
	}
	switch tmp.Level {
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
		return NoLevel, nil
	}
}

// GetWriter Get the log output channel.
func (s *MyWrite) GetWriter() io.Writer {
	return s.writer
}

// SetWriter Set the log output channel.
func (s *MyWrite) SetWriter(writer io.Writer) *MyWrite {
	if writer != nil {
		s.writer = writer
	}
	return s
}

// ColorWrite Set different colors for different levels of log output, but Fatal and Panic are both blue.
func (s *MyWrite) ColorWrite(b []byte) (int, error) {
	level, err := s.ParseLevel(b)
	if err != nil {
		return 0, err
	}
	if s.allowLevel <= PanicLevel && level < s.allowLevel {
		return len(b), nil
	}
	ctrl := ctrlWhite
	switch level {
	case TraceLevel:
		ctrl = ctrlMagenta
	case DebugLevel:
		ctrl = ctrlCyan
	case InfoLevel:
		ctrl = ctrlGreen
	case WarnLevel:
		ctrl = ctrlYellow
	case ErrorLevel:
		ctrl = ctrlRed
	case FatalLevel, PanicLevel:
		ctrl = fmt.Sprintf("%s%s%s", ctrlBlue, ctrlBold, ctrlUnderline)
	default:
		return s.GetWriter().Write(b)
	}
	write := make([]byte, 0, len(ctrl)+len(b)+len(ctrlReset))
	write = append(write, ctrl...)
	write = append(write, b...)
	write = append(write, ctrlReset...)
	return s.GetWriter().Write(write)
}

// LimitWrite Limit the output of logs to Warn and above by default, usually use cache for current limiting control.
// Of course, you can also completely customize the current limiting logic.
func (s *MyWrite) LimitWrite(
	getCaller func(b []byte) (string, error),
	getExists func(key string) (exists bool, err error),
	setDuration time.Duration,
	setValue func(key string, value string, duration time.Duration) error,
) func([]byte) (int, error) {
	return func(b []byte) (int, error) {
		level, err := s.ParseLevel(b)
		if err != nil {
			return 0, err
		}
		if s.allowLevel <= PanicLevel && level < s.allowLevel {
			return len(b), nil
		}
		if level < WarnLevel || level > PanicLevel {
			return s.GetWriter().Write(b)
		}
		caller, err := getCaller(b)
		if err != nil {
			return 0, err
		}
		if caller == "" {
			return s.GetWriter().Write(b)
		}
		has, err := getExists(caller)
		if err != nil {
			return 0, err
		}
		if has {
			// limited
			return len(b), nil
		}
		n, err := s.GetWriter().Write(b)
		if err != nil {
			return 0, err
		}
		if err = setValue(caller, strconv.FormatInt(time.Now().UnixMilli(), 10), setDuration); err != nil {
			return 0, err
		}
		return n, nil
	}
}

// PriorityWrite Customize the log priority output method.
func (s *MyWrite) PriorityWrite(write func(b []byte) (int, error)) *MyWrite {
	if write != nil {
		s.priorityWrite = write
	}
	return s
}

// Write Rewrite the log output method, the priorityWrite method takes precedence.
func (s *MyWrite) Write(b []byte) (int, error) {
	if s.priorityWrite != nil {
		return s.priorityWrite(b)
	}
	return s.GetWriter().Write(b)
}

func NewMyWrite(writer io.Writer) *MyWrite {
	if writer == nil {
		writer = os.Stdout
	}
	newMyWrite := &MyWrite{
		allowLevel: TraceLevel,
		writer:     writer,
	}
	newMyWrite.ParseLevel = newMyWrite.parseLevel
	return newMyWrite
}
