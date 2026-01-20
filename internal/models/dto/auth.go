package dto

import "time"

type RequestOTPRequest struct {
	Phone   string `json:"phone" validate:"required"`
	Purpose string `json:"purpose" validate:"required"`
}

type RequestOTPResponse struct {
	RefCode   string    `json:"ref_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type VerifyOTPRequest struct {
	Phone   string `json:"phone" validate:"required"`
	RefCode string `json:"ref_code" validate:"required"`
	OTPCode string `json:"otp_code" validate:"required"`
	Purpose string `json:"purpose" validate:"required"`
}

type VerifyOTPResponse struct {
	Verified bool `json:"verified"`
}

type RegisterRequest struct {
	Phone     string `json:"phone" validate:"required"`
	RefCode   string `json:"ref_code" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequestOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type ForgotPasswordConfirmRequest struct {
	Phone       string `json:"phone" validate:"required"`
	RefCode     string `json:"ref_code" validate:"required"`
	OTPCode     string `json:"otp_code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
