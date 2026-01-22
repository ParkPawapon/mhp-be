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

type PreferenceRepository interface {
	Upsert(ctx context.Context, pref *db.UserPreference) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*db.UserPreference, error)
	ListWeeklyReminderUsers(ctx context.Context) ([]uuid.UUID, error)
}

type preferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(dbConn *gorm.DB) PreferenceRepository {
	return &preferenceRepository{db: dbConn}
}

func (r *preferenceRepository) Upsert(ctx context.Context, pref *db.UserPreference) error {
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(pref).Error; err != nil {
		return domain.WrapError(constants.InternalError, "upsert preferences failed", err)
	}
	return nil
}

func (r *preferenceRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*db.UserPreference, error) {
	var pref db.UserPreference
	if err := r.db.WithContext(ctx).First(&pref, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.UserNotFound, "preferences not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find preferences failed", err)
	}
	return &pref, nil
}

func (r *preferenceRepository) ListWeeklyReminderUsers(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id").
		Joins("LEFT JOIN user_preferences ON user_preferences.user_id = users.id").
		Where("users.role = ? AND users.is_active = ? AND (user_preferences.weekly_reminder_enabled IS NULL OR user_preferences.weekly_reminder_enabled = ?)", constants.RolePatient, true, true).
		Scan(&ids).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list reminder users failed", err)
	}
	return ids, nil
}
