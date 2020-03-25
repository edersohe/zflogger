## zerolog + fiber = zflogger

The middleware contains functionality of requestid + logger(zerolog) + recover for request traceability

## basic example

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

    app.Use(zflogger.Middleware(log, nil))

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

## advanced example

```go

package main

import (
    "os"

    "github.com/edersohe/zflogger"
    "github.com/gofiber/fiber"
)


// this filter apply logger middleware when the context path is "/error",
// else the zflogger.Middleware are skipped and the flow continue
func filter(c *fiber.ctx) bool {
    return c.Path() != "/error"
}

func main() {
    app := fiber.New()
    log := zflogger.New(os.Stderr, "debug")

    app.Use(zflogger.Middleware(log, filter))

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

## test

```sh
curl http://localhost:3000
curl http://localhost:3000/error
```