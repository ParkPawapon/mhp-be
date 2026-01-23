package handlers

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type healthServiceStub struct{}

func (healthServiceStub) CreateHealthRecord(ctx context.Context, userID string, req dto.CreateHealthRecordRequest) (dto.HealthRecordResponse, error) {
	return dto.HealthRecordResponse{ID: uuid.New().String(), RecordDate: req.RecordDate}, nil
}
func (healthServiceStub) ListHealthRecords(ctx context.Context, userID, from, to string) ([]dto.HealthRecordResponse, error) {
	return []dto.HealthRecordResponse{{ID: uuid.New().String()}}, nil
}
func (healthServiceStub) CreateDailyAssessment(ctx context.Context, userID string, req dto.CreateDailyAssessmentRequest) (dto.DailyAssessmentResponse, error) {
	return dto.DailyAssessmentResponse{ID: uuid.New().String(), LogDate: req.LogDate}, nil
}
func (healthServiceStub) ListDailyAssessments(ctx context.Context, userID, from, to string) ([]dto.DailyAssessmentResponse, error) {
	return []dto.DailyAssessmentResponse{{ID: uuid.New().String()}}, nil
}

func TestHealthRecordHandlers(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RolePatient, actorID))
	handler := NewHealthRecordsHandler(healthServiceStub{}, caregiverServiceStubSimple{})

	router.POST("/health/records", handler.CreateHealthRecord)
	router.GET("/health/records", handler.ListHealthRecords)
	router.POST("/assessments/daily", handler.CreateDailyAssessment)
	router.GET("/assessments/daily", handler.ListDailyAssessments)

	recordPayload := dto.CreateHealthRecordRequest{RecordDate: time.Now().Format("2006-01-02"), TimePeriod: "morning"}
	resp := performRequest(router, http.MethodPost, "/health/records", recordPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/health/records?from=2025-01-01&to=2025-01-31", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	assessmentPayload := dto.CreateDailyAssessmentRequest{LogDate: time.Now().Format("2006-01-02")}
	resp = performRequest(router, http.MethodPost, "/assessments/daily", assessmentPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/assessments/daily?from=2025-01-01&to=2025-01-31", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
