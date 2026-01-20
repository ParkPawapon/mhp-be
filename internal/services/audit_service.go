package services

import (
	"context"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
)

type AuditService interface {
	ListAuditLogs(ctx context.Context, page, pageSize int, from, to, actorID, actionType string) ([]dto.AuditLogResponse, int64, error)
}

type auditService struct{}

func NewAuditService() AuditService {
	return &auditService{}
}

func (s *auditService) ListAuditLogs(ctx context.Context, page, pageSize int, from, to, actorID, actionType string) ([]dto.AuditLogResponse, int64, error) {
	return nil, 0, domain.NewError(constants.InternalNotImplemented, "audit logs not implemented")
}
