package main

import (
	"errors"
	"fmt"
	"github.com/cd365/logger/v9"
	"github.com/rs/zerolog"
	"runtime"
	"strings"
	"time"
)

func loggerCallerString(frames *runtime.Frames) string {
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

func tryLogger() {
	l := logger.NewLogger(nil)

	fmt.Println(l.GetLevel().String())
	l.SetLevel(logger.DebugLevel)
	fmt.Println(l.GetLevel().String())

	l.WithContext(func(ctx zerolog.Context) zerolog.Logger {
		return ctx.
			Str("project", "project-a").
			Str("module", "module-a").
			Int64("program_start_at", time.Now().Unix()).
			Logger()
	})

	// Note the difference in the `unix_milli` values in the two log outputs
	l.AddEvent(func(event *logger.Event) {

		// // The value displays log information at warning level and above
		// if event.Level() <= logger.InfoLevel {
		// 	event.Discard()
		// 	return
		// }

		event.Int64("unix_milli", time.Now().UnixMilli())

		// Get stack information.
		callsString := loggerCallerString(logger.Callers(0))
		event.Str("stack", callsString)

	})

	l.Info().Msg("123")

	<-time.After(time.Second)

	l.Error().Str("tips", "compare the value of unix_milli").Err(errors.New("error value"))

	l.Trace().Msg("There will be no output here")
}
