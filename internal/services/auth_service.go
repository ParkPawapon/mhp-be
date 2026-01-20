package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
	"github.com/ParkPawapon/mhp-be/internal/utils"
)

type AuthService interface {
	RequestOTP(ctx context.Context, phone, purpose, ip string) (dto.RequestOTPResponse, error)
	VerifyOTP(ctx context.Context, phone, refCode, otpCode, purpose string) error
	Register(ctx context.Context, req dto.RegisterRequest) (dto.TokenResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error)
	ForgotPasswordRequestOTP(ctx context.Context, phone, ip string) (dto.RequestOTPResponse, error)
	ForgotPasswordConfirm(ctx context.Context, req dto.ForgotPasswordConfirmRequest) error
	Refresh(ctx context.Context, refreshToken string) (dto.TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	cfg     config.Config
	authRepo repositories.AuthRepository
	userRepo repositories.UserRepository
	redis   *redis.Client
	sms     SmsSender
}

func NewAuthService(cfg config.Config, authRepo repositories.AuthRepository, userRepo repositories.UserRepository, redisClient *redis.Client, sms SmsSender) AuthService {
	return &authService{
		cfg:      cfg,
		authRepo: authRepo,
		userRepo: userRepo,
		redis:    redisClient,
		sms:      sms,
	}
}

func (s *authService) RequestOTP(ctx context.Context, phone, purpose, ip string) (dto.RequestOTPResponse, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return dto.RequestOTPResponse{}, domain.NewError(constants.ValidationFailed, "phone required")
	}

	if err := s.checkRateLimit(ctx, phone, ip); err != nil {
		return dto.RequestOTPResponse{}, err
	}

	otpCode, err := utils.RandomDigits(s.cfg.OTP.Digits)
	if err != nil {
		return dto.RequestOTPResponse{}, domain.WrapError(constants.InternalError, "generate otp failed", err)
	}
	refCode, err := utils.RandomRefCode(s.cfg.OTP.RefCodeLength)
	if err != nil {
		return dto.RequestOTPResponse{}, domain.WrapError(constants.InternalError, "generate ref code failed", err)
	}

	expiresAt := time.Now().UTC().Add(s.cfg.OTP.TTL)
	record := &db.AuthOtpCode{
		PhoneNumber: phone,
		OtpCode:     otpCode,
		RefCode:     refCode,
		ExpiredAt:   expiresAt,
		IsUsed:      false,
	}

	if err := s.authRepo.CreateOTP(ctx, record); err != nil {
		return dto.RequestOTPResponse{}, err
	}

	if s.sms != nil {
		_ = s.sms.SendOTP(phone, otpCode, refCode)
	}

	return dto.RequestOTPResponse{RefCode: refCode, ExpiresAt: expiresAt}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, phone, refCode, otpCode, purpose string) error {
	phone = strings.TrimSpace(phone)
	refCode = strings.TrimSpace(refCode)
	otpCode = strings.TrimSpace(otpCode)
	if phone == "" || refCode == "" || otpCode == "" {
		return domain.NewError(constants.ValidationFailed, "invalid input")
	}

	record, err := s.authRepo.FindOTP(ctx, phone, refCode)
	if err != nil {
		return err
	}
	if record.OtpCode != otpCode {
		return domain.NewError(constants.AuthOTPInvalid, "otp invalid")
	}

	if err := s.authRepo.MarkOTPUsed(ctx, record.ID); err != nil {
		return err
	}

	key := verifiedOTPKey(purpose, phone, refCode)
	if err := s.redis.Set(ctx, key, "1", s.cfg.OTP.TTL).Err(); err != nil {
		return domain.WrapError(constants.InternalError, "set otp verified failed", err)
	}

	return nil
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (dto.TokenResponse, error) {
	phone := strings.TrimSpace(req.Phone)
	refCode := strings.TrimSpace(req.RefCode)
	if err := s.requireOTPVerified(ctx, "register", phone, refCode); err != nil {
		return dto.TokenResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.TokenResponse{}, domain.WrapError(constants.InternalError, "hash password failed", err)
	}

	user := &db.User{
		Username:     phone,
		PasswordHash: string(hash),
		Role:         constants.RolePatient,
		IsActive:     true,
		IsVerified:   true,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return dto.TokenResponse{}, err
	}

	return s.issueTokens(ctx, user.ID, user.Role)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error) {
	phone := strings.TrimSpace(req.Phone)
	user, err := s.userRepo.FindByUsername(ctx, phone)
	if err != nil {
		return dto.TokenResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.TokenResponse{}, domain.NewError(constants.AuthInvalidCredentials, "invalid credentials")
	}
	return s.issueTokens(ctx, user.ID, user.Role)
}

func (s *authService) ForgotPasswordRequestOTP(ctx context.Context, phone, ip string) (dto.RequestOTPResponse, error) {
	return s.RequestOTP(ctx, phone, "forgot_password", ip)
}

func (s *authService) ForgotPasswordConfirm(ctx context.Context, req dto.ForgotPasswordConfirmRequest) error {
	phone := strings.TrimSpace(req.Phone)
	if err := s.VerifyOTP(ctx, phone, req.RefCode, req.OTPCode, "forgot_password"); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return domain.WrapError(constants.InternalError, "hash password failed", err)
	}

	user, err := s.userRepo.FindByUsername(ctx, phone)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, user.ID, string(hash))
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (dto.TokenResponse, error) {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWT)
	if err != nil {
		return dto.TokenResponse{}, domain.NewError(constants.AuthTokenInvalid, "invalid token")
	}
	if claims.TokenType != utils.TokenTypeRefresh {
		return dto.TokenResponse{}, domain.NewError(constants.AuthTokenInvalid, "invalid token type")
	}
	if !claims.Role.IsValid() {
		return dto.TokenResponse{}, domain.NewError(constants.AuthTokenInvalid, "invalid role")
	}
	if claims.SessionID == "" {
		return dto.TokenResponse{}, domain.NewError(constants.AuthTokenInvalid, "invalid session")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return dto.TokenResponse{}, domain.NewError(constants.AuthTokenInvalid, "invalid subject")
	}

	if err := s.verifyRefreshSession(ctx, userID, claims.SessionID, refreshToken); err != nil {
		return dto.TokenResponse{}, err
	}

	return s.rotateTokens(ctx, userID, claims.Role, claims.SessionID)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWT)
	if err != nil {
		return domain.NewError(constants.AuthTokenInvalid, "invalid token")
	}
	if claims.TokenType != utils.TokenTypeRefresh {
		return domain.NewError(constants.AuthTokenInvalid, "invalid token type")
	}
	if !claims.Role.IsValid() {
		return domain.NewError(constants.AuthTokenInvalid, "invalid role")
	}
	if claims.SessionID == "" {
		return domain.NewError(constants.AuthTokenInvalid, "invalid session")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return domain.NewError(constants.AuthTokenInvalid, "invalid subject")
	}

	key := refreshSessionKey(userID, claims.SessionID)
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return domain.WrapError(constants.InternalError, "logout failed", err)
	}
	return nil
}

