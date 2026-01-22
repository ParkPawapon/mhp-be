package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type ContentService interface {
	ListHealthContent(ctx context.Context, publishedOnly bool) ([]dto.HealthContentResponse, error)
	CreateHealthContent(ctx context.Context, req dto.CreateHealthContentRequest) (dto.HealthContentResponse, error)
	UpdateHealthContent(ctx context.Context, id string, req dto.UpdateHealthContentRequest) error
	PublishHealthContent(ctx context.Context, id string, req dto.PublishHealthContentRequest) error
	ListHealthCategories(ctx context.Context) []string
}

type contentService struct {
	repo repositories.ContentRepository
}

func NewContentService(repo repositories.ContentRepository) ContentService {
	return &contentService{repo: repo}
}

func (s *contentService) ListHealthContent(ctx context.Context, publishedOnly bool) ([]dto.HealthContentResponse, error) {
	items, err := s.repo.ListHealthContent(ctx, publishedOnly)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.HealthContentResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.HealthContentResponse{
			ID:               item.ID.String(),
			Title:            item.Title,
			BodyContent:      item.BodyContent,
			ThumbnailURL:     item.ThumbnailURL,
			ExternalVideoURL: item.ExternalVideoURL,
			Category:         item.Category,
			IsPublished:      item.IsPublished,
			CreatedAt:        item.CreatedAt,
		})
	}
	return resp, nil
}

func (s *contentService) CreateHealthContent(ctx context.Context, req dto.CreateHealthContentRequest) (dto.HealthContentResponse, error) {
	if req.Category != nil {
		if !isAllowed(*req.Category, constants.HealthContentCategories) {
			return dto.HealthContentResponse{}, domain.NewError(constants.ContentInvalid, "invalid category")
		}
	}

	content := &db.HealthContent{
		Title:            strings.TrimSpace(req.Title),
		BodyContent:      req.BodyContent,
		ThumbnailURL:     req.ThumbnailURL,
		ExternalVideoURL: req.ExternalVideoURL,
		Category:         req.Category,
		IsPublished:      false,
	}

	if err := s.repo.CreateHealthContent(ctx, content); err != nil {
		return dto.HealthContentResponse{}, err
	}

	return dto.HealthContentResponse{
		ID:               content.ID.String(),
		Title:            content.Title,
		BodyContent:      content.BodyContent,
		ThumbnailURL:     content.ThumbnailURL,
		ExternalVideoURL: content.ExternalVideoURL,
		Category:         content.Category,
		IsPublished:      content.IsPublished,
		CreatedAt:        content.CreatedAt,
	}, nil
}

func (s *contentService) UpdateHealthContent(ctx context.Context, id string, req dto.UpdateHealthContentRequest) error {
	contentID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid id")
	}

	updates := map[string]any{}
	if req.Title != nil {
		value := strings.TrimSpace(*req.Title)
		if value == "" {
			return domain.NewError(constants.ValidationFailed, "title required")
		}
		updates["title"] = value
	}
	if req.BodyContent != nil {
		updates["body_content"] = req.BodyContent
	}
	if req.ThumbnailURL != nil {
		updates["thumbnail_url"] = req.ThumbnailURL
	}
	if req.ExternalVideoURL != nil {
		updates["external_video_url"] = req.ExternalVideoURL
	}
	if req.Category != nil {
		if !isAllowed(*req.Category, constants.HealthContentCategories) {
			return domain.NewError(constants.ContentInvalid, "invalid category")
		}
		updates["category"] = req.Category
	}
	if len(updates) == 0 {
		return domain.NewError(constants.ValidationFailed, "no fields to update")
	}

	return s.repo.UpdateHealthContent(ctx, contentID, updates)
}

func (s *contentService) PublishHealthContent(ctx context.Context, id string, req dto.PublishHealthContentRequest) error {
	contentID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid id")
	}
	return s.repo.SetPublished(ctx, contentID, req.IsPublished)
}

func (s *contentService) ListHealthCategories(ctx context.Context) []string {
	return constants.HealthContentCategories
}
