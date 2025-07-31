package middleware

import (
	"context"
	"weatherApi/internal/common/constants"

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
