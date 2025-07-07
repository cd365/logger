package logger

import (
	"fmt"
	"io"
	"os"
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

// Limiter Temporarily save the output log.
type Limiter interface {
	// Exists Check if a certain data exists.
	Exists(key string) (exists bool, err error)

	// Set Writing data.
	Set(key string, value []byte, duration ...time.Duration) error
}

// CallLimit Log output frequency limit (allows you to set the minimum log level for frequency limit).
type CallLimit struct {
	// duration Log frequency limit period.
	duration time.Duration

	// level The minimum level for limiting the frequency of log output. default: TraceLevel, that is, all log levels limit the frequency.
	level Level

	// limiter Current Limiter.
	limiter Limiter

	// writer The final output channel of the log.
	writer io.Writer

	// parser Parsing log data.
	parser Parser
}

func (s *CallLimit) GetLevel() Level {
	return s.level
}

func (s *CallLimit) SetLevel(level Level) *CallLimit {
	s.level = level
	return s
}

func (s *CallLimit) SetParser(parser Parser) *CallLimit {
	if parser != nil {
		s.parser = parser
	}
	return s
}

func (s *CallLimit) Write(content []byte) (int, error) {
	level, err := s.parser.GetLevel(content)
	if err != nil {
		return 0, err
	}
	if level < s.level {
		return s.writer.Write(content)
	}
	caller := s.parser.GetCaller(content)
	exists, err := s.limiter.Exists(caller)
	if err != nil {
		return 0, err
	}
	if exists {
		return len(content), nil
	}
	n, err := s.writer.Write(content)
	if err != nil {
		return 0, err
	}
	err = s.limiter.Set(
		caller,
		fmt.Appendf(nil, "%d", time.Now().UnixMilli()),
		s.duration,
	)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func NewCallLimit(duration time.Duration, limiter Limiter, writer io.Writer) *CallLimit {
	if limiter == nil {
		panic("logger: limiter is nil")
	}
	if writer == nil {
		writer = os.Stdout
	}
	return &CallLimit{
		duration: duration,
		level:    TraceLevel,
		limiter:  limiter,
		writer:   writer,
		parser:   &parser{},
	}
}

func colorByLevel(content []byte, level Level) []byte {
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
		return content
	}
	result := make([]byte, 0, len(ctrl)+len(content)+len(ctrlReset))
	result = append(result, ctrl...)
	result = append(result, content...)
	result = append(result, ctrlReset...)
	return result
}

// LevelColor Output log data in different colors according to the log level.
// It is often used to output to the terminal and is friendly to debugging and running programs.
type LevelColor struct {
	// level The lowest level of data output, usually set to WarnLevel, ErrorLevel, default: TraceLevel
	level Level

	// writer The final output channel of the log, usually the terminal.
	writer io.Writer

	// parser Parsing log data.
	parser Parser
}

func (s *LevelColor) GetLevel() Level {
	return s.level
}

func (s *LevelColor) SetLevel(level Level) *LevelColor {
	s.level = level
	return s
}

func (s *LevelColor) SetParser(parser Parser) *LevelColor {
	if parser != nil {
		s.parser = parser
	}
	return s
}

func (s *LevelColor) Write(content []byte) (int, error) {
	level, err := s.parser.GetLevel(content)
	if err != nil {
		return 0, err
	}
	if level < s.level {
		return len(content), nil
	}
	return s.writer.Write(colorByLevel(content, level))
}

func NewLevelColor(writer io.Writer) *LevelColor {
	if writer == nil {
		writer = os.Stdout
	}
	return &LevelColor{
		level:  TraceLevel,
		writer: writer,
		parser: &parser{},
	}
}
