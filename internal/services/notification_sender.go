package services

import (
	"context"

	"go.uber.org/zap"

	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type NotificationSender interface {
	Send(ctx context.Context, event db.NotificationEvent, template db.NotificationTemplate) error
}

type ConsoleNotificationSender struct {
	Logger *zap.Logger
}

func (s ConsoleNotificationSender) Send(ctx context.Context, event db.NotificationEvent, template db.NotificationTemplate) error {
	if s.Logger == nil {
		return nil
	}

	s.Logger.Info("notification send", zap.String("request_id", "job"), zap.String("user_id", event.UserID.String()), zap.String("template_code", event.TemplateCode), zap.Time("scheduled_at", event.ScheduledAt))
	return nil
}
