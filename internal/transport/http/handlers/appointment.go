package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type AppointmentHandler struct {
	service    services.AppointmentService
	caregivers services.CaregiverService
}

func NewAppointmentHandler(service services.AppointmentService, caregivers services.CaregiverService) *AppointmentHandler {
	return &AppointmentHandler{service: service, caregivers: caregivers}
}

func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	userID := c.Query("user_id")
	resolvedUserID, err := authorizePatientAccess(c, h.caregivers, userID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.ListAppointments(c.Request.Context(), resolvedUserID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)
	var req dto.CreateAppointmentRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}
	resp, err := h.service.CreateAppointment(c.Request.Context(), actorID.String(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateAppointmentStatusRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}
	if err := h.service.UpdateStatus(c.Request.Context(), id, req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"updated": true})
}

func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteAppointment(c.Request.Context(), id); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"deleted": true})
}

func (h *AppointmentHandler) CreateNurseVisitNote(c *gin.Context) {
	id := c.Param("id")
	actorID, _ := middleware.GetActorID(c)
	var req dto.CreateNurseVisitNoteRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}
	if err := h.service.CreateNurseVisitNote(c.Request.Context(), id, actorID.String(), req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, gin.H{"created": true})
}

func (h *AppointmentHandler) ListVisitHistory(c *gin.Context) {
	userID := c.Query("user_id")
	resolvedUserID, err := authorizePatientAccess(c, h.caregivers, userID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.ListVisitHistory(c.Request.Context(), resolvedUserID)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}
