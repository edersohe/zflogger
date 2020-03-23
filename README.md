## zerolog + fiber = zflogger

The middleware contains functionality of requestid + logger(zerolog) + recover for request traceability

*Open browser http://localhost:3000 and http://localhost:3000/error and check the output in the console*

## example

```go

package main

import (
	"os"

	"github.com/edersohe/zflogger"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()
	log := zflogger.New(os.Stderr, "debug")

	app.Use(zflogger.Middleware(log))

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("Hello, World!")
	})

	app.Get("/error", func(c *fiber.Ctx) {
		a := 0
		c.JSON(1 / a)
	})

	log.Fatal().Err(app.Listen(3000)).Str("tag", "server").Send()
}
```