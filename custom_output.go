package logger

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
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

type parseLevel struct {
	Level string `json:"level"`
}

type CustomOutput struct {
	allowLevel func(level string) bool
	parseLevel func(b []byte) (string, error)
	custom     func(b []byte) (int, error)
	writer     io.Writer
}

func (s *CustomOutput) AllowLevel(level zerolog.Level) *CustomOutput {
	s.allowLevel = func(tmp string) bool {
		if val, err := zerolog.ParseLevel(tmp); err != nil {
			return false
		} else {
			if val < level {
				return false
			}
		}
		return true
	}
	return s
}

func (s *CustomOutput) ParseLevel(parseLevel func(b []byte) (string, error)) *CustomOutput {
	if parseLevel != nil {
		s.parseLevel = parseLevel
	}
	return s
}

func (s *CustomOutput) Custom(custom func(b []byte) (int, error)) *CustomOutput {
	if custom != nil {
		s.custom = custom
	}
	return s
}

func (s *CustomOutput) write(b []byte) (int, error) {
	level, err := s.parseLevel(b)
	if err != nil {
		return 0, err
	}
	if !s.allowLevel(level) {
		return len(b), nil
	}
	ctrl := ctrlWhite
	switch level {
	case zerolog.LevelTraceValue:
		ctrl = ctrlMagenta
	case zerolog.LevelDebugValue:
		ctrl = ctrlCyan
	case zerolog.LevelInfoValue:
		ctrl = ctrlGreen
	case zerolog.LevelWarnValue:
		ctrl = ctrlYellow
	case zerolog.LevelErrorValue:
		ctrl = ctrlRed
	case zerolog.LevelFatalValue, zerolog.LevelPanicValue:
		ctrl = fmt.Sprintf("%s%s%s", ctrlBlue, ctrlBold, ctrlUnderline)
	default:
		return s.writer.Write(b)
	}
	write := make([]byte, 0, len(ctrl)+len(b)+len(ctrlReset))
	write = append(write, ctrl...)
	write = append(write, b...)
	write = append(write, ctrlReset...)
	return s.writer.Write(write)
}

func (s *CustomOutput) Write(b []byte) (int, error) {
	return s.custom(b)
}

func NewCustomOutput(writer io.Writer) *CustomOutput {
	if writer == nil {
		writer = os.Stdout
	}
	result := &CustomOutput{
		writer: writer,
	}
	result.ParseLevel(func(b []byte) (string, error) {
		tmp := &parseLevel{}
		if err := json.Unmarshal(b, tmp); err != nil {
			return "", err
		}
		return tmp.Level, nil
	})
	result.AllowLevel(zerolog.WarnLevel)
	result.custom = result.write
	return result
}
