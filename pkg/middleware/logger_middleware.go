package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Config struct {
	// Next defines a function to skip this middleware when returned true
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// TimeFormat https://programming.guide/go/format-parse-string-time-date-example.html
	//
	// Optional. Default: 2006-01-02 15:04:05
	TimeFormat string
}

var ConfigDefault = Config{
	Next:       nil,
	TimeFormat: "2006-01-02 15:04:05",
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = ConfigDefault.TimeFormat
	}
	return cfg
}

// Inspiration from:
// https://github.com/edersohe/zflogger/blob/087c6cbef12b25269934b9883d3881d3933f900e/zflogger.go#L49
// https://github.com/gofiber/fiber/blob/master/middleware/logger/logger.go
func ZapLogger(config ...Config) fiber.Handler {

	cfg := configDefault(config...)

	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	return func(c *fiber.Ctx) error {
		// Dont execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Set error handler once
		once.Do(func() {
			errHandler = c.App().Config().ErrorHandler
		})

		start := time.Now()

		chainErr := c.Next()

		requestId := c.Get(fiber.HeaderXRequestID)
		if requestId == "" {
			requestId = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, requestId)
		}

		// Not using sugared logger since performance is important in request logging
		// TODO: Life time of logger object? Should this be created outside?
		logger := zap.L()
		defer logger.Sync()

		// Manually call error handler
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		defer func() {
			rvr := recover()

			var errorMsg error = nil
			errorStack := ""

			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				errorMsg = err
				errorStack = string(debug.Stack())
				c.Status(fiber.StatusInternalServerError)
				c.JSON(map[string]interface{}{
					"status": http.StatusText(http.StatusInternalServerError),
				})
			}

			message := ""
			switch {
			case rvr != nil:
				message = "panic recover"
			case c.Response().StatusCode() >= 500:
				message = "server error"
			case c.Response().StatusCode() >= 400:
				message = "client error"
			case c.Response().StatusCode() >= 300:
				message = "redirect"
			case c.Response().StatusCode() >= 200:
				message = "success"
			case c.Response().StatusCode() >= 100:
				message = "informative"
			default:
				message = "unknown status"
			}

			// TODO: pass log structure as config
			logger.Info(message,
				zap.String("Timestamp", start.Format(cfg.TimeFormat)),
				zap.Int("Status", c.Response().StatusCode()),
				zap.Duration("Duration", time.Since(start)),
				zap.String("IP", c.IP()),
				zap.String("RequestID", requestId),
				zap.String("Method", c.Method()),
				zap.String("Path", c.Path()),
				zap.String("Stacktrace", errorStack),
				zap.NamedError("Error", errorMsg),
			)
		}()

		// c.Next has already been called, so we return nil here
		return nil
	}
}
