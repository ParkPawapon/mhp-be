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
)

type userRepoStub struct {
	findByID func(ctx context.Context, id uuid.UUID) (*db.User, error)
}

func (s userRepoStub) Create(ctx context.Context, user *db.User) error {
	panic("not used")
}
func (s userRepoStub) FindByUsername(ctx context.Context, username string) (*db.User, error) {
	panic("not used")
}
func (s userRepoStub) FindByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	return s.findByID(ctx, id)
}
func (s userRepoStub) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	panic("not used")
}

type profileRepoStub struct {
	findByUserID func(ctx context.Context, userID uuid.UUID) (*db.UserProfile, error)
	upsert       func(ctx context.Context, profile *db.UserProfile) error
}

func (s profileRepoStub) FindByUserID(ctx context.Context, userID uuid.UUID) (*db.UserProfile, error) {
	return s.findByUserID(ctx, userID)
}
func (s profileRepoStub) Upsert(ctx context.Context, profile *db.UserProfile) error {
	return s.upsert(ctx, profile)
}

type deviceTokenRepoStub struct {
	save func(ctx context.Context, token *db.DeviceToken) error
}

func (s deviceTokenRepoStub) Save(ctx context.Context, token *db.DeviceToken) error {
	return s.save(ctx, token)
}

type preferenceRepoStub struct {
	upsert func(ctx context.Context, pref *db.UserPreference) error
}

func (s preferenceRepoStub) Upsert(ctx context.Context, pref *db.UserPreference) error {
	return s.upsert(ctx, pref)
}
func (s preferenceRepoStub) FindByUserID(ctx context.Context, userID uuid.UUID) (*db.UserPreference, error) {
	panic("not used")
}
func (s preferenceRepoStub) ListWeeklyReminderUsers(ctx context.Context) ([]uuid.UUID, error) {
	panic("not used")
}

type notificationStub struct {
	cancelWeekly func(ctx context.Context, userID uuid.UUID) error
}

func (s notificationStub) ScheduleMedicineReminders(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, mealTiming *string, timeSlot time.Time) error {
	panic("not used")
}
func (s notificationStub) CancelMedicineAfterMealReminder(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate time.Time) error {
	panic("not used")
}
func (s notificationStub) ScheduleAppointmentReminders(ctx context.Context, appt *db.Appointment) error {
	panic("not used")
}
func (s notificationStub) CancelAppointmentReminders(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	panic("not used")
}
func (s notificationStub) ListUpcoming(ctx context.Context, userID string, from, to string) ([]dto.NotificationUpcomingItem, error) {
	panic("not used")
}
func (s notificationStub) EnsureWeeklyReminders(ctx context.Context) error {
	panic("not used")
}
func (s notificationStub) ProcessDue(ctx context.Context) error {
	panic("not used")
}
func (s notificationStub) CancelWeeklyReminders(ctx context.Context, userID uuid.UUID) error {
	if s.cancelWeekly != nil {
		return s.cancelWeekly(ctx, userID)
	}
	return nil
}

func TestUserServiceGetMeMasking(t *testing.T) {
	actorID := uuid.New()
	citizenID := "1234567890123"

	userRepo := userRepoStub{findByID: func(ctx context.Context, id uuid.UUID) (*db.User, error) {
		return &db.User{ID: actorID, Role: constants.RolePatient}, nil
	}}
	profileRepo := profileRepoStub{findByUserID: func(ctx context.Context, userID uuid.UUID) (*db.UserProfile, error) {
		return &db.UserProfile{
			UserID:      actorID,
			FirstName:   "A",
			LastName:    "B",
			CitizenID:   &citizenID,
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil
	}, upsert: func(ctx context.Context, profile *db.UserProfile) error { return nil }}

	svc := NewUserService(userRepo, profileRepo, nil, nil, nil)

	resp, err := svc.GetMe(context.Background(), actorID, constants.RolePatient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Profile.CitizenID == nil || *resp.Profile.CitizenID == citizenID {
		t.Fatalf("expected masked citizen id")
	}

	respAdmin, err := svc.GetMe(context.Background(), actorID, constants.RoleAdmin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if respAdmin.Profile.CitizenID == nil || *respAdmin.Profile.CitizenID != citizenID {
		t.Fatalf("expected unmasked citizen id")
	}

	respCaregiver, err := svc.GetMe(context.Background(), actorID, constants.RoleCaregiver)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if respCaregiver.Profile.CitizenID != nil {
		t.Fatalf("expected citizen id hidden for caregiver")
	}
}

func TestUserServiceUpdateProfileRequiresFieldsOnNew(t *testing.T) {
	actorID := uuid.New()
	userRepo := userRepoStub{findByID: func(ctx context.Context, id uuid.UUID) (*db.User, error) { return &db.User{ID: actorID}, nil }}

	profileRepo := profileRepoStub{
		findByUserID: func(ctx context.Context, userID uuid.UUID) (*db.UserProfile, error) {
			return nil, domain.NewError(constants.UserNotFound, "profile not found")
		},
		upsert: func(ctx context.Context, profile *db.UserProfile) error { return nil },
	}

	svc := NewUserService(userRepo, profileRepo, nil, nil, nil)

	if err := svc.UpdateProfile(context.Background(), actorID, dto.UpdateProfileRequest{}); err == nil {
		t.Fatalf("expected validation error for missing required fields")
	}

	first := "Jane"
	last := "Doe"
	dob := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := svc.UpdateProfile(context.Background(), actorID, dto.UpdateProfileRequest{
		FirstName:   &first,
		LastName:    &last,
		DateOfBirth: &dob,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserServiceSaveDeviceToken(t *testing.T) {
	actorID := uuid.New()
	called := false
	deviceRepo := deviceTokenRepoStub{save: func(ctx context.Context, token *db.DeviceToken) error {
		called = true
		if token.UserID != actorID {
			return domain.NewError(constants.InternalError, "bad user")
		}
		return nil
	}}

	svc := NewUserService(userRepoStub{findByID: func(ctx context.Context, id uuid.UUID) (*db.User, error) { return &db.User{}, nil }}, profileRepoStub{}, deviceRepo, preferenceRepoStub{}, nil)

	if err := svc.SaveDeviceToken(context.Background(), actorID, dto.DeviceTokenRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}

	req := dto.DeviceTokenRequest{Platform: "ios", DeviceToken: "token"}
	if err := svc.SaveDeviceToken(context.Background(), actorID, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected device token to be saved")
	}
}

func TestUserServiceUpdatePreferences(t *testing.T) {
	actorID := uuid.New()
	cancelled := false
	prefRepo := preferenceRepoStub{upsert: func(ctx context.Context, pref *db.UserPreference) error { return nil }}
	notify := notificationStub{cancelWeekly: func(ctx context.Context, userID uuid.UUID) error {
		cancelled = true
		return nil
	}}

	svc := NewUserService(userRepoStub{findByID: func(ctx context.Context, id uuid.UUID) (*db.User, error) { return &db.User{}, nil }}, profileRepoStub{}, deviceTokenRepoStub{}, prefRepo, notify)

	if _, err := svc.UpdatePreferences(context.Background(), actorID, dto.UpdatePreferencesRequest{}); err == nil {
		t.Fatalf("expected validation error")
	}

	flag := false
	resp, err := svc.UpdatePreferences(context.Background(), actorID, dto.UpdatePreferencesRequest{WeeklyReminderEnabled: &flag})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.WeeklyReminderEnabled != false {
		t.Fatalf("expected weekly_reminder_enabled false")
	}
	if !cancelled {
		t.Fatalf("expected weekly reminders cancelled")
	}
}
