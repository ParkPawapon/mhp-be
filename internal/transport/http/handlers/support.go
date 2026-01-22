package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type SupportHandler struct {
	service services.SupportService
}

func NewSupportHandler(service services.SupportService) *SupportHandler {
	return &SupportHandler{service: service}
}

func (h *SupportHandler) CreateChatRequest(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)

	var req dto.SupportChatRequestCreateRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.CreateChatRequest(c.Request.Context(), actorID.String(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *SupportHandler) ListChatRequests(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.service.ListChatRequests(c.Request.Context(), page, pageSize)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	meta := httpx.PaginationMeta(middleware.GetRequestID(c), page, pageSize, total)
	c.JSON(200, httpx.SuccessResponse{Data: items, Meta: meta})
}

func (h *SupportHandler) EmergencyInfo(c *gin.Context) {
	resp := h.service.GetEmergencyInfo(c.Request.Context())
	httpx.OK(c, resp)
}
