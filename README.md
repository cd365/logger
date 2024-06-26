# FOR EXAMPLE

```go
package main

import (
	"github.com/cd365/logger/v2"
	"log/slog"
	"os"
)

func main() {
	l := logger.New(logger.LevelAll)
	l.SetLevel(logger.LevelDebug)
	msg := "test"
	args := []any{"title", "title-value"}
	l.Trace(msg, args...)
	l.Debug(msg, args...)
	l.Info(msg, args...)
	l.Warn(msg, args...)
	l.Error(msg, args...)
	l.Fatal(msg, args...)
	l.SetHandler(slog.NewJSONHandler(os.Stdout, l.GetOptions()))
	l.Trace(msg, args...)
	l.Debug(msg, args...)
	l.Info(msg, args...)
	l.Warn(msg, args...)
	l.Error(msg, args...)
	l.Fatal(msg, args...)
}
```
