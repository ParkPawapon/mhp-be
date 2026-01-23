package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type caregiverServiceStub struct{}

func (caregiverServiceStub) CreateAssignment(ctx context.Context, req dto.CreateCaregiverAssignmentRequest) (dto.CaregiverAssignmentResponse, error) {
	return dto.CaregiverAssignmentResponse{ID: uuid.New().String(), PatientID: req.PatientID, CaregiverID: req.CaregiverID, Relationship: req.Relationship}, nil
}
func (caregiverServiceStub) ListAssignments(ctx context.Context, patientID uuid.UUID) ([]dto.CaregiverAssignmentResponse, error) {
	return []dto.CaregiverAssignmentResponse{{ID: uuid.New().String(), PatientID: patientID.String(), CaregiverID: uuid.New().String(), Relationship: "family"}}, nil
}
func (caregiverServiceStub) IsAssigned(ctx context.Context, caregiverID, patientID uuid.UUID) (bool, error) {
	return true, nil
}

func TestCaregiverHandlers(t *testing.T) {
	router := newTestRouter()
	handler := NewCaregiverHandler(caregiverServiceStub{})

	router.POST("/caregivers/assignments", handler.CreateAssignment)
	router.GET("/caregivers/assignments", handler.ListAssignments)

	payload := dto.CreateCaregiverAssignmentRequest{PatientID: uuid.New().String(), CaregiverID: uuid.New().String(), Relationship: "family"}
	resp := performRequest(router, http.MethodPost, "/caregivers/assignments", payload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/caregivers/assignments?patient_id="+payload.PatientID, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var meta envelopeMeta
	if err := json.Unmarshal(resp.Body.Bytes(), &meta); err != nil {
		t.Fatalf("invalid json")
	}
	if meta.Meta.RequestID != testRequestID {
		t.Fatalf("expected request_id")
	}
}
