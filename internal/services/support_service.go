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

type SupportService interface {
	CreateChatRequest(ctx context.Context, userID string, req dto.SupportChatRequestCreateRequest) (dto.SupportChatRequestResponse, error)
	ListChatRequests(ctx context.Context, page, pageSize int) ([]dto.SupportChatRequestItem, int64, error)
	GetEmergencyInfo(ctx context.Context) dto.SupportEmergencyResponse
}

type supportService struct {
	repo repositories.SupportRepository
}

func NewSupportService(repo repositories.SupportRepository) SupportService {
	return &supportService{repo: repo}
}

func (s *supportService) CreateChatRequest(ctx context.Context, userID string, req dto.SupportChatRequestCreateRequest) (dto.SupportChatRequestResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return dto.SupportChatRequestResponse{}, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}
	category := strings.ToUpper(strings.TrimSpace(req.Category))
	if !isAllowed(category, constants.SupportCategories) {
		return dto.SupportChatRequestResponse{}, domain.NewError(constants.ValidationFailed, "invalid category")
	}

	message := strings.TrimSpace(req.Message)
	if message == "" {
		return dto.SupportChatRequestResponse{}, domain.NewError(constants.ValidationFailed, "message required")
	}

	item := &db.SupportChatRequest{
		UserID:        uid,
		Message:       message,
		Category:      category,
		AttachmentURL: req.AttachmentURL,
		Status:        "OPEN",
	}
	if err := s.repo.CreateChatRequest(ctx, item); err != nil {
		return dto.SupportChatRequestResponse{}, err
	}

	return dto.SupportChatRequestResponse{ID: item.ID.String(), Status: item.Status}, nil
}

func (s *supportService) ListChatRequests(ctx context.Context, page, pageSize int) ([]dto.SupportChatRequestItem, int64, error) {
	items, total, err := s.repo.ListChatRequests(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]dto.SupportChatRequestItem, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.SupportChatRequestItem{
			ID:            item.ID.String(),
			UserID:        item.UserID.String(),
			Message:       item.Message,
			Category:      item.Category,
			AttachmentURL: item.AttachmentURL,
			Status:        item.Status,
			CreatedAt:     item.CreatedAt,
		})
	}
	return resp, total, nil
}

func (s *supportService) GetEmergencyInfo(ctx context.Context) dto.SupportEmergencyResponse {
	return dto.SupportEmergencyResponse{Hotline: "1669", DisplayName: "Emergency 1669"}
}
