package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type appointmentRepoStub struct {
	appointment *db.Appointment
}

func (s *appointmentRepoStub) ListAppointments(ctx context.Context, userID uuid.UUID) ([]db.Appointment, error) {
	panic("not used")
}
func (s *appointmentRepoStub) CreateAppointment(ctx context.Context, appt *db.Appointment) error {
	appt.ID = uuid.New()
	appt.CreatedAt = time.Now().UTC()
	s.appointment = appt
	return nil
}
func (s *appointmentRepoStub) FindByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error) {
	if s.appointment == nil {
		return nil, domain.NewError(constants.ApptNotFound, "appointment not found")
	}
	return s.appointment, nil
}
func (s *appointmentRepoStub) UpdateStatus(ctx context.Context, id uuid.UUID, status constants.AppointmentStatus) error {
	if s.appointment == nil {
		return domain.NewError(constants.ApptNotFound, "appointment not found")
	}
	s.appointment.Status = status
	return nil
}
func (s *appointmentRepoStub) DeleteAppointment(ctx context.Context, id uuid.UUID) error {
	if s.appointment == nil {
		return domain.NewError(constants.ApptNotFound, "appointment not found")
	}
	return nil
}
func (s *appointmentRepoStub) CreateNurseVisitNote(ctx context.Context, note *db.NurseVisitNote) error {
	return nil
}
func (s *appointmentRepoStub) ListVisitHistory(ctx context.Context, userID uuid.UUID) ([]repositories.VisitHistoryRow, error) {
	panic("not used")
}

type notificationCancelStub struct {
	cancelled bool
}

func (s *notificationCancelStub) ScheduleMedicineReminders(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, mealTiming *string, timeSlot time.Time) error {
	panic("not used")
}
func (s *notificationCancelStub) CancelMedicineAfterMealReminder(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate time.Time) error {
	panic("not used")
}
func (s *notificationCancelStub) ScheduleAppointmentReminders(ctx context.Context, appt *db.Appointment) error {
	return nil
}
func (s *notificationCancelStub) CancelAppointmentReminders(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	s.cancelled = true
	return nil
}
func (s *notificationCancelStub) ListUpcoming(ctx context.Context, userID string, from, to string) ([]dto.NotificationUpcomingItem, error) {
	panic("not used")
}
func (s *notificationCancelStub) EnsureWeeklyReminders(ctx context.Context) error {
	panic("not used")
}
func (s *notificationCancelStub) ProcessDue(ctx context.Context) error {
	panic("not used")
}
func (s *notificationCancelStub) CancelWeeklyReminders(ctx context.Context, userID uuid.UUID) error {
	panic("not used")
}

func TestCreateAppointmentValidation(t *testing.T) {
	repo := &appointmentRepoStub{}
	svc := NewAppointmentService(repo, nil)

	_, err := svc.CreateAppointment(context.Background(), uuid.New().String(), dto.CreateAppointmentRequest{
		Title:        "Visit",
		ApptType:     constants.ApptHospital,
		ApptDateTime: "bad-time",
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestUpdateStatusCancelsWhenCancelled(t *testing.T) {
	repo := &appointmentRepoStub{appointment: &db.Appointment{ID: uuid.New(), UserID: uuid.New(), ApptType: constants.ApptHospital}}
	notify := &notificationCancelStub{}
	svc := NewAppointmentService(repo, notify)

	if err := svc.UpdateStatus(context.Background(), repo.appointment.ID.String(), dto.UpdateAppointmentStatusRequest{Status: constants.ApptCancelled}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !notify.cancelled {
		t.Fatalf("expected reminders cancelled")
	}
}

func TestDeleteAppointmentCancels(t *testing.T) {
	repo := &appointmentRepoStub{appointment: &db.Appointment{ID: uuid.New(), UserID: uuid.New(), ApptType: constants.ApptHospital}}
	notify := &notificationCancelStub{}
	svc := NewAppointmentService(repo, notify)

	if err := svc.DeleteAppointment(context.Background(), repo.appointment.ID.String()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !notify.cancelled {
		t.Fatalf("expected reminders cancelled")
	}
}
