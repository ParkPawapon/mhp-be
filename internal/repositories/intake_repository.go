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

type IntakeRepository interface {
	Create(ctx context.Context, intake *db.IntakeHistory) error
	ListHistory(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]db.IntakeHistory, error)
}

type intakeRepository struct {
	db *gorm.DB
}

func NewIntakeRepository(dbConn *gorm.DB) IntakeRepository {
	return &intakeRepository{db: dbConn}
}

func (r *intakeRepository) Create(ctx context.Context, intake *db.IntakeHistory) error {
	if err := r.db.WithContext(ctx).Create(intake).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create intake failed", err)
	}
	return nil
}

func (r *intakeRepository) ListHistory(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]db.IntakeHistory, error) {
	var items []db.IntakeHistory
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if !from.IsZero() {
		query = query.Where("target_date >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("target_date <= ?", to)
	}
	if err := query.Order("target_date desc, created_at desc").Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list intake history failed", err)
	}
	return items, nil
}
