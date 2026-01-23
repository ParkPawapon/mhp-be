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

type appointmentServiceStub struct{}

func (appointmentServiceStub) ListAppointments(ctx context.Context, userID string) ([]dto.AppointmentResponse, error) {
	return []dto.AppointmentResponse{{ID: uuid.New().String(), UserID: userID}}, nil
}
func (appointmentServiceStub) CreateAppointment(ctx context.Context, userID string, req dto.CreateAppointmentRequest) (dto.AppointmentResponse, error) {
	return dto.AppointmentResponse{ID: uuid.New().String(), UserID: userID, Title: req.Title, ApptType: req.ApptType, ApptDateTime: time.Now().UTC(), Status: constants.ApptPending}, nil
}
func (appointmentServiceStub) UpdateStatus(ctx context.Context, id string, req dto.UpdateAppointmentStatusRequest) error {
	return nil
}
func (appointmentServiceStub) DeleteAppointment(ctx context.Context, id string) error {
	return nil
}
func (appointmentServiceStub) CreateNurseVisitNote(ctx context.Context, appointmentID, nurseID string, req dto.CreateNurseVisitNoteRequest) error {
	return nil
}
func (appointmentServiceStub) ListVisitHistory(ctx context.Context, userID string) ([]dto.VisitHistoryItem, error) {
	return []dto.VisitHistoryItem{{AppointmentID: uuid.New().String(), VisitNoteID: uuid.New().String()}}, nil
}

func TestAppointmentHandlers(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RolePatient, actorID))
	handler := NewAppointmentHandler(appointmentServiceStub{}, caregiverServiceStubSimple{})

	router.GET("/appointments", handler.ListAppointments)
	router.POST("/appointments", handler.CreateAppointment)
	router.PATCH("/appointments/:id/status", handler.UpdateStatus)
	router.DELETE("/appointments/:id", handler.DeleteAppointment)
	router.POST("/appointments/:id/notes", handler.CreateNurseVisitNote)
	router.GET("/visits/history", handler.ListVisitHistory)

	resp := performRequest(router, http.MethodGet, "/appointments", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("list appointments expected 200, got %d", resp.Code)
	}

	createPayload := dto.CreateAppointmentRequest{Title: "Visit", ApptType: constants.ApptHospital, ApptDateTime: time.Now().UTC().Format(time.RFC3339)}
	resp = performRequest(router, http.MethodPost, "/appointments", createPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create appointment expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodPatch, "/appointments/"+uuid.New().String()+"/status", dto.UpdateAppointmentStatusRequest{Status: constants.ApptConfirmed})
	if resp.Code != http.StatusOK {
		t.Fatalf("update status expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodDelete, "/appointments/"+uuid.New().String(), nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("delete appointment expected 200, got %d", resp.Code)
	}

	notePayload := dto.CreateNurseVisitNoteRequest{VisitDetails: "Details"}
	resp = performRequest(router, http.MethodPost, "/appointments/"+uuid.New().String()+"/notes", notePayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create note expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/visits/history", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("visit history expected 200, got %d", resp.Code)
	}
}
