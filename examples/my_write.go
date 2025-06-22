package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cd365/logger/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type parseCaller struct {
	Caller string `json:"caller"`
}

// myWriteColor Log output color control.
func myWriteColor() {
	output := logger.NewMyWrite(nil)
	// output.SetAllowLevel(logger.WarnLevel)

	output.PriorityWrite(output.ColorWrite)
	// output.SetWriter(os.Stdout)

	l := logger.NewLogger(output)
	l.AddEvent(func(event *logger.Event) { event.Str("custom_write", "color") })
	msg := "Hello World"
	l.Trace().Msg(msg)
	l.Debug().Msg(msg)
	l.Info().Msg(msg)
	l.Warn().Msg(msg)
	l.Error().Msg(msg)
}

// myWriteLimit Log output limiting control.
func myWriteLimit() {
	output := logger.NewMyWrite(nil)
	// output.SetAllowLevel(logger.ErrorLevel)

	redisClient := newRedisClient()
	output.PriorityWrite(output.LimitWrite(
		func(b []byte) (string, error) {
			tmp := &parseCaller{}
			if err := json.Unmarshal(b, tmp); err != nil {
				return "", err
			}
			return tmp.Caller, nil
		},
		func(key string) (exists bool, err error) {
			cmd := redisClient.Get(context.Background(), key)
			if cmd.Err() != nil {
				if errors.Is(cmd.Err(), redis.Nil) {
					return false, nil
				}
				return false, err
			}
			return true, nil
		},
		time.Minute,
		func(key string, value string, duration time.Duration) error {
			return redisClient.Set(context.Background(), key, value, duration).Err()
		},
	))
	// output.SetWriter(os.Stdout)

	l := logger.NewLogger(output)
	l.AddEvent(func(event *logger.Event) { event.Str("custom_write", "limit") })
	msg := "Hello World"
	l.Trace().Msg(msg)
	l.Debug().Msg(msg)
	l.Info().Msg(msg)
	l.Warn().Msg(msg)
	l.Error().Msg(msg)
}

func tryMyWrite() {
	myWriteColor()
	// myWriteLimit() // Please make sure the redis connection is ok first.
}
