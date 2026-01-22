package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type DeviceTokenRepository interface {
	Save(ctx context.Context, token *db.DeviceToken) error
}

type deviceTokenRepository struct {
	db *gorm.DB
}

func NewDeviceTokenRepository(dbConn *gorm.DB) DeviceTokenRepository {
	return &deviceTokenRepository{db: dbConn}
}

func (r *deviceTokenRepository) Save(ctx context.Context, token *db.DeviceToken) error {
	var existing db.DeviceToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND platform = ? AND token = ?", token.UserID, token.Platform, token.Token).
		First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
				return domain.WrapError(constants.InternalError, "create device token failed", err)
			}
			return nil
		}
		return domain.WrapError(constants.InternalError, "find device token failed", err)
	}

	if err := r.db.WithContext(ctx).Model(&db.DeviceToken{}).Where("id = ?", existing.ID).Update("is_active", true).Error; err != nil {
		return domain.WrapError(constants.InternalError, "update device token failed", err)
	}
	return nil
}
