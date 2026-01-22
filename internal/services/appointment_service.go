package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type AppointmentService interface {
	ListAppointments(ctx context.Context, userID string) ([]dto.AppointmentResponse, error)
	CreateAppointment(ctx context.Context, userID string, req dto.CreateAppointmentRequest) (dto.AppointmentResponse, error)
	UpdateStatus(ctx context.Context, id string, req dto.UpdateAppointmentStatusRequest) error
	DeleteAppointment(ctx context.Context, id string) error
	CreateNurseVisitNote(ctx context.Context, appointmentID, nurseID string, req dto.CreateNurseVisitNoteRequest) error
	ListVisitHistory(ctx context.Context, userID string) ([]dto.VisitHistoryItem, error)
}

type appointmentService struct {
	repo   repositories.AppointmentRepository
	notify NotificationService
}

func NewAppointmentService(repo repositories.AppointmentRepository, notify NotificationService) AppointmentService {
	return &appointmentService{repo: repo, notify: notify}
}

func (s *appointmentService) ListAppointments(ctx context.Context, userID string) ([]dto.AppointmentResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	items, err := s.repo.ListAppointments(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.AppointmentResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.AppointmentResponse{
			ID:           item.ID.String(),
			UserID:       item.UserID.String(),
			Title:        item.Title,
			ApptType:     item.ApptType,
			ApptDateTime: item.ApptDateTime,
			LocationName: item.LocationName,
			SlipImageURL: item.SlipImageURL,
			Status:       item.Status,
			CreatedAt:    item.CreatedAt,
		})
	}
	return resp, nil
}

func (s *appointmentService) CreateAppointment(ctx context.Context, userID string, req dto.CreateAppointmentRequest) (dto.AppointmentResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return dto.AppointmentResponse{}, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	apptTime, err := parseRFC3339(req.ApptDateTime)
	if err != nil {
		return dto.AppointmentResponse{}, domain.NewError(constants.ValidationFailed, "invalid appt_datetime")
	}

	appt := &db.Appointment{
		UserID:       uid,
		Title:        req.Title,
		ApptType:     req.ApptType,
		ApptDateTime: apptTime.UTC(),
		LocationName: req.LocationName,
		SlipImageURL: req.SlipImageURL,
		Status:       constants.ApptPending,
	}

	if err := s.repo.CreateAppointment(ctx, appt); err != nil {
		return dto.AppointmentResponse{}, err
	}

	if s.notify != nil {
		_ = s.notify.ScheduleAppointmentReminders(ctx, appt)
	}

	return dto.AppointmentResponse{
		ID:           appt.ID.String(),
		UserID:       appt.UserID.String(),
		Title:        appt.Title,
		ApptType:     appt.ApptType,
		ApptDateTime: appt.ApptDateTime,
		LocationName: appt.LocationName,
		SlipImageURL: appt.SlipImageURL,
		Status:       appt.Status,
		CreatedAt:    appt.CreatedAt,
	}, nil
}

func (s *appointmentService) UpdateStatus(ctx context.Context, id string, req dto.UpdateAppointmentStatusRequest) error {
	apptID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid id")
	}

	appt, err := s.repo.FindByID(ctx, apptID)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateStatus(ctx, apptID, req.Status); err != nil {
		return err
	}

	if req.Status == constants.ApptCancelled && s.notify != nil {
		_ = s.notify.CancelAppointmentReminders(ctx, appt.UserID, appt.ID)
	}
	return nil
}

func (s *appointmentService) DeleteAppointment(ctx context.Context, id string) error {
	apptID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid id")
	}

	appt, err := s.repo.FindByID(ctx, apptID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteAppointment(ctx, apptID); err != nil {
		return err
	}
	if s.notify != nil {
		_ = s.notify.CancelAppointmentReminders(ctx, appt.UserID, appt.ID)
	}
	return nil
}

func (s *appointmentService) CreateNurseVisitNote(ctx context.Context, appointmentID, nurseID string, req dto.CreateNurseVisitNoteRequest) error {
	apptID, err := uuid.Parse(appointmentID)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid appointment_id")
	}

	nID, err := uuid.Parse(nurseID)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid nurse_id")
	}

	if _, err := s.repo.FindByID(ctx, apptID); err != nil {
		return err
	}

	var summary datatypes.JSON
	if req.VitalSignsSummary != nil {
		payload, err := json.Marshal(req.VitalSignsSummary)
		if err != nil {
			return domain.NewError(constants.ValidationFailed, "invalid vital_signs_summary")
		}
		summary = payload
	}

	note := &db.NurseVisitNote{
		AppointmentID:     apptID,
		NurseID:           nID,
		VisitDetails:      req.VisitDetails,
		VitalSignsSummary: summary,
		NextActionPlan:    req.NextActionPlan,
	}

	return s.repo.CreateNurseVisitNote(ctx, note)
}

func (s *appointmentService) ListVisitHistory(ctx context.Context, userID string) ([]dto.VisitHistoryItem, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	rows, err := s.repo.ListVisitHistory(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.VisitHistoryItem, 0, len(rows))
	for _, row := range rows {
		var summary any
		if len(row.VitalSignsSummary) > 0 {
			_ = json.Unmarshal(row.VitalSignsSummary, &summary)
		}
		resp = append(resp, dto.VisitHistoryItem{
			AppointmentID:     row.AppointmentID.String(),
			VisitNoteID:       row.VisitNoteID.String(),
			ApptDateTime:      row.ApptDateTime,
			Title:             row.Title,
			LocationName:      row.LocationName,
			NurseID:           row.NurseID.String(),
			VisitDetails:      row.VisitDetails,
			VitalSignsSummary: summary,
			NextActionPlan:    row.NextActionPlan,
			CreatedAt:         row.CreatedAt,
		})
	}
	return resp, nil
}

func parseRFC3339(value string) (time.Time, error) {
	if ts, err := time.Parse(time.RFC3339, value); err == nil {
		return ts, nil
	}
	return time.Parse(time.RFC3339Nano, value)
}
