package main

import (
	"context"
	"github.com/cd365/logger/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client struct {
	client *redis.Client
}

func (s *Client) Exists(key string) (bool, error) {
	val, err := s.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (s *Client) Set(key string, value []byte, duration ...time.Duration) error {
	dur := time.Duration(0)
	if len(duration) > 0 {
		dur = duration[0]
	}
	return s.client.Set(context.Background(), key, value, dur).Err()
}

func newRedisClient() *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Client{client: client}
}

// levelColor Log output color control.
func levelColor() {
	output := logger.NewLevelColor(nil)
	// output.SetLevel(logger.WarnLevel)

	l := logger.NewLogger(output)
	l.AddEvent(func(event *logger.Event) { event.Str("category", "color") })
	msg := "Hello World"
	l.Trace().Msg(msg)
	l.Debug().Msg(msg)
	l.Info().Msg(msg)
	l.Warn().Msg(msg)
	l.Error().Msg(msg)
}

// callLimit Log output limiting control.
func callLimit() {
	redisClient := newRedisClient()

	output := logger.NewCallLimit(time.Second*15, redisClient, nil)
	output.SetLevel(logger.WarnLevel)

	l := logger.NewLogger(output)
	l.AddEvent(func(event *logger.Event) { event.Str("category", "limit") })
	msg := "Hello World"
	l.Trace().Msg(msg)
	l.Debug().Msg(msg)
	l.Info().Msg(msg)
	l.Warn().Msg(msg)
	l.Error().Msg(msg)
}

func tryWrite() {
	levelColor()
	callLimit() // Please make sure the redis connection is ok first.
}
