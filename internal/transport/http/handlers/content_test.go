package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type contentServiceStub struct{}

func (contentServiceStub) ListHealthContent(ctx context.Context, publishedOnly bool) ([]dto.HealthContentResponse, error) {
	return []dto.HealthContentResponse{{ID: uuid.New().String(), Title: "Title"}}, nil
}
func (contentServiceStub) CreateHealthContent(ctx context.Context, req dto.CreateHealthContentRequest) (dto.HealthContentResponse, error) {
	return dto.HealthContentResponse{ID: uuid.New().String(), Title: req.Title}, nil
}
func (contentServiceStub) UpdateHealthContent(ctx context.Context, id string, req dto.UpdateHealthContentRequest) error {
	return nil
}
func (contentServiceStub) PublishHealthContent(ctx context.Context, id string, req dto.PublishHealthContentRequest) error {
	return nil
}
func (contentServiceStub) ListHealthCategories(ctx context.Context) []string {
	return constants.HealthContentCategories
}

func TestContentHandlers(t *testing.T) {
	actorID := uuid.New()
	router := newTestRouter(withActor(constants.RoleNurse, actorID))
	handler := NewContentHandler(contentServiceStub{})

	router.GET("/content/health/categories", handler.ListHealthCategories)
	router.GET("/content/health", handler.ListHealthContent)
	router.POST("/content/health", handler.CreateHealthContent)
	router.PATCH("/content/health/:id", handler.UpdateHealthContent)
	router.POST("/content/health/:id/publish", handler.PublishHealthContent)

	resp := performRequest(router, http.MethodGet, "/content/health/categories", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("categories expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodGet, "/content/health?published=true", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("list content expected 200, got %d", resp.Code)
	}

	createPayload := dto.CreateHealthContentRequest{Title: "Title"}
	resp = performRequest(router, http.MethodPost, "/content/health", createPayload)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create content expected 201, got %d", resp.Code)
	}

	updatePayload := dto.UpdateHealthContentRequest{Title: strPtr("Updated")}
	resp = performRequest(router, http.MethodPatch, "/content/health/"+uuid.New().String(), updatePayload)
	if resp.Code != http.StatusOK {
		t.Fatalf("update content expected 200, got %d", resp.Code)
	}

	resp = performRequest(router, http.MethodPost, "/content/health/"+uuid.New().String()+"/publish", dto.PublishHealthContentRequest{IsPublished: true})
	if resp.Code != http.StatusOK {
		t.Fatalf("publish content expected 200, got %d", resp.Code)
	}
}
