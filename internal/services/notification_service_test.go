package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type fakeNotificationRepo struct {
	created []db.NotificationEvent
}

func (f *fakeNotificationRepo) WithTx(tx *gorm.DB) repositories.NotificationRepository {
	return f
}

func (f *fakeNotificationRepo) CreateEvents(ctx context.Context, events []db.NotificationEvent) error {
	f.created = append(f.created, events...)
	return nil
}

func (f *fakeNotificationRepo) ListUpcoming(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]db.NotificationEvent, error) {
	return nil, nil
}

func (f *fakeNotificationRepo) ListDueForUpdate(ctx context.Context, now time.Time, limit int) ([]db.NotificationEvent, error) {
	return nil, nil
}

func (f *fakeNotificationRepo) UpdateEventStatus(ctx context.Context, id uuid.UUID, status constants.NotificationStatus, sentAt *time.Time) error {
	return nil
}

func (f *fakeNotificationRepo) FindTemplateByCode(ctx context.Context, code string) (*db.NotificationTemplate, error) {
	return &db.NotificationTemplate{}, nil
}

func (f *fakeNotificationRepo) CancelPendingBySchedule(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate string) error {
	return nil
}

func (f *fakeNotificationRepo) CancelPendingByAppointment(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	return nil
}

func (f *fakeNotificationRepo) CancelPendingByTemplate(ctx context.Context, userID uuid.UUID, templateCode string) error {
	return nil
}

func TestScheduleAppointmentRemindersCreatesTwoEvents(t *testing.T) {
	repo := &fakeNotificationRepo{}
	cfg := config.NotificationConfig{ScheduleDays: 1, Timezone: "UTC"}
	svc := NewNotificationService(cfg, nil, repo, nil, nil, zap.NewNop())

	impl, ok := svc.(*notificationService)
	if !ok {
		t.Fatalf("expected notificationService")
	}

	fixedNow := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	impl.now = func() time.Time { return fixedNow }

	appt := &db.Appointment{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		ApptType:     constants.ApptHospital,
		ApptDateTime: fixedNow.AddDate(0, 0, 10),
	}

	if err := impl.ScheduleAppointmentReminders(context.Background(), appt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.created) != 2 {
		t.Fatalf("expected 2 events, got %d", len(repo.created))
	}

	codes := map[string]bool{}
	for _, event := range repo.created {
		codes[event.TemplateCode] = true
	}

	if !codes[constants.TemplateAppt5Days] || !codes[constants.TemplateAppt1Day] {
		t.Fatalf("expected appointment template codes")
	}
}

func TestScheduleMedicineRemindersBeforeMealCreatesTwoEvents(t *testing.T) {
	repo := &fakeNotificationRepo{}
	cfg := config.NotificationConfig{ScheduleDays: 1, Timezone: "UTC"}
	svc := NewNotificationService(cfg, nil, repo, nil, nil, zap.NewNop())

	impl, ok := svc.(*notificationService)
	if !ok {
		t.Fatalf("expected notificationService")
	}

	fixedNow := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	impl.now = func() time.Time { return fixedNow }

	timeSlot := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
	mealTiming := constants.MealTimingBeforeMeal

	if err := impl.ScheduleMedicineReminders(context.Background(), uuid.New(), uuid.New(), &mealTiming, timeSlot); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.created) != 2 {
		t.Fatalf("expected 2 events, got %d", len(repo.created))
	}

	codes := map[string]bool{}
	for _, event := range repo.created {
		codes[event.TemplateCode] = true
	}

	if !codes[constants.TemplateMedBeforeMeal5Min] || !codes[constants.TemplateMedBeforeMeal20Min] {
		t.Fatalf("expected medicine template codes")
	}
}
