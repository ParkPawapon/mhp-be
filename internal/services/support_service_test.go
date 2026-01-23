package services

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type supportRepoStub struct {
	created *db.SupportChatRequest
}

func (s *supportRepoStub) CreateChatRequest(ctx context.Context, req *db.SupportChatRequest) error {
	s.created = req
	req.ID = uuid.New()
	return nil
}
func (s *supportRepoStub) ListChatRequests(ctx context.Context, page, pageSize int) ([]db.SupportChatRequest, int64, error) {
	return []db.SupportChatRequest{}, 0, nil
}

func TestSupportServiceValidation(t *testing.T) {
	repo := &supportRepoStub{}
	svc := NewSupportService(repo)

	_, err := svc.CreateChatRequest(context.Background(), "bad", dto.SupportChatRequestCreateRequest{Message: "hi", Category: "GENERAL"})
	if err == nil {
		t.Fatalf("expected invalid user_id error")
	}

	_, err = svc.CreateChatRequest(context.Background(), uuid.New().String(), dto.SupportChatRequestCreateRequest{Message: "hi", Category: "UNKNOWN"})
	if err == nil {
		t.Fatalf("expected invalid category error")
	}

	resp, err := svc.CreateChatRequest(context.Background(), uuid.New().String(), dto.SupportChatRequestCreateRequest{Message: "hi", Category: constants.SupportCategories[0]})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID == "" || resp.Status != "OPEN" {
		t.Fatalf("unexpected response")
	}
	if repo.created == nil || repo.created.Category != constants.SupportCategories[0] {
		t.Fatalf("expected category stored")
	}
}
