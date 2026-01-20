package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ParkPawapon/mhp-be/internal/observability"
)

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		rid := GetRequestID(c)
		actorID, _ := GetActorID(c)
		role, _ := GetRole(c)
		status := c.Writer.Status()
		latencyMs := time.Since(start).Milliseconds()
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		observability.HTTPRequests.WithLabelValues(c.Request.Method, route, statusText(status)).Inc()
		observability.HTTPDuration.WithLabelValues(c.Request.Method, route, statusText(status)).Observe(float64(latencyMs) / 1000.0)

		logger.Info("request",
			zap.String("request_id", rid),
			zap.String("actor_id", actorID.String()),
			zap.String("role", string(role)),
			zap.String("route", route),
			zap.String("method", c.Request.Method),
			zap.Int("status", status),
			zap.Int64("latency_ms", latencyMs),
			zap.String("ip", c.ClientIP()),
		)
	}
}

func statusText(status int) string {
	return strconv.Itoa(status)
}
