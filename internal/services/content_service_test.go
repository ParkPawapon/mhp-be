package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type contentRepoStub struct {
	created *db.HealthContent
	updated map[string]any
}

func (s *contentRepoStub) ListHealthContent(ctx context.Context, publishedOnly bool) ([]db.HealthContent, error) {
	return []db.HealthContent{}, nil
}
func (s *contentRepoStub) CreateHealthContent(ctx context.Context, content *db.HealthContent) error {
	content.ID = uuid.New()
	content.CreatedAt = time.Now().UTC()
	s.created = content
	return nil
}
func (s *contentRepoStub) UpdateHealthContent(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	s.updated = updates
	return nil
}
func (s *contentRepoStub) SetPublished(ctx context.Context, id uuid.UUID, published bool) error {
	return nil
}
func (s *contentRepoStub) FindByID(ctx context.Context, id uuid.UUID) (*db.HealthContent, error) {
	return s.created, nil
}

func TestContentCategoryValidation(t *testing.T) {
	repo := &contentRepoStub{}
	svc := NewContentService(repo)

	_, err := svc.CreateHealthContent(context.Background(), dto.CreateHealthContentRequest{
		Title:    "Title",
		Category: strPtr("INVALID"),
	})
	if err == nil {
		t.Fatalf("expected invalid category error")
	}

	valid := constants.HealthContentCategories[0]
	_, err = svc.CreateHealthContent(context.Background(), dto.CreateHealthContentRequest{
		Title:    "Title",
		Category: &valid,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContentUpdateValidation(t *testing.T) {
	repo := &contentRepoStub{}
	svc := NewContentService(repo)

	if err := svc.UpdateHealthContent(context.Background(), "bad-id", dto.UpdateHealthContentRequest{}); err == nil {
		t.Fatalf("expected invalid id error")
	}

	invalid := "INVALID"
	if err := svc.UpdateHealthContent(context.Background(), uuid.New().String(), dto.UpdateHealthContentRequest{Category: &invalid}); err == nil {
		t.Fatalf("expected invalid category error")
	}

	if err := svc.UpdateHealthContent(context.Background(), uuid.New().String(), dto.UpdateHealthContentRequest{}); err == nil {
		t.Fatalf("expected no fields error")
	}

	value := "  Updated "
	if err := svc.UpdateHealthContent(context.Background(), uuid.New().String(), dto.UpdateHealthContentRequest{Title: &value}); err != nil {
		if appErr, ok := domain.AsAppError(err); !ok || appErr.Code != constants.ValidationFailed {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if repo.updated == nil || repo.updated["title"] != "Updated" {
		t.Fatalf("expected trimmed title update")
	}
}
