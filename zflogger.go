package zflogger

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

//Filter type used by the Middleware
type Filter func(*fiber.Ctx) bool

//New logger with level specified by string, if isn't a valid level string returns a logger with zerolog.NoLevel
//
//Valid level strings: trace, debug, info, warn, error, fatal and panic
func New(out io.Writer, level string) *zerolog.Logger {
	lvl, _ := zerolog.ParseLevel(level)
	logger := zerolog.New(out).With().Timestamp().Logger().Level(lvl)
	return &logger
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type logFields struct {
	ID         string
	RemoteIP   string
	Host       string
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Latency    float64
	Error      error
	Stack      []byte
}

func (lf *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", lf.ID).
		Str("remote_ip", lf.RemoteIP).
		Str("host", lf.Host).
		Str("method", lf.Method).
		Str("path", lf.Path).
		Str("protocol", lf.Protocol).
		Int("status_code", lf.StatusCode).
		Float64("latency", lf.Latency).
		Str("tag", "request")

	if lf.Error != nil {
		e.Err(lf.Error)
	}

	if lf.Stack != nil {
		e.Bytes("stack", lf.Stack)
	}
}

//Marshal a interface with jsoniter
func Marshal(i interface{}) []byte {
	marshaled, err := json.Marshal(&i)

	if err != nil {
		return []byte("Marshal error: " + err.Error())
	}

	return marshaled
}

//Middleware RequestID, Logger and Recover by convenience
func Middleware(log *zerolog.Logger, filter Filter) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		if filter != nil && filter(c) {
			c.Next()
			return
		}

		start := time.Now()

		rid := c.Get(fiber.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, rid)
		}

		fields := &logFields{
			ID:       rid,
			RemoteIP: c.IP(),
			Method:   c.Method(),
			Host:     c.Hostname(),
			Path:     c.Path(),
			Protocol: c.Protocol(),
		}

		defer func() {
			rvr := recover()

			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				fields.Error = err
				fields.Stack = debug.Stack()

				c.Status(http.StatusInternalServerError)
				c.JSON(map[string]interface{}{
					"status": http.StatusText(http.StatusInternalServerError),
				})
			}

			fields.StatusCode = c.Fasthttp.Response.StatusCode()
			fields.Latency = time.Since(start).Seconds()

			switch {
			case rvr != nil:
				log.Error().EmbedObject(fields).Msg("panic recover")
			case fields.StatusCode >= 500:
				log.Error().EmbedObject(fields).Msg("server error")
			case fields.StatusCode >= 400:
				log.Error().EmbedObject(fields).Msg("client error")
			case fields.StatusCode >= 300:
				log.Warn().EmbedObject(fields).Msg("redirect")
			case fields.StatusCode >= 200:
				log.Info().EmbedObject(fields).Msg("success")
			case fields.StatusCode >= 100:
				log.Info().EmbedObject(fields).Msg("informative")
			default:
				log.Warn().EmbedObject(fields).Msg("unknown status")
			}
		}()

		c.Next()
	}
}
