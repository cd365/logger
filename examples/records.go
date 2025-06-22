package main

import (
	"errors"
	"github.com/cd365/logger/v9"
)

type services struct {
	MySQL func(event *logger.Event)
	Redis func(event *logger.Event)
	// ...
}

func newService(service string) func(event *logger.Event) {
	return func(event *logger.Event) { event.Str("service", service) }
}

func newServices() *services {
	return &services{
		MySQL: newService("MySQL"),
		Redis: newService("Redis"),
		// ...
	}
}

func tryRecords() {
	srv := newServices()

	lg := logger.NewLogger(nil)

	lg.AddEvent(func(event *logger.Event) {
		// Add custom default data information.
		event.Str("project", "project_1")
		event.Str("module", "module_a")
		event.Str("version", "v1.0.0")
		if event.Level() >= logger.ErrorLevel {
			// TODO ...
			// Add call stack?
		}
	})

	rc := logger.NewRecords(lg)

	// rc.SetTagName("___tag").SetDataName("___data") // Try uncommenting it and watching the log output

	err := errors.New("hello error")
	rc.Debug().Msg("")
	rc.Info().Msg("hello")
	rc.Warn().Set("a", 1).Set("b", 10.08).Msg("hello")
	rc.Error(srv.Redis).Err(err)
	rc.Error(srv.MySQL).Tag("database query").Err(err)
	// rc.Fatal().Err(err) // Try uncommenting it
	// rc.Panic().Msg("hello")
}
