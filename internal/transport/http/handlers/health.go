package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/database/postgres"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redisClient}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	httpx.OK(c, gin.H{"status": "ok"})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := postgres.Ping(ctx, h.db); err != nil {
		dbStatus = "error"
	}

	redisStatus := "ok"
	if err := h.redis.Ping(ctx).Err(); err != nil {
		redisStatus = "error"
	}

	statusCode := http.StatusOK
	if dbStatus != "ok" || redisStatus != "ok" {
		statusCode = http.StatusServiceUnavailable
	}

	httpx.Respond(c, statusCode, gin.H{
		"database": dbStatus,
		"redis":    redisStatus,
	})
}
