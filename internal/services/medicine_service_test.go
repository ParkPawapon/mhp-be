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

type medicineRepoStub struct {
	master          *db.MedicineMaster
	categoryItem    *db.MedicineCategoryItem
	patientMedicine *db.PatientMedicine
	createdMedicine *db.PatientMedicine
	createdSchedule *db.MedicineSchedule
}

func (s *medicineRepoStub) ListMaster(ctx context.Context, page, pageSize int) ([]db.MedicineMaster, int64, error) {
	panic("not used")
}
func (s *medicineRepoStub) GetMasterByID(ctx context.Context, id uuid.UUID) (*db.MedicineMaster, error) {
	if s.master == nil {
		return nil, domain.NewError(constants.MedNotFound, "not found")
	}
	return s.master, nil
}
func (s *medicineRepoStub) CreatePatientMedicine(ctx context.Context, med *db.PatientMedicine) error {
	s.createdMedicine = med
	med.ID = uuid.New()
	med.CreatedAt = time.Now().UTC()
	return nil
}
func (s *medicineRepoStub) ListPatientMedicines(ctx context.Context, userID uuid.UUID) ([]db.PatientMedicine, error) {
	panic("not used")
}
func (s *medicineRepoStub) GetPatientMedicineByID(ctx context.Context, id uuid.UUID) (*db.PatientMedicine, error) {
	if s.patientMedicine == nil {
		return nil, domain.NewError(constants.MedNotFound, "not found")
	}
	return s.patientMedicine, nil
}
func (s *medicineRepoStub) UpdatePatientMedicine(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	panic("not used")
}
func (s *medicineRepoStub) DeletePatientMedicine(ctx context.Context, id uuid.UUID) error {
	panic("not used")
}
func (s *medicineRepoStub) CreateSchedule(ctx context.Context, schedule *db.MedicineSchedule) error {
	s.createdSchedule = schedule
	schedule.ID = uuid.New()
	schedule.CreatedAt = time.Now().UTC()
	return nil
}
func (s *medicineRepoStub) GetScheduleByID(ctx context.Context, id uuid.UUID) (*db.MedicineSchedule, error) {
	panic("not used")
}
func (s *medicineRepoStub) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	panic("not used")
}
func (s *medicineRepoStub) ListCategories(ctx context.Context) ([]db.MedicineCategory, error) {
	panic("not used")
}
func (s *medicineRepoStub) ListCategoryItems(ctx context.Context, categoryID uuid.UUID) ([]db.MedicineCategoryItem, error) {
	panic("not used")
}
func (s *medicineRepoStub) GetCategoryItemByID(ctx context.Context, id uuid.UUID) (*db.MedicineCategoryItem, error) {
	if s.categoryItem == nil {
		return nil, domain.NewError(constants.MedNotFound, "not found")
	}
	return s.categoryItem, nil
}

type notificationScheduleStub struct {
	called bool
}

func (s *notificationScheduleStub) ScheduleMedicineReminders(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, mealTiming *string, timeSlot time.Time) error {
	s.called = true
	return nil
}
func (s *notificationScheduleStub) CancelMedicineAfterMealReminder(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate time.Time) error {
	panic("not used")
}
func (s *notificationScheduleStub) ScheduleAppointmentReminders(ctx context.Context, appt *db.Appointment) error {
	panic("not used")
}
func (s *notificationScheduleStub) CancelAppointmentReminders(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	panic("not used")
}
func (s *notificationScheduleStub) ListUpcoming(ctx context.Context, userID string, from, to string) ([]dto.NotificationUpcomingItem, error) {
	panic("not used")
}
func (s *notificationScheduleStub) EnsureWeeklyReminders(ctx context.Context) error {
	panic("not used")
}
func (s *notificationScheduleStub) ProcessDue(ctx context.Context) error {
	panic("not used")
}
func (s *notificationScheduleStub) CancelWeeklyReminders(ctx context.Context, userID uuid.UUID) error {
	panic("not used")
}

func TestCreatePatientMedicineRequiresSource(t *testing.T) {
	repo := &medicineRepoStub{}
	svc := NewMedicineService(repo, nil)

	_, err := svc.CreatePatientMedicine(context.Background(), uuid.New().String(), dto.CreatePatientMedicineRequest{
		DosageAmount: "1",
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestCreatePatientMedicineFromCategoryItem(t *testing.T) {
	itemID := uuid.New()
	itemName := "Amlodipine 5 mg"
	dosage := "1"
	repo := &medicineRepoStub{
		categoryItem: &db.MedicineCategoryItem{
			ID:                itemID,
			DisplayName:       itemName,
			DefaultDosageText: &dosage,
		},
	}
	svc := NewMedicineService(repo, nil)

	resp, err := svc.CreatePatientMedicine(context.Background(), uuid.New().String(), dto.CreatePatientMedicineRequest{
		CategoryItemID: &[]string{itemID.String()}[0],
		DosageAmount:   "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CustomName == nil || *resp.CustomName != itemName {
		t.Fatalf("expected custom name from category item")
	}
	if resp.DosageAmount != dosage {
		t.Fatalf("expected dosage default, got %q", resp.DosageAmount)
	}
}

func TestCreateScheduleValidatesMealTiming(t *testing.T) {
	medID := uuid.New()
	repo := &medicineRepoStub{patientMedicine: &db.PatientMedicine{ID: medID, UserID: uuid.New()}}
	notify := &notificationScheduleStub{}
	svc := NewMedicineService(repo, notify)

	_, err := svc.CreateSchedule(context.Background(), medID.String(), dto.CreateMedicineScheduleRequest{
		TimeSlot:   "08:00",
		MealTiming: strPtr("INVALID"),
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	_, err = svc.CreateSchedule(context.Background(), medID.String(), dto.CreateMedicineScheduleRequest{
		TimeSlot:   "08:00",
		MealTiming: strPtr(constants.MealTimingBeforeMeal),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !notify.called {
		t.Fatalf("expected notification schedule called")
	}
}
