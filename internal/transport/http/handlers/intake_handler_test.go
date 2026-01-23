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

type intakeServiceStub struct{}

func (intakeServiceStub) CreateIntake(ctx context.Context, userID string, req dto.CreateIntakeRequest) (dto.IntakeHistoryResponse, error) {
	return dto.IntakeHistoryResponse{ID: uuid.New().String(), UserID: userID, Status: req.Status, TargetDate: time.Now().Format("2006-01-02")}, nil
}
func (intakeServiceStub) ListHistory(ctx context.Context, userID, from, to string) ([]dto.IntakeHistoryResponse, error) {
	return []dto.IntakeHistoryResponse{{ID: uuid.New().String(), UserID: userID, Status: constants.MedTaken}}, nil
}

type caregiverServiceStubSimple struct{}

func (caregiverServiceStubSimple) CreateAssignment(ctx context.Context, req dto.CreateCaregiverAssignmentRequest) (dto.CaregiverAssignmentResponse, error) {
	panic("not used")
}
func (caregiverServiceStubSimple) ListAssignments(ctx context.Context, patientID uuid.UUID) ([]dto.CaregiverAssignmentResponse, error) {
	panic("not used")
}
func (caregiverServiceStubSimple) IsAssigned(ctx context.Context, caregiverID, patientID uuid.UUID) (bool, error) {
	return true, nil
}

func TestIntakeHandlers(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RolePatient, actorID))
	handler := NewIntakeHandler(intakeServiceStub{}, caregiverServiceStubSimple{})

	router.POST("/intake", handler.CreateIntake)
	router.GET("/intake/history", handler.ListHistory)

	payload := dto.CreateIntakeRequest{Status: constants.MedTaken, TargetDate: time.Now().Format("2006-01-02")}
	resp := performRequest(router, http.MethodPost, "/intake", payload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/intake/history?from=2025-01-01&to=2025-01-31", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
