package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type auditServiceStub struct{}

func (auditServiceStub) ListAuditLogs(ctx context.Context, page, pageSize int, from, to, actorID, actionType string) ([]dto.AuditLogResponse, int64, error) {
	return []dto.AuditLogResponse{{ID: uuid.New().String()}}, 1, nil
}

func TestAuditHandlers(t *testing.T) {
	router := newTestRouter()
	handler := NewAuditHandler(auditServiceStub{})

	router.GET("/admin/audit-logs", handler.ListAuditLogs)

	resp := performRequest(router, http.MethodGet, "/admin/audit-logs?page=1&page_size=20", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
