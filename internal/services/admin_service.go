package services

import (
	"context"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type AdminService interface {
	StaffLogin(ctx context.Context, req dto.StaffLoginRequest) (dto.TokenResponse, error)
	ListPatients(ctx context.Context, page, pageSize int) ([]dto.PatientSummaryResponse, int64, error)
	GetPatient(ctx context.Context, id string) (dto.PatientDetailResponse, error)
	ListAdherence(ctx context.Context, patientID, from, to string) ([]dto.IntakeHistoryResponse, error)
}

type adminService struct{}

func NewAdminService() AdminService {
	return &adminService{}
}

func (s *adminService) StaffLogin(ctx context.Context, req dto.StaffLoginRequest) (dto.TokenResponse, error) {
	return dto.TokenResponse{}, domain.NewError(constants.InternalNotImplemented, "staff login not implemented")
}

func (s *adminService) ListPatients(ctx context.Context, page, pageSize int) ([]dto.PatientSummaryResponse, int64, error) {
	return nil, 0, domain.NewError(constants.InternalNotImplemented, "admin patients not implemented")
}

func (s *adminService) GetPatient(ctx context.Context, id string) (dto.PatientDetailResponse, error) {
	return dto.PatientDetailResponse{}, domain.NewError(constants.InternalNotImplemented, "admin patients not implemented")
}

func (s *adminService) ListAdherence(ctx context.Context, patientID, from, to string) ([]dto.IntakeHistoryResponse, error) {
	return nil, domain.NewError(constants.InternalNotImplemented, "admin adherence not implemented")
}
