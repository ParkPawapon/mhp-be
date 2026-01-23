package services

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type authRepoStub struct {
	otp *db.AuthOtpCode
}

func (s *authRepoStub) CreateOTP(ctx context.Context, otp *db.AuthOtpCode) error {
	s.otp = otp
	s.otp.ID = uuid.New()
	return nil
}
func (s *authRepoStub) FindOTP(ctx context.Context, phone, refCode string) (*db.AuthOtpCode, error) {
	if s.otp == nil || s.otp.PhoneNumber != phone || s.otp.RefCode != refCode {
		return nil, domain.NewError(constants.AuthOTPInvalid, "otp not found")
	}
	if s.otp.IsUsed {
		return nil, domain.NewError(constants.AuthOTPUsed, "otp already used")
	}
	if time.Now().UTC().After(s.otp.ExpiredAt) {
		return nil, domain.NewError(constants.AuthOTPExpired, "otp expired")
	}
	return s.otp, nil
}
func (s *authRepoStub) MarkOTPUsed(ctx context.Context, id uuid.UUID) error {
	if s.otp != nil {
		s.otp.IsUsed = true
	}
	return nil
}

type userRepoStubAuth struct{}

func (userRepoStubAuth) Create(ctx context.Context, user *db.User) error {
	return nil
}
func (userRepoStubAuth) FindByUsername(ctx context.Context, username string) (*db.User, error) {
	return &db.User{ID: uuid.New(), Username: username, PasswordHash: "hash", Role: constants.RolePatient}, nil
}
func (userRepoStubAuth) FindByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	return &db.User{ID: id, Username: "phone", PasswordHash: "hash", Role: constants.RolePatient}, nil
}
func (userRepoStubAuth) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	return nil
}

type smsSenderStub struct{}

func (smsSenderStub) SendOTP(ctx context.Context, phone, otpCode, refCode string) error {
	return nil
}

func newTestAuthService(t *testing.T) (AuthService, *authRepoStub, *redis.Client) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	cfg := config.Config{
		OTP:       config.OTPConfig{TTL: 5 * time.Minute, Digits: 6, RefCodeLength: 6},
		JWT:       config.JWTConfig{Issuer: "test", Secret: "test-secret-123456789012345678901234567890", AccessTTL: 15 * time.Minute, RefreshTTL: 24 * time.Hour},
		RateLimit: config.RateLimitConfig{OTPPerPhone: 1, OTPPerIP: 1, Window: time.Minute},
	}
	repo := &authRepoStub{}
	svc := NewAuthService(cfg, repo, userRepoStubAuth{}, rdb, smsSenderStub{})
	return svc, repo, rdb
}

func TestRequestOTPRateLimit(t *testing.T) {
	svc, _, _ := newTestAuthService(t)

	_, err := svc.RequestOTP(context.Background(), "0800000000", "register", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.RequestOTP(context.Background(), "0800000000", "register", "127.0.0.1")
	if err == nil {
		t.Fatalf("expected rate limit error")
	}
}

func TestVerifyOTPSetsRedisFlag(t *testing.T) {
	svc, repo, rdb := newTestAuthService(t)
	phone := "0800000000"
	refCode := "ref123"
	otp := "123456"
	repo.otp = &db.AuthOtpCode{ID: uuid.New(), PhoneNumber: phone, RefCode: refCode, OtpCode: otp, ExpiredAt: time.Now().Add(time.Minute)}

	if err := svc.VerifyOTP(context.Background(), phone, refCode, otp, "register"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	key := "otp:verified:register:" + phone + ":" + refCode
	if val, err := rdb.Get(context.Background(), key).Result(); err != nil || val == "" {
		t.Fatalf("expected verified key set")
	}
}

func TestRequestOTPMissingPhone(t *testing.T) {
	svc, _, _ := newTestAuthService(t)
	_, err := svc.RequestOTP(context.Background(), "", "register", "")
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if appErr, ok := domain.AsAppError(err); !ok || appErr.Code != constants.ValidationFailed {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForgotPasswordConfirmUsesOTP(t *testing.T) {
	svc, repo, _ := newTestAuthService(t)
	phone := "0800000000"
	refCode := "ref123"
	otp := "123456"
	repo.otp = &db.AuthOtpCode{ID: uuid.New(), PhoneNumber: phone, RefCode: refCode, OtpCode: otp, ExpiredAt: time.Now().Add(time.Minute)}

	err := svc.ForgotPasswordConfirm(context.Background(), dto.ForgotPasswordConfirmRequest{
		Phone:       phone,
		RefCode:     refCode,
		OTPCode:     otp,
		NewPassword: "newpass",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
