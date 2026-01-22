package jobs

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/ParkPawapon/mhp-be/internal/services"
)

type NotificationWorker struct {
	service  services.NotificationService
	interval time.Duration
	logger   *zap.Logger
}

func NewNotificationWorker(service services.NotificationService, interval time.Duration, logger *zap.Logger) *NotificationWorker {
	return &NotificationWorker{service: service, interval: interval, logger: logger}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	if w.service == nil {
		return
	}

	if w.interval <= 0 {
		w.interval = time.Minute
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *NotificationWorker) runOnce(ctx context.Context) {
	if err := w.service.EnsureWeeklyReminders(ctx); err != nil && w.logger != nil {
		w.logger.Warn("weekly reminders failed", zap.String("request_id", "job"), zap.Error(err))
	}
	if err := w.service.ProcessDue(ctx); err != nil && w.logger != nil {
		w.logger.Warn("process notifications failed", zap.String("request_id", "job"), zap.Error(err))
	}
}
