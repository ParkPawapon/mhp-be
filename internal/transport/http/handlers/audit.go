package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type AuditHandler struct {
	service services.AuditService
}

func NewAuditHandler(service services.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	page, pageSize := parsePagination(c)
	from := c.Query("from")
	to := c.Query("to")
	actorID := c.Query("actor_id")
	actionType := c.Query("action_type")

	items, total, err := h.service.ListAuditLogs(c.Request.Context(), page, pageSize, from, to, actorID, actionType)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	meta := httpx.PaginationMeta(middleware.GetRequestID(c), page, pageSize, total)
	c.JSON(200, httpx.SuccessResponse{Data: items, Meta: meta})
}
