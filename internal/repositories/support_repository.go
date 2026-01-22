package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type SupportRepository interface {
	CreateChatRequest(ctx context.Context, req *db.SupportChatRequest) error
	ListChatRequests(ctx context.Context, page, pageSize int) ([]db.SupportChatRequest, int64, error)
}

type supportRepository struct {
	db *gorm.DB
}

func NewSupportRepository(dbConn *gorm.DB) SupportRepository {
	return &supportRepository{db: dbConn}
}

func (r *supportRepository) CreateChatRequest(ctx context.Context, req *db.SupportChatRequest) error {
	if err := r.db.WithContext(ctx).Create(req).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create support chat request failed", err)
	}
	return nil
}

func (r *supportRepository) ListChatRequests(ctx context.Context, page, pageSize int) ([]db.SupportChatRequest, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&db.SupportChatRequest{}).Count(&total).Error; err != nil {
		return nil, 0, domain.WrapError(constants.InternalError, "count support chat requests failed", err)
	}

	var items []db.SupportChatRequest
	if err := r.db.WithContext(ctx).
		Order("created_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, domain.WrapError(constants.InternalError, "list support chat requests failed", err)
	}
	return items, total, nil
}
