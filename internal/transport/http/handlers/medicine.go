package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type MedicineHandler struct {
	service services.MedicineService
}

func NewMedicineHandler(service services.MedicineService) *MedicineHandler {
	return &MedicineHandler{service: service}
}

func (h *MedicineHandler) ListMaster(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.service.ListMaster(c.Request.Context(), page, pageSize)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	meta := httpx.PaginationMeta(middleware.GetRequestID(c), page, pageSize, total)
	c.JSON(200, httpx.SuccessResponse{Data: items, Meta: meta})
}

func (h *MedicineHandler) ListCategories(c *gin.Context) {
	resp, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *MedicineHandler) ListCategoryItems(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.service.ListCategoryItems(c.Request.Context(), id)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *MedicineHandler) GetDosageOptions(c *gin.Context) {
	httpx.OK(c, h.service.GetDosageOptions(c.Request.Context()))
}

func (h *MedicineHandler) GetMealTimingOptions(c *gin.Context) {
	httpx.OK(c, h.service.GetMealTimingOptions(c.Request.Context()))
}

func (h *MedicineHandler) CreatePatientMedicine(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)

	var req dto.CreatePatientMedicineRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.CreatePatientMedicine(c.Request.Context(), actorID.String(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *MedicineHandler) ListPatientMedicines(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)
	resp, err := h.service.ListPatientMedicines(c.Request.Context(), actorID.String())
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *MedicineHandler) UpdatePatientMedicine(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdatePatientMedicineRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	if err := h.service.UpdatePatientMedicine(c.Request.Context(), id, req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"updated": true})
}

func (h *MedicineHandler) DeletePatientMedicine(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePatientMedicine(c.Request.Context(), id); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"deleted": true})
}

func (h *MedicineHandler) CreateSchedule(c *gin.Context) {
	id := c.Param("id")

	var req dto.CreateMedicineScheduleRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.service.CreateSchedule(c.Request.Context(), id, req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *MedicineHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSchedule(c.Request.Context(), id); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"deleted": true})
}
