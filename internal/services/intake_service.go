package services

import (
	"context"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type IntakeService interface {
	CreateIntake(ctx context.Context, userID string, req dto.CreateIntakeRequest) (dto.IntakeHistoryResponse, error)
	ListHistory(ctx context.Context, userID string, from, to string) ([]dto.IntakeHistoryResponse, error)
}

type intakeService struct{}

func NewIntakeService() IntakeService {
	return &intakeService{}
}

func (s *intakeService) CreateIntake(ctx context.Context, userID string, req dto.CreateIntakeRequest) (dto.IntakeHistoryResponse, error) {
	return dto.IntakeHistoryResponse{}, domain.NewError(constants.InternalNotImplemented, "intake not implemented")
}

func (s *intakeService) ListHistory(ctx context.Context, userID string, from, to string) ([]dto.IntakeHistoryResponse, error) {
	return nil, domain.NewError(constants.InternalNotImplemented, "intake history not implemented")
}
