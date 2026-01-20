package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/utils"
)

type authErrorResponse struct {
	Error authErrorObject `json:"error"`
	Meta  authMeta        `json:"meta"`
}

type authErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

type authMeta struct {
	RequestID string `json:"request_id"`
}

func RequireAuth(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearer(c.GetHeader("Authorization"))
		if tokenString == "" {
			respondAuthError(c, http.StatusUnauthorized, constants.AuthUnauthorized, "unauthorized")
			return
		}

		claims, err := utils.ParseToken(tokenString, cfg)
		if err != nil {
			respondAuthError(c, http.StatusUnauthorized, constants.AuthTokenInvalid, "invalid token")
			return
		}

		if claims.TokenType != utils.TokenTypeAccess {
			respondAuthError(c, http.StatusUnauthorized, constants.AuthTokenInvalid, "invalid token type")
			return
		}

		actorID, err := uuid.Parse(claims.Subject)
		if err != nil {
			respondAuthError(c, http.StatusUnauthorized, constants.AuthTokenInvalid, "invalid subject")
			return
		}

		role := claims.Role
		if !role.IsValid() {
			respondAuthError(c, http.StatusUnauthorized, constants.AuthTokenInvalid, "invalid role")
			return
		}

		SetActor(c, actorID, role)
		c.Next()
	}
}

func OptionalAuth(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearer(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.Next()
			return
		}

		claims, err := utils.ParseToken(tokenString, cfg)
		if err != nil || claims.TokenType != utils.TokenTypeAccess {
			c.Next()
			return
		}

		actorID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.Next()
			return
		}

		role := claims.Role
		if !role.IsValid() {
			c.Next()
			return
		}

		SetActor(c, actorID, role)
		c.Next()
	}
}

func extractBearer(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func respondAuthError(c *gin.Context, status int, code, message string) {
	rid := GetRequestID(c)
	c.AbortWithStatusJSON(status, authErrorResponse{
		Error: authErrorObject{
			Code:    code,
			Message: message,
			Details: nil,
		},
		Meta: authMeta{RequestID: rid},
	})
}
