package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type supportServiceStub struct{}

func (supportServiceStub) CreateChatRequest(ctx context.Context, userID string, req dto.SupportChatRequestCreateRequest) (dto.SupportChatRequestResponse, error) {
	return dto.SupportChatRequestResponse{ID: uuid.New().String(), Status: "OPEN"}, nil
}

func (supportServiceStub) ListChatRequests(ctx context.Context, page, pageSize int) ([]dto.SupportChatRequestItem, int64, error) {
	return []dto.SupportChatRequestItem{{ID: uuid.New().String(), Status: "OPEN"}}, 1, nil
}

func (supportServiceStub) GetEmergencyInfo(ctx context.Context) dto.SupportEmergencyResponse {
	return dto.SupportEmergencyResponse{Hotline: "1669", DisplayName: "Emergency 1669"}
}

func TestSupportEmergency(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(constants.RequestIDKey, "req-1")
		c.Next()
	})

	handler := NewSupportHandler(supportServiceStub{})
	router.GET("/support/emergency", handler.EmergencyInfo)

	req, _ := http.NewRequest(http.MethodGet, "/support/emergency", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var body struct {
		Data dto.SupportEmergencyResponse `json:"data"`
		Meta httpx.Meta                   `json:"meta"`
	}

	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if body.Data.Hotline != "1669" || body.Data.DisplayName != "Emergency 1669" {
		t.Fatalf("unexpected payload: %#v", body.Data)
	}
	if body.Meta.RequestID != "req-1" {
		t.Fatalf("expected request_id req-1, got %s", body.Meta.RequestID)
	}
}

func TestSupportChatRequests(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RolePatient, actorID))
	handler := NewSupportHandler(supportServiceStub{})

	router.POST("/support/chat/requests", handler.CreateChatRequest)
	router.GET("/support/chat/requests", handler.ListChatRequests)

	resp := performRequest(router, http.MethodPost, "/support/chat/requests", dto.SupportChatRequestCreateRequest{Message: "Help", Category: "GENERAL"})
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/support/chat/requests?page=1&page_size=20", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
