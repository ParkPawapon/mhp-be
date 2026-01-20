package services

import (
	"context"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type HealthService interface {
	CreateHealthRecord(ctx context.Context, userID string, req dto.CreateHealthRecordRequest) (dto.HealthRecordResponse, error)
	ListHealthRecords(ctx context.Context, userID, from, to string) ([]dto.HealthRecordResponse, error)
	CreateDailyAssessment(ctx context.Context, userID string, req dto.CreateDailyAssessmentRequest) (dto.DailyAssessmentResponse, error)
	ListDailyAssessments(ctx context.Context, userID, from, to string) ([]dto.DailyAssessmentResponse, error)
}

type healthService struct{}

func NewHealthService() HealthService {
	return &healthService{}
}

func (s *healthService) CreateHealthRecord(ctx context.Context, userID string, req dto.CreateHealthRecordRequest) (dto.HealthRecordResponse, error) {
	return dto.HealthRecordResponse{}, domain.NewError(constants.InternalNotImplemented, "health records not implemented")
}

func (s *healthService) ListHealthRecords(ctx context.Context, userID, from, to string) ([]dto.HealthRecordResponse, error) {
	return nil, domain.NewError(constants.InternalNotImplemented, "health records not implemented")
}

func (s *healthService) CreateDailyAssessment(ctx context.Context, userID string, req dto.CreateDailyAssessmentRequest) (dto.DailyAssessmentResponse, error) {
	return dto.DailyAssessmentResponse{}, domain.NewError(constants.InternalNotImplemented, "daily assessments not implemented")
}

func (s *healthService) ListDailyAssessments(ctx context.Context, userID, from, to string) ([]dto.DailyAssessmentResponse, error) {
	return nil, domain.NewError(constants.InternalNotImplemented, "daily assessments not implemented")
}
