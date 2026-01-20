package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := strings.TrimSpace(c.GetHeader("X-Request-Id"))
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(constants.RequestIDKey, rid)
		c.Writer.Header().Set("X-Request-Id", rid)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(constants.RequestIDKey); ok {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}
