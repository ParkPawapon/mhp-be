package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/transport/httpx"
)

type supportServiceStub struct{}

func (supportServiceStub) CreateChatRequest(ctx context.Context, userID string, req dto.SupportChatRequestCreateRequest) (dto.SupportChatRequestResponse, error) {
	return dto.SupportChatRequestResponse{}, nil
}

func (supportServiceStub) ListChatRequests(ctx context.Context, page, pageSize int) ([]dto.SupportChatRequestItem, int64, error) {
	return nil, 0, nil
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
