package logger

import (
	"errors"
	"github.com/rs/zerolog"
	"testing"
)

type services struct {
	MySQL func(event *zerolog.Event)
	Redis func(event *zerolog.Event)
	// ...
}

func newService(service string) func(event *zerolog.Event) {
	return func(event *zerolog.Event) { event.Str("service", service) }
}

func newServices() *services {
	return &services{
		MySQL: newService("MySQL"),
		Redis: newService("Redis"),
		// ...
	}
}

func TestNewRecords(t *testing.T) {
	srv := newServices()

	lg := NewLogger(nil)

	rc := NewRecords(lg)

	rc.Func(func(event *zerolog.Event) {
		// Add custom default data information.
		event.CallerSkipFrame(1)
		event.Str("project", "project_1")
		event.Str("module", "module_a")
		event.Str("version", "v1.0.0")
	})
	// rc.SetTagName("___tag").SetDataName("___data") // Try uncommenting it and watching the log output

	err := errors.New("hello error")
	rc.Debug().Msg("")
	rc.Info().Msg("hello")
	rc.Warn().Set("a", 1).Set("b", 10.08).Msg("hello")
	rc.Error().Func(srv.Redis).Err(err)
	rc.Error().Func(srv.MySQL).Tag("database query").Err(err)
	rc.Fatal().Err(err) // Try commenting it out
	rc.Panic().Msg("hello")
}
