package logger

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"runtime"
	"strings"
	"testing"
	"time"
)

func callerString(frames *runtime.Frames) string {
	b := &strings.Builder{}
	index := 0
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		if index > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("%d %s %s:%d", index, frame.Function, frame.File, frame.Line))
		index++
	}
	return b.String()
}

func TestNewLogger(t *testing.T) {
	l := NewLogger(nil)

	fmt.Println(l.GetLevel().String())
	l.SetLevel(zerolog.DebugLevel)
	fmt.Println(l.GetLevel().String())

	l.CustomContext(func(ctx zerolog.Context) zerolog.Logger {
		return ctx.
			Str("project", "project-a").
			Str("module", "module-a").
			Int64("program_start_at", time.Now().Unix()).
			Logger()
	})

	// Note the difference in the `unix_milli` values in the two log outputs
	l.CustomEvent(func(event *zerolog.Event, level zerolog.Level) *zerolog.Event {
		return event.Int64("unix_milli", time.Now().UnixMilli())
	})
	// callers
	l.CustomEvent(func(event *zerolog.Event, level zerolog.Level) *zerolog.Event {
		if level > zerolog.InfoLevel {
			callsString := callerString(Callers(1))
			fmt.Println(callsString)
			event.Str("callers", callsString)
		}
		return event
	})

	l.Info().Msg("123")

	<-time.After(time.Second * 2)

	l.Error().Err(errors.New("321")).Send()

	l.Trace().Msg("000")

}
