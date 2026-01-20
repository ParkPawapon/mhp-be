package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type AuthRepository interface {
	CreateOTP(ctx context.Context, otp *db.AuthOtpCode) error
	FindOTP(ctx context.Context, phone, refCode string) (*db.AuthOtpCode, error)
	MarkOTPUsed(ctx context.Context, id uuid.UUID) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(dbConn *gorm.DB) AuthRepository {
	return &authRepository{db: dbConn}
}

func (r *authRepository) CreateOTP(ctx context.Context, otp *db.AuthOtpCode) error {
	if err := r.db.WithContext(ctx).Create(otp).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create otp failed", err)
	}
	return nil
}

func (r *authRepository) FindOTP(ctx context.Context, phone, refCode string) (*db.AuthOtpCode, error) {
	var otp db.AuthOtpCode
	err := r.db.WithContext(ctx).
		Where("phone_number = ? AND ref_code = ?", phone, refCode).
		Order("created_at desc").
		First(&otp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.AuthOTPInvalid, "otp not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find otp failed", err)
	}
	if otp.IsUsed {
		return nil, domain.NewError(constants.AuthOTPUsed, "otp already used")
	}
	if time.Now().UTC().After(otp.ExpiredAt) {
		return nil, domain.NewError(constants.AuthOTPExpired, "otp expired")
	}
	return &otp, nil
}

func (r *authRepository) MarkOTPUsed(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&db.AuthOtpCode{}).Where("id = ?", id).Update("is_used", true).Error; err != nil {
		return domain.WrapError(constants.InternalError, "mark otp used failed", err)
	}
	return nil
}
