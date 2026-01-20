package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type IntakeHandler struct {
	service   services.IntakeService
	caregivers services.CaregiverService
}

func NewIntakeHandler(service services.IntakeService, caregivers services.CaregiverService) *IntakeHandler {
	return &IntakeHandler{service: service, caregivers: caregivers}
}

func (h *IntakeHandler) CreateIntake(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)
	var req dto.CreateIntakeRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.CreateIntake(c.Request.Context(), actorID.String(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *IntakeHandler) ListHistory(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	userID := c.Query("user_id")

	resolvedUserID, err := authorizePatientAccess(c, h.caregivers, userID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.ListHistory(c.Request.Context(), resolvedUserID, from, to)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}
