package middleware

import (
	"context"
	"strconv"
	"time"
	"weatherApi/internal/common/constants"
	"weatherApi/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.NewString()
		ctx := context.WithValue(c.Request.Context(), constants.TraceID, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// PrometheusMiddleware returns a gin.HandlerFunc that records HTTP metrics using the provided HTTPMetrics.
func PrometheusMiddleware(httpMetrics *metrics.HTTPMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" || c.FullPath() == "/api/v1/health" || c.FullPath() == "/assets/*filepath" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start).Seconds()
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		httpMetrics.ObserveRequestDuration(c.Request.Method, route, strconv.Itoa(status), duration)
	}
}
