package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type recoveryErrorResponse struct {
	Error recoveryErrorObject `json:"error"`
	Meta  recoveryMeta        `json:"meta"`
}

type recoveryErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

type recoveryMeta struct {
	RequestID string `json:"request_id"`
}

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", zap.Any("error", r))
				rid := GetRequestID(c)
				c.AbortWithStatusJSON(http.StatusInternalServerError, recoveryErrorResponse{
					Error: recoveryErrorObject{
						Code:    constants.InternalError,
						Message: "internal error",
						Details: nil,
					},
					Meta: recoveryMeta{RequestID: rid},
				})
			}
		}()
		c.Next()
	}
}
