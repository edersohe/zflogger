## zerolog + fiber = zflogger

The middleware contains functionality of requestid + logger + recover for request traceability

## basic example

```go
package main

import (
    "os"

    "github.com/edersohe/zflogger"
    "github.com/gofiber/fiber"
    "github.com/rs/zerolog"
)

func main() {
    log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)

    app := fiber.New()

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

## example with filter

```go
package main

import (
    "os"

    "github.com/edersohe/zflogger"
    "github.com/gofiber/fiber"
    "github.com/rs/zerolog"
)

// this filter apply logger middleware when the context path is "/error",
// else the zflogger.Middleware are skipped and the flow continue
func filter(c *fiber.Ctx) bool {
    return c.Path() != "/error"
}

func main() {
    log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)

    app := fiber.New()

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

## benchmark 

* Mac BooK Pro (Retina, 13-inch, Mid 2014) 
* Procesor 2.6 GHz Dual-Core Intel Core i5
* Memory 8 GB 1600 MHz DDR3

```sh
wrk -t12 -c250 -d30s http://127.0.0.1:3000/
Running 30s test @ http://127.0.0.1:3000/
  12 threads and 250 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    17.83ms   23.41ms 240.55ms   91.84%
    Req/Sec     1.58k   603.26     2.89k    66.73%
  564256 requests in 30.09s, 97.94MB read
Requests/sec:  18749.94
Transfer/sec:      3.25MB
```