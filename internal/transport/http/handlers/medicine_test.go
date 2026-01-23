package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type medicineServiceStub struct{}

func (medicineServiceStub) ListMaster(ctx context.Context, page, pageSize int) ([]dto.MedicineMasterResponse, int64, error) {
	return []dto.MedicineMasterResponse{{ID: uuid.New().String(), TradeName: "A"}}, 1, nil
}
func (medicineServiceStub) CreatePatientMedicine(ctx context.Context, userID string, req dto.CreatePatientMedicineRequest) (dto.PatientMedicineResponse, error) {
	return dto.PatientMedicineResponse{ID: uuid.New().String(), UserID: userID, DosageAmount: req.DosageAmount}, nil
}
func (medicineServiceStub) ListPatientMedicines(ctx context.Context, userID string) ([]dto.PatientMedicineResponse, error) {
	return []dto.PatientMedicineResponse{{ID: uuid.New().String(), UserID: userID, DosageAmount: "1"}}, nil
}
func (medicineServiceStub) UpdatePatientMedicine(ctx context.Context, id string, req dto.UpdatePatientMedicineRequest) error {
	return nil
}
func (medicineServiceStub) DeletePatientMedicine(ctx context.Context, id string) error {
	return nil
}
func (medicineServiceStub) CreateSchedule(ctx context.Context, patientMedicineID string, req dto.CreateMedicineScheduleRequest) (dto.MedicineScheduleResponse, error) {
	return dto.MedicineScheduleResponse{ID: uuid.New().String(), PatientMedicineID: patientMedicineID, TimeSlot: req.TimeSlot, CreatedAt: time.Now().UTC()}, nil
}
func (medicineServiceStub) DeleteSchedule(ctx context.Context, id string) error {
	return nil
}
func (medicineServiceStub) ListCategories(ctx context.Context) ([]dto.MedicineCategoryResponse, error) {
	return []dto.MedicineCategoryResponse{{ID: uuid.New().String(), Name: "Hypertension", Code: "BP"}}, nil
}
func (medicineServiceStub) ListCategoryItems(ctx context.Context, categoryID string) ([]dto.MedicineCategoryItemResponse, error) {
	return []dto.MedicineCategoryItemResponse{{ID: uuid.New().String(), CategoryID: categoryID, DisplayName: "Amlodipine 5 mg"}}, nil
}
func (medicineServiceStub) GetDosageOptions(ctx context.Context) []string {
	return []string{"1/2", "1"}
}
func (medicineServiceStub) GetMealTimingOptions(ctx context.Context) []string {
	return []string{constants.MealTimingBeforeMeal}
}

func TestMedicineHandlers(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RolePatient, actorID))
	handler := NewMedicineHandler(medicineServiceStub{})

	router.GET("/medicines/master", handler.ListMaster)
	router.GET("/medicines/categories", handler.ListCategories)
	router.GET("/medicines/categories/:id/items", handler.ListCategoryItems)
	router.GET("/medicines/dosage-options", handler.GetDosageOptions)
	router.GET("/medicines/meal-timing-options", handler.GetMealTimingOptions)
	router.POST("/medicines/patient", handler.CreatePatientMedicine)
	router.GET("/medicines/patient", handler.ListPatientMedicines)
	router.PATCH("/medicines/patient/:id", handler.UpdatePatientMedicine)
	router.DELETE("/medicines/patient/:id", handler.DeletePatientMedicine)
	router.POST("/medicines/patient/:id/schedules", handler.CreateSchedule)
	router.DELETE("/medicines/schedules/:id", handler.DeleteSchedule)

	resp := performRequest(router, http.MethodGet, "/medicines/master?page=1&page_size=20", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("master: expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/medicines/categories", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("categories: expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/medicines/categories/"+uuid.New().String()+"/items", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("category items: expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/medicines/dosage-options", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("dosage options: expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/medicines/meal-timing-options", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("meal timing: expected 200, got %d", resp.Code)
	}

	createPayload := dto.CreatePatientMedicineRequest{DosageAmount: "1", CustomName: strPtr("Custom")}
	resp = performRequest(router, http.MethodPost, "/medicines/patient", createPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create patient medicine: expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/medicines/patient", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("list patient medicines: expected 200, got %d", resp.Code)
	}

	updatePayload := dto.UpdatePatientMedicineRequest{DosageAmount: strPtr("2")}
	resp = performRequest(router, http.MethodPatch, "/medicines/patient/"+uuid.New().String(), updatePayload)
	if resp.Code != http.StatusOK {
		t.Fatalf("update patient medicine: expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodDelete, "/medicines/patient/"+uuid.New().String(), nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("delete patient medicine: expected 200, got %d", resp.Code)
	}

	schedulePayload := dto.CreateMedicineScheduleRequest{TimeSlot: "08:00", MealTiming: strPtr(constants.MealTimingBeforeMeal)}
	resp = performRequest(router, http.MethodPost, "/medicines/patient/"+uuid.New().String()+"/schedules", schedulePayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create schedule: expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodDelete, "/medicines/schedules/"+uuid.New().String(), nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("delete schedule: expected 200, got %d", resp.Code)
	}

	var meta envelopeMeta
	if err := json.Unmarshal(resp.Body.Bytes(), &meta); err != nil {
		t.Fatalf("invalid json")
	}
	if meta.Meta.RequestID != testRequestID {
		t.Fatalf("expected request_id")
	}
}
