package services

import (
	"context"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type ContentService interface {
	ListHealthContent(ctx context.Context, publishedOnly bool) ([]dto.HealthContentResponse, error)
	CreateHealthContent(ctx context.Context, req dto.CreateHealthContentRequest) (dto.HealthContentResponse, error)
	UpdateHealthContent(ctx context.Context, id string, req dto.UpdateHealthContentRequest) error
	PublishHealthContent(ctx context.Context, id string, req dto.PublishHealthContentRequest) error
}

type contentService struct{}

func NewContentService() ContentService {
	return &contentService{}
}

func (s *contentService) ListHealthContent(ctx context.Context, publishedOnly bool) ([]dto.HealthContentResponse, error) {
	return nil, domain.NewError(constants.InternalNotImplemented, "content not implemented")
}

func (s *contentService) CreateHealthContent(ctx context.Context, req dto.CreateHealthContentRequest) (dto.HealthContentResponse, error) {
	return dto.HealthContentResponse{}, domain.NewError(constants.InternalNotImplemented, "content not implemented")
}

func (s *contentService) UpdateHealthContent(ctx context.Context, id string, req dto.UpdateHealthContentRequest) error {
	return domain.NewError(constants.InternalNotImplemented, "content not implemented")
}

func (s *contentService) PublishHealthContent(ctx context.Context, id string, req dto.PublishHealthContentRequest) error {
	return domain.NewError(constants.InternalNotImplemented, "content not implemented")
}
