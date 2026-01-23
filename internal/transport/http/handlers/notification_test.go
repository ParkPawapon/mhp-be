package handlers

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type notificationServiceStub struct{}

func (notificationServiceStub) ListUpcoming(ctx context.Context, userID string, from, to string) ([]dto.NotificationUpcomingItem, error) {
	return []dto.NotificationUpcomingItem{{ID: uuid.New().String(), TemplateCode: "T", ScheduledAt: time.Now().UTC(), Status: "PENDING"}}, nil
}
func (notificationServiceStub) ScheduleMedicineReminders(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, mealTiming *string, timeSlot time.Time) error {
	panic("not used")
}
func (notificationServiceStub) CancelMedicineAfterMealReminder(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate time.Time) error {
	panic("not used")
}
func (notificationServiceStub) ScheduleAppointmentReminders(ctx context.Context, appt *db.Appointment) error {
	panic("not used")
}
func (notificationServiceStub) CancelAppointmentReminders(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	panic("not used")
}
func (notificationServiceStub) EnsureWeeklyReminders(ctx context.Context) error {
	panic("not used")
}
func (notificationServiceStub) ProcessDue(ctx context.Context) error {
	panic("not used")
}
func (notificationServiceStub) CancelWeeklyReminders(ctx context.Context, userID uuid.UUID) error {
	panic("not used")
}

func TestNotificationHandlers(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RolePatient, actorID))
	handler := NewNotificationHandler(notificationServiceStub{})

	router.GET("/notifications/upcoming", handler.ListUpcoming)

	resp := performRequest(router, http.MethodGet, "/notifications/upcoming?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
