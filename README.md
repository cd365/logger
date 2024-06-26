# FOR EXAMPLE

```go
package main

import (
	"github.com/cd365/logger/v3"
	"log/slog"
	"os"
)

func main() {
	l := logger.New(logger.LevelAll)
	l.LevelVar.Set(logger.LevelWarn)
	msg := "test"
	args := []any{"title", "title-value"}
	l.Trace(msg, args...)
	l.Debug(msg, args...)
	l.Info(msg, args...)
	l.Warn(msg, args...)
	l.Error(msg, args...)
	l.Fatal(msg, args...)
	l.Logger = slog.New(slog.NewJSONHandler(os.Stdout, l.HandlerOptions))
	l.LevelVar.Set(logger.LevelInfo)
	l.Trace(msg, args...)
	l.Debug(msg, args...)
	l.Info(msg, args...)
	l.Warn(msg, args...)
	l.Error(msg, args...)
	l.Fatal(msg, args...)
	logger.Trace(msg, args...)
	logger.Debug(msg, args...)
	logger.Info(msg, args...)
	logger.Warn(msg, args...)
	logger.Error(msg, args...)
	logger.Fatal(msg, args...)
}
```
