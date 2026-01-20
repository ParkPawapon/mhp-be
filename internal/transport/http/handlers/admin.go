package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type AdminHandler struct {
	service services.AdminService
}

func NewAdminHandler(service services.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) StaffLogin(c *gin.Context) {
	var req dto.StaffLoginRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.StaffLogin(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *AdminHandler) ListPatients(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.service.ListPatients(c.Request.Context(), page, pageSize)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	meta := httpx.PaginationMeta(middleware.GetRequestID(c), page, pageSize, total)
	c.JSON(200, httpx.SuccessResponse{Data: items, Meta: meta})
}

func (h *AdminHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.service.GetPatient(c.Request.Context(), id)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *AdminHandler) ListAdherence(c *gin.Context) {
	patientID := c.Query("patient_id")
	from := c.Query("from")
	to := c.Query("to")
	if patientID == "" {
		httpx.Fail(c, domain.NewError(constants.ValidationFailed, "patient_id required"))
		return
	}
	resp, err := h.service.ListAdherence(c.Request.Context(), patientID, from, to)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}
