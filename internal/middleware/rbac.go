package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type rbacErrorResponse struct {
	Error rbacErrorObject `json:"error"`
	Meta  rbacMeta        `json:"meta"`
}

type rbacErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

type rbacMeta struct {
	RequestID string `json:"request_id"`
}

func RequireRoles(roles ...constants.Role) gin.HandlerFunc {
	roleSet := map[constants.Role]struct{}{}
	for _, role := range roles {
		roleSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role, ok := GetRole(c)
		if !ok {
			respondRBACError(c, http.StatusUnauthorized, constants.AuthUnauthorized, "unauthorized")
			return
		}
		if _, allowed := roleSet[role]; !allowed {
			respondRBACError(c, http.StatusForbidden, constants.AuthForbidden, "forbidden")
			return
		}
		c.Next()
	}
}

func respondRBACError(c *gin.Context, status int, code, message string) {
	rid := GetRequestID(c)
	c.AbortWithStatusJSON(status, rbacErrorResponse{
		Error: rbacErrorObject{
			Code:    code,
			Message: message,
			Details: nil,
		},
		Meta: rbacMeta{RequestID: rid},
	})
}
