package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type HealthRecordHandler struct {
	service    services.HealthService
	caregivers services.CaregiverService
}

func NewHealthRecordsHandler(service services.HealthService, caregivers services.CaregiverService) *HealthRecordHandler {
	return &HealthRecordHandler{service: service, caregivers: caregivers}
}

func (h *HealthRecordHandler) CreateHealthRecord(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)
	var req dto.CreateHealthRecordRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.CreateHealthRecord(c.Request.Context(), actorID.String(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *HealthRecordHandler) ListHealthRecords(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	userID := c.Query("user_id")

	resolvedUserID, err := authorizePatientAccess(c, h.caregivers, userID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.ListHealthRecords(c.Request.Context(), resolvedUserID, from, to)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *HealthRecordHandler) CreateDailyAssessment(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)
	var req dto.CreateDailyAssessmentRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.CreateDailyAssessment(c.Request.Context(), actorID.String(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *HealthRecordHandler) ListDailyAssessments(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	userID := c.Query("user_id")

	resolvedUserID, err := authorizePatientAccess(c, h.caregivers, userID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.ListDailyAssessments(c.Request.Context(), resolvedUserID, from, to)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}
