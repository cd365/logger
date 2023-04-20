# FOR EXAMPLE

```go
package main

import (
	"github.com/cd365/logger"
	"log"
	"os"
)

func main() {
	l := logger.New(os.Stdout, log.Ldate|log.Lmicroseconds|log.Lshortfile, logger.ALL)
	defer l.Close()
	l.Info(112233)
	l.Infof("Hello World")
	l.Infof("Hello %s", "Jerry")
}
```
