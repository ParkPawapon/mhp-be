package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type ProfileRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*db.UserProfile, error)
	Upsert(ctx context.Context, profile *db.UserProfile) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(dbConn *gorm.DB) ProfileRepository {
	return &profileRepository{db: dbConn}
}

func (r *profileRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*db.UserProfile, error) {
	var profile db.UserProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.UserNotFound, "profile not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find profile failed", err)
	}
	return &profile, nil
}

func (r *profileRepository) Upsert(ctx context.Context, profile *db.UserProfile) error {
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(profile).Error; err != nil {
		return domain.WrapError(constants.InternalError, "upsert profile failed", err)
	}
	return nil
}
