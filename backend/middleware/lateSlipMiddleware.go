package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()

		// Set request ID in header and context
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		c.Next()
	}
}