func (s *authService) issueTokens(ctx context.Context, userID uuid.UUID, role constants.Role) (dto.TokenResponse, error) {
	sessionID := uuid.New()
	accessToken, err := utils.NewAccessToken(userID, role, sessionID, s.cfg.JWT)
	if err != nil {
		return dto.TokenResponse{}, domain.WrapError(constants.InternalError, "access token failed", err)
	}
	refreshToken, err := utils.NewRefreshToken(userID, role, sessionID, s.cfg.JWT)
	if err != nil {
		return dto.TokenResponse{}, domain.WrapError(constants.InternalError, "refresh token failed", err)
	}

	if err := s.storeRefreshSession(ctx, userID, sessionID.String(), refreshToken); err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authService) rotateTokens(ctx context.Context, userID uuid.UUID, role constants.Role, oldSessionID string) (dto.TokenResponse, error) {
	oldKey := refreshSessionKey(userID, oldSessionID)
	if err := s.redis.Del(ctx, oldKey).Err(); err != nil {
		return dto.TokenResponse{}, domain.WrapError(constants.InternalError, "revoke old session failed", err)
	}
	return s.issueTokens(ctx, userID, role)
}

func (s *authService) storeRefreshSession(ctx context.Context, userID uuid.UUID, sessionID string, refreshToken string) error {
	key := refreshSessionKey(userID, sessionID)
	hash := utils.HashToken(refreshToken)
	if err := s.redis.Set(ctx, key, hash, s.cfg.JWT.RefreshTTL).Err(); err != nil {
		return domain.WrapError(constants.InternalError, "store refresh session failed", err)
	}
	return nil
}

func (s *authService) verifyRefreshSession(ctx context.Context, userID uuid.UUID, sessionID string, refreshToken string) error {
	key := refreshSessionKey(userID, sessionID)
	stored, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.NewError(constants.AuthTokenInvalid, "session revoked")
		}
		return domain.WrapError(constants.InternalError, "refresh session lookup failed", err)
	}
	if stored != utils.HashToken(refreshToken) {
		return domain.NewError(constants.AuthTokenInvalid, "session invalid")
	}
	return nil
}

func (s *authService) checkRateLimit(ctx context.Context, phone, ip string) error {
	if err := s.rateLimitKey(ctx, fmt.Sprintf("otp:rate:phone:%s", phone), s.cfg.RateLimit.OTPPerPhone); err != nil {
		return err
	}
	if ip != "" {
		if err := s.rateLimitKey(ctx, fmt.Sprintf("otp:rate:ip:%s", ip), s.cfg.RateLimit.OTPPerIP); err != nil {
			return err
		}
	}
	return nil
}

func (s *authService) rateLimitKey(ctx context.Context, key string, limit int) error {
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return domain.WrapError(constants.InternalError, "rate limit failed", err)
	}
	if count == 1 {
		_ = s.redis.Expire(ctx, key, s.cfg.RateLimit.Window).Err()
	}
	if count > int64(limit) {
		return domain.NewError(constants.RateLimited, "rate limited")
	}
	return nil
}

func (s *authService) requireOTPVerified(ctx context.Context, purpose, phone, refCode string) error {
	key := verifiedOTPKey(purpose, phone, refCode)
	val, err := s.redis.Get(ctx, key).Result()
	if err != nil || val == "" {
		return domain.NewError(constants.AuthOTPInvalid, "otp not verified")
	}
	return nil
}

func verifiedOTPKey(purpose, phone, refCode string) string {
	return fmt.Sprintf("otp:verified:%s:%s:%s", purpose, phone, refCode)
}

func refreshSessionKey(userID uuid.UUID, sessionID string) string {
	return fmt.Sprintf("auth:refresh:%s:%s", userID.String(), sessionID)
}
