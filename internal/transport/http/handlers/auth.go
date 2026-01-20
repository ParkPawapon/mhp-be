package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type AuthHandler struct {
	auth services.AuthService
}

func NewAuthHandler(auth services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req dto.RequestOTPRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.auth.RequestOTP(c.Request.Context(), req.Phone, req.Purpose, c.ClientIP())
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Respond(c, http.StatusAccepted, resp)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	if err := h.auth.VerifyOTP(c.Request.Context(), req.Phone, req.RefCode, req.OTPCode, req.Purpose); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, dto.VerifyOTPResponse{Verified: true})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.auth.Register(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Created(c, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.auth.Login(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *AuthHandler) ForgotPasswordRequestOTP(c *gin.Context) {
	var req dto.ForgotPasswordRequestOTPRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.auth.ForgotPasswordRequestOTP(c.Request.Context(), req.Phone, c.ClientIP())
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.Respond(c, http.StatusAccepted, resp)
}

func (h *AuthHandler) ForgotPasswordConfirm(c *gin.Context) {
	var req dto.ForgotPasswordConfirmRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	if err := h.auth.ForgotPasswordConfirm(c.Request.Context(), req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"updated": true})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	resp, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := bindAndValidateJSON(c, &req); err != nil {
		httpx.Fail(c, err)
		return
	}

	if err := h.auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"revoked": true})
}
