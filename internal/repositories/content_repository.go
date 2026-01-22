package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type ContentRepository interface {
	ListHealthContent(ctx context.Context, publishedOnly bool) ([]db.HealthContent, error)
	CreateHealthContent(ctx context.Context, content *db.HealthContent) error
	UpdateHealthContent(ctx context.Context, id uuid.UUID, updates map[string]any) error
	SetPublished(ctx context.Context, id uuid.UUID, published bool) error
	FindByID(ctx context.Context, id uuid.UUID) (*db.HealthContent, error)
}

type contentRepository struct {
	db *gorm.DB
}

func NewContentRepository(dbConn *gorm.DB) ContentRepository {
	return &contentRepository{db: dbConn}
}

func (r *contentRepository) ListHealthContent(ctx context.Context, publishedOnly bool) ([]db.HealthContent, error) {
	var items []db.HealthContent
	query := r.db.WithContext(ctx)
	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}
	if err := query.Order("created_at desc").Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list health content failed", err)
	}
	return items, nil
}

func (r *contentRepository) CreateHealthContent(ctx context.Context, content *db.HealthContent) error {
	if err := r.db.WithContext(ctx).Create(content).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create health content failed", err)
	}
	return nil
}

func (r *contentRepository) UpdateHealthContent(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&db.HealthContent{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return domain.WrapError(constants.InternalError, "update health content failed", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.NewError(constants.ContentNotFound, "health content not found")
	}
	return nil
}

func (r *contentRepository) SetPublished(ctx context.Context, id uuid.UUID, published bool) error {
	result := r.db.WithContext(ctx).Model(&db.HealthContent{}).Where("id = ?", id).Update("is_published", published)
	if result.Error != nil {
		return domain.WrapError(constants.InternalError, "update publish status failed", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.NewError(constants.ContentNotFound, "health content not found")
	}
	return nil
}

func (r *contentRepository) FindByID(ctx context.Context, id uuid.UUID) (*db.HealthContent, error) {
	var item db.HealthContent
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.ContentNotFound, "health content not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find health content failed", err)
	}
	return &item, nil
}
