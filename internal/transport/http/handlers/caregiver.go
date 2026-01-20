package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type CaregiverHandler struct {
	service services.CaregiverService
}

func NewCaregiverHandler(service services.CaregiverService) *CaregiverHandler {
	return &CaregiverHandler{service: service}
}

func (h *CaregiverHandler) CreateAssignment(c *gin.Context) {
	var req dto.CreateCaregiverAssignmentRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.CreateAssignment(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *CaregiverHandler) ListAssignments(c *gin.Context) {
	patientID := c.Query("patient_id")
	if patientID == "" {
		httpx.Fail(c, domain.NewError(constants.ValidationFailed, "patient_id required"))
		return
	}
	pid, err := uuid.Parse(patientID)
	if err != nil {
		httpx.Fail(c, domain.NewError(constants.ValidationFailed, "invalid patient_id"))
		return
	}

	resp, err := h.service.ListAssignments(c.Request.Context(), pid)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}
