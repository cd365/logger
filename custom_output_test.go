package logger

import (
	"github.com/rs/zerolog"
	"testing"
)

func TestNewCustomOutput(t *testing.T) {
	output := NewCustomOutput(nil)
	output.AllowLevel(zerolog.ErrorLevel)
	l := NewLogger(output)
	msg := "Hello World"
	l.Trace().Msg(msg)
	l.Debug().Msg(msg)
	l.Info().Msg(msg)
	l.Warn().Msg(msg)
	l.Error().Msg(msg)
	l.Fatal().Msg(msg)
	l.Panic().Msg(msg)
}
