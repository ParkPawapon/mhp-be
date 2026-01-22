package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type UserHandler struct {
	users services.UserService
}

func NewUserHandler(users services.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Me(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)
	role, _ := middleware.GetRole(c)

	resp, err := h.users.GetMe(c.Request.Context(), actorID, role)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)

	var req dto.UpdateProfileRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	if err := h.users.UpdateProfile(c.Request.Context(), actorID, req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"updated": true})
}

func (h *UserHandler) SaveDeviceToken(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)

	var req dto.DeviceTokenRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	if err := h.users.SaveDeviceToken(c.Request.Context(), actorID, req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"saved": true})
}

func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	actorID, _ := middleware.GetActorID(c)

	var req dto.UpdatePreferencesRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.users.UpdatePreferences(c.Request.Context(), actorID, req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}
