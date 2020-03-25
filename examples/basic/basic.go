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
