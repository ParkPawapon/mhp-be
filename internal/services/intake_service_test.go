package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type fakeIntakeRepo struct {
	created *db.IntakeHistory
}

func (f *fakeIntakeRepo) Create(ctx context.Context, intake *db.IntakeHistory) error {
	f.created = intake
	return nil
}

func (f *fakeIntakeRepo) ListHistory(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]db.IntakeHistory, error) {
	return nil, nil
}

type fakeNotificationService struct {
	cancelCalled bool
	gotSchedule  uuid.UUID
	gotDate      time.Time
}

func (f *fakeNotificationService) ScheduleMedicineReminders(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, mealTiming *string, timeSlot time.Time) error {
	return nil
}

func (f *fakeNotificationService) CancelMedicineAfterMealReminder(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate time.Time) error {
	f.cancelCalled = true
	f.gotSchedule = scheduleID
	f.gotDate = targetDate
	return nil
}

func (f *fakeNotificationService) ScheduleAppointmentReminders(ctx context.Context, appt *db.Appointment) error {
	return nil
}

func (f *fakeNotificationService) CancelAppointmentReminders(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	return nil
}

func (f *fakeNotificationService) ListUpcoming(ctx context.Context, userID string, from, to string) ([]dto.NotificationUpcomingItem, error) {
	return nil, nil
}

func (f *fakeNotificationService) EnsureWeeklyReminders(ctx context.Context) error {
	return nil
}

func (f *fakeNotificationService) ProcessDue(ctx context.Context) error {
	return nil
}

func (f *fakeNotificationService) CancelWeeklyReminders(ctx context.Context, userID uuid.UUID) error {
	return nil
}

var _ repositories.IntakeRepository = (*fakeIntakeRepo)(nil)

func TestCreateIntakeCancelsAfterMealReminderWhenTaken(t *testing.T) {
	repo := &fakeIntakeRepo{}
	notify := &fakeNotificationService{}
	svc := NewIntakeService(repo, notify)

	scheduleID := uuid.New().String()
	userID := uuid.New().String()
	req := dto.CreateIntakeRequest{
		ScheduleID: &scheduleID,
		TargetDate: "2026-01-20",
		Status:     constants.MedTaken,
	}

	_, err := svc.CreateIntake(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.created == nil || repo.created.TakenAt == nil {
		t.Fatalf("expected intake record with taken_at")
	}
	if !notify.cancelCalled {
		t.Fatalf("expected cancel reminder to be called")
	}
	if notify.gotSchedule.String() != scheduleID {
		t.Fatalf("expected schedule id %s, got %s", scheduleID, notify.gotSchedule.String())
	}
	if notify.gotDate.Format("2006-01-02") != "2026-01-20" {
		t.Fatalf("expected target date 2026-01-20, got %s", notify.gotDate.Format("2006-01-02"))
	}
}
