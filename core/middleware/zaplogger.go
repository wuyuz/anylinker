package middleware

import (
	"anylinker/common/log"
	"anylinker/core/utils/resp"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// ZapLogger zap logget middle
func ZapLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		url := c.BaseURL()
		c.Next()
		latency := time.Now().Sub(start) * 1000
		statuscode,err := strconv.Atoi(c.Get("statuscode","10000"))
		reqip := c.IP()
		method := c.Method()
		bodySize := len(c.Body())

		fields := []zap.Field{
			zap.String("uid", c.Get("uid")),
			zap.String("method", strings.ToLower(method)),
			zap.Int("statuscode", statuscode),
			zap.String("reqip", reqip),
			zap.Duration("latency(ms)", latency),
			zap.String("url", url),
		}

		if bodySize > 0 {
			fields = append(fields, zap.Int("respsize", bodySize))
		}

		switch {
		case statuscode < resp.ErrBadRequest:
			log.Debug("fiber", fields...)
		case statuscode < resp.ErrInternalServer:
			log.Warn("fiber", fields...)
		case statuscode >= resp.ErrInternalServer:
			log.Error("fiber", fields...)
		}
		return err
	}
}
