package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

func SetActor(c *gin.Context, actorID uuid.UUID, role constants.Role) {
	c.Set(constants.ActorIDKey, actorID)
	c.Set(constants.RoleKey, role)
}

func GetActorID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(constants.ActorIDKey)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

func GetRole(c *gin.Context) (constants.Role, bool) {
	v, ok := c.Get(constants.RoleKey)
	if !ok {
		return "", false
	}
	role, ok := v.(constants.Role)
	return role, ok
}
