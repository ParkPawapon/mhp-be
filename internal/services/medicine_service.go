package services

import (
	"context"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type MedicineService interface {
	ListMaster(ctx context.Context, page, pageSize int) ([]dto.MedicineMasterResponse, int64, error)
	CreatePatientMedicine(ctx context.Context, userID string, req dto.CreatePatientMedicineRequest) (dto.PatientMedicineResponse, error)
	ListPatientMedicines(ctx context.Context, userID string) ([]dto.PatientMedicineResponse, error)
	UpdatePatientMedicine(ctx context.Context, id string, req dto.UpdatePatientMedicineRequest) error
	DeletePatientMedicine(ctx context.Context, id string) error
	CreateSchedule(ctx context.Context, patientMedicineID string, req dto.CreateMedicineScheduleRequest) (dto.MedicineScheduleResponse, error)
	DeleteSchedule(ctx context.Context, id string) error
}

type medicineService struct{}

func NewMedicineService() MedicineService {
	return &medicineService{}
}

func (s *medicineService) ListMaster(ctx context.Context, page, pageSize int) ([]dto.MedicineMasterResponse, int64, error) {
	return nil, 0, domain.NewError(constants.InternalNotImplemented, "medicine master not implemented")
}

func (s *medicineService) CreatePatientMedicine(ctx context.Context, userID string, req dto.CreatePatientMedicineRequest) (dto.PatientMedicineResponse, error) {
	return dto.PatientMedicineResponse{}, domain.NewError(constants.InternalNotImplemented, "patient medicine not implemented")
}

func (s *medicineService) ListPatientMedicines(ctx context.Context, userID string) ([]dto.PatientMedicineResponse, error) {
	return nil, domain.NewError(constants.InternalNotImplemented, "patient medicine not implemented")
}

func (s *medicineService) UpdatePatientMedicine(ctx context.Context, id string, req dto.UpdatePatientMedicineRequest) error {
	return domain.NewError(constants.InternalNotImplemented, "patient medicine not implemented")
}

func (s *medicineService) DeletePatientMedicine(ctx context.Context, id string) error {
	return domain.NewError(constants.InternalNotImplemented, "patient medicine not implemented")
}

func (s *medicineService) CreateSchedule(ctx context.Context, patientMedicineID string, req dto.CreateMedicineScheduleRequest) (dto.MedicineScheduleResponse, error) {
	return dto.MedicineScheduleResponse{}, domain.NewError(constants.InternalNotImplemented, "medicine schedule not implemented")
}

func (s *medicineService) DeleteSchedule(ctx context.Context, id string) error {
	return domain.NewError(constants.InternalNotImplemented, "medicine schedule not implemented")
}
