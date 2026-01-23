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

type userServiceStub struct{}

func (userServiceStub) GetMe(ctx context.Context, actorID uuid.UUID, role constants.Role) (dto.MeResponse, error) {
	return dto.MeResponse{ID: actorID.String(), Role: role, Profile: dto.ProfileResponse{FirstName: "A", LastName: "B"}}, nil
}
func (userServiceStub) UpdateProfile(ctx context.Context, actorID uuid.UUID, req dto.UpdateProfileRequest) error {
	return nil
}
func (userServiceStub) SaveDeviceToken(ctx context.Context, actorID uuid.UUID, req dto.DeviceTokenRequest) error {
	return nil
}
func (userServiceStub) UpdatePreferences(ctx context.Context, actorID uuid.UUID, req dto.UpdatePreferencesRequest) (dto.PreferencesResponse, error) {
	return dto.PreferencesResponse{WeeklyReminderEnabled: true}, nil
}

func TestUserHandlers(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RolePatient, actorID))
	handler := NewUserHandler(userServiceStub{})

	router.GET("/me", handler.Me)
	router.PATCH("/me/profile", handler.UpdateProfile)
	router.PATCH("/me/preferences", handler.UpdatePreferences)
	router.POST("/me/device-tokens", handler.SaveDeviceToken)

	resp := performRequest(router, http.MethodGet, "/me", nil)
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

	first := "Jane"
	last := "Doe"
	dob := time.Now().AddDate(-20, 0, 0)
	resp = performRequest(router, http.MethodPatch, "/me/profile", dto.UpdateProfileRequest{FirstName: &first, LastName: &last, DateOfBirth: &dob})
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodPatch, "/me/preferences", dto.UpdatePreferencesRequest{WeeklyReminderEnabled: boolPtr(true)})
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodPost, "/me/device-tokens", dto.DeviceTokenRequest{Platform: "ios", DeviceToken: "token"})
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func boolPtr(v bool) *bool {
	return &v
}
