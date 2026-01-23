package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type authServiceStub struct{}

func (authServiceStub) RequestOTP(ctx context.Context, phone, purpose, ip string) (dto.RequestOTPResponse, error) {
	return dto.RequestOTPResponse{RefCode: "ref", ExpiresAt: time.Now().UTC()}, nil
}
func (authServiceStub) VerifyOTP(ctx context.Context, phone, refCode, otpCode, purpose string) error {
	return nil
}
func (authServiceStub) Register(ctx context.Context, req dto.RegisterRequest) (dto.TokenResponse, error) {
	return dto.TokenResponse{AccessToken: "a", RefreshToken: "r"}, nil
}
func (authServiceStub) Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error) {
	return dto.TokenResponse{AccessToken: "a", RefreshToken: "r"}, nil
}
func (authServiceStub) ForgotPasswordRequestOTP(ctx context.Context, phone, ip string) (dto.RequestOTPResponse, error) {
	return dto.RequestOTPResponse{RefCode: "ref", ExpiresAt: time.Now().UTC()}, nil
}
func (authServiceStub) ForgotPasswordConfirm(ctx context.Context, req dto.ForgotPasswordConfirmRequest) error {
	return nil
}
func (authServiceStub) Refresh(ctx context.Context, refreshToken string) (dto.TokenResponse, error) {
	return dto.TokenResponse{AccessToken: "a", RefreshToken: "r"}, nil
}
func (authServiceStub) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func TestAuthHandlers(t *testing.T) {
	router := newTestRouter()
	handler := NewAuthHandler(authServiceStub{})

	router.POST("/auth/request-otp", handler.RequestOTP)
	router.POST("/auth/verify-otp", handler.VerifyOTP)
	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/forgot-password/request-otp", handler.ForgotPasswordRequestOTP)
	router.POST("/auth/forgot-password/confirm", handler.ForgotPasswordConfirm)
	router.POST("/auth/refresh", handler.Refresh)
	router.POST("/auth/logout", handler.Logout)

	cases := []struct {
		name       string
		path       string
		payload    any
		wantStatus int
	}{
		{"request-otp", "/auth/request-otp", dto.RequestOTPRequest{Phone: "0800000000", Purpose: "register"}, http.StatusAccepted},
		{"verify-otp", "/auth/verify-otp", dto.VerifyOTPRequest{Phone: "0800000000", RefCode: "ref", OTPCode: "123456", Purpose: "register"}, http.StatusOK},
		{"register", "/auth/register", dto.RegisterRequest{Phone: "0800000000", RefCode: "ref", Password: "pass", FirstName: "A", LastName: "B"}, http.StatusCreated},
		{"login", "/auth/login", dto.LoginRequest{Phone: "0800000000", Password: "pass"}, http.StatusOK},
		{"forgot-otp", "/auth/forgot-password/request-otp", dto.ForgotPasswordRequestOTPRequest{Phone: "0800000000"}, http.StatusAccepted},
		{"forgot-confirm", "/auth/forgot-password/confirm", dto.ForgotPasswordConfirmRequest{Phone: "0800000000", RefCode: "ref", OTPCode: "123456", NewPassword: "new"}, http.StatusOK},
		{"refresh", "/auth/refresh", dto.RefreshRequest{RefreshToken: "token"}, http.StatusOK},
		{"logout", "/auth/logout", dto.LogoutRequest{RefreshToken: "token"}, http.StatusOK},
	}

	for _, tc := range cases {
		resp := performRequest(router, http.MethodPost, tc.path, tc.payload)
		if resp.Code != tc.wantStatus {
			t.Fatalf("%s: expected %d got %d", tc.name, tc.wantStatus, resp.Code)
		}
		var meta envelopeMeta
		if err := json.Unmarshal(resp.Body.Bytes(), &meta); err != nil {
			t.Fatalf("%s: invalid json", tc.name)
		}
		if meta.Meta.RequestID != testRequestID {
			t.Fatalf("%s: expected request_id", tc.name)
		}
	}
}
