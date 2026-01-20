package services

import (
	"context"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type AppointmentService interface {
	ListAppointments(ctx context.Context, userID string) ([]dto.AppointmentResponse, error)
	CreateAppointment(ctx context.Context, userID string, req dto.CreateAppointmentRequest) (dto.AppointmentResponse, error)
	UpdateStatus(ctx context.Context, id string, req dto.UpdateAppointmentStatusRequest) error
	DeleteAppointment(ctx context.Context, id string) error
	CreateNurseVisitNote(ctx context.Context, appointmentID, nurseID string, req dto.CreateNurseVisitNoteRequest) error
}

type appointmentService struct{}

func NewAppointmentService() AppointmentService {
	return &appointmentService{}
}

func (s *appointmentService) ListAppointments(ctx context.Context, userID string) ([]dto.AppointmentResponse, error) {
	return nil, domain.NewError(constants.InternalNotImplemented, "appointments not implemented")
}

func (s *appointmentService) CreateAppointment(ctx context.Context, userID string, req dto.CreateAppointmentRequest) (dto.AppointmentResponse, error) {
	return dto.AppointmentResponse{}, domain.NewError(constants.InternalNotImplemented, "appointments not implemented")
}

func (s *appointmentService) UpdateStatus(ctx context.Context, id string, req dto.UpdateAppointmentStatusRequest) error {
	return domain.NewError(constants.InternalNotImplemented, "appointments not implemented")
}

func (s *appointmentService) DeleteAppointment(ctx context.Context, id string) error {
	return domain.NewError(constants.InternalNotImplemented, "appointments not implemented")
}

func (s *appointmentService) CreateNurseVisitNote(ctx context.Context, appointmentID, nurseID string, req dto.CreateNurseVisitNoteRequest) error {
	return domain.NewError(constants.InternalNotImplemented, "appointments not implemented")
}
