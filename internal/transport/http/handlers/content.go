package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type ContentHandler struct {
	service services.ContentService
}

func NewContentHandler(service services.ContentService) *ContentHandler {
	return &ContentHandler{service: service}
}

func (h *ContentHandler) ListHealthContent(c *gin.Context) {
	publishedOnly := false
	if v := c.Query("published"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			publishedOnly = parsed
		}
	}

	resp, err := h.service.ListHealthContent(c.Request.Context(), publishedOnly)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *ContentHandler) CreateHealthContent(c *gin.Context) {
	var req dto.CreateHealthContentRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}
	resp, err := h.service.CreateHealthContent(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *ContentHandler) UpdateHealthContent(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateHealthContentRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}
	if err := h.service.UpdateHealthContent(c.Request.Context(), id, req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"updated": true})
}

func (h *ContentHandler) PublishHealthContent(c *gin.Context) {
	id := c.Param("id")
	var req dto.PublishHealthContentRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}
	if err := h.service.PublishHealthContent(c.Request.Context(), id, req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"published": req.IsPublished})
}
