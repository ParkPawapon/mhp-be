package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type adminServiceStub struct{}

func (adminServiceStub) StaffLogin(ctx context.Context, req dto.StaffLoginRequest) (dto.TokenResponse, error) {
	return dto.TokenResponse{AccessToken: "a", RefreshToken: "r"}, nil
}
func (adminServiceStub) ListPatients(ctx context.Context, page, pageSize int) ([]dto.PatientSummaryResponse, int64, error) {
	return []dto.PatientSummaryResponse{{ID: uuid.New().String()}}, 1, nil
}
func (adminServiceStub) GetPatient(ctx context.Context, id string) (dto.PatientDetailResponse, error) {
	return dto.PatientDetailResponse{ID: id}, nil
}
func (adminServiceStub) ListAdherence(ctx context.Context, patientID, from, to string) ([]dto.IntakeHistoryResponse, error) {
	return []dto.IntakeHistoryResponse{{ID: uuid.New().String(), UserID: patientID}}, nil
}

func TestAdminHandlers(t *testing.T) {
	router := newTestRouter()
	handler := NewAdminHandler(adminServiceStub{})

	router.POST("/staff/login", handler.StaffLogin)
	router.GET("/admin/patients", handler.ListPatients)
	router.GET("/admin/patients/:id", handler.GetPatient)
	router.GET("/admin/adherence", handler.ListAdherence)

	resp := performRequest(router, http.MethodPost, "/staff/login", dto.StaffLoginRequest{Username: "admin", Password: "pass"})
	if resp.Code != http.StatusOK {
		t.Fatalf("staff login expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/admin/patients?page=1&page_size=10", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("list patients expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/admin/patients/"+uuid.New().String(), nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("get patient expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/admin/adherence?patient_id="+uuid.New().String(), nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("adherence expected 200, got %d", resp.Code)
	}
}
