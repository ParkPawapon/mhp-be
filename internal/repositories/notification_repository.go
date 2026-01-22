package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type NotificationRepository interface {
	WithTx(tx *gorm.DB) NotificationRepository
	CreateEvents(ctx context.Context, events []db.NotificationEvent) error
	ListUpcoming(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]db.NotificationEvent, error)
	ListDueForUpdate(ctx context.Context, now time.Time, limit int) ([]db.NotificationEvent, error)
	UpdateEventStatus(ctx context.Context, id uuid.UUID, status constants.NotificationStatus, sentAt *time.Time) error
	FindTemplateByCode(ctx context.Context, code string) (*db.NotificationTemplate, error)
	CancelPendingBySchedule(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate string) error
	CancelPendingByAppointment(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error
	CancelPendingByTemplate(ctx context.Context, userID uuid.UUID, templateCode string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(dbConn *gorm.DB) NotificationRepository {
	return &notificationRepository{db: dbConn}
}

func (r *notificationRepository) WithTx(tx *gorm.DB) NotificationRepository {
	return &notificationRepository{db: tx}
}

func (r *notificationRepository) CreateEvents(ctx context.Context, events []db.NotificationEvent) error {
	if len(events) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "template_code"}, {Name: "scheduled_at"}},
			DoNothing: true,
		}).
		CreateInBatches(events, 100).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create notification events failed", err)
	}
	return nil
}

func (r *notificationRepository) ListUpcoming(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]db.NotificationEvent, error) {
	var items []db.NotificationEvent
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if !from.IsZero() {
		query = query.Where("scheduled_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("scheduled_at <= ?", to)
	}
	if err := query.Order("scheduled_at asc").Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list notification events failed", err)
	}
	return items, nil
}

func (r *notificationRepository) ListDueForUpdate(ctx context.Context, now time.Time, limit int) ([]db.NotificationEvent, error) {
	var items []db.NotificationEvent
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ? AND scheduled_at <= ?", constants.NotificationPending, now).
		Order("scheduled_at asc").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list due notifications failed", err)
	}
	return items, nil
}

func (r *notificationRepository) UpdateEventStatus(ctx context.Context, id uuid.UUID, status constants.NotificationStatus, sentAt *time.Time) error {
	updates := map[string]any{"status": status}
	if sentAt != nil {
		updates["sent_at"] = *sentAt
	}
	if err := r.db.WithContext(ctx).Model(&db.NotificationEvent{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return domain.WrapError(constants.InternalError, "update notification status failed", err)
	}
	return nil
}

func (r *notificationRepository) FindTemplateByCode(ctx context.Context, code string) (*db.NotificationTemplate, error) {
	var tpl db.NotificationTemplate
	if err := r.db.WithContext(ctx).Where("code = ? AND is_active = ?", code, true).First(&tpl).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.InternalError, "notification template not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find notification template failed", err)
	}
	return &tpl, nil
}

func (r *notificationRepository) CancelPendingBySchedule(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate string) error {
	if err := r.db.WithContext(ctx).
		Model(&db.NotificationEvent{}).
		Where("user_id = ? AND status = ? AND template_code = ? AND payload->>'schedule_id' = ? AND payload->>'target_date' = ?", userID, constants.NotificationPending, constants.TemplateMedBeforeMeal20Min, scheduleID.String(), targetDate).
		Update("status", constants.NotificationCancelled).Error; err != nil {
		return domain.WrapError(constants.InternalError, "cancel medicine reminder failed", err)
	}
	return nil
}

func (r *notificationRepository) CancelPendingByAppointment(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&db.NotificationEvent{}).
		Where("user_id = ? AND status = ? AND (template_code = ? OR template_code = ?) AND payload->>'appointment_id' = ?", userID, constants.NotificationPending, constants.TemplateAppt5Days, constants.TemplateAppt1Day, appointmentID.String()).
		Update("status", constants.NotificationCancelled).Error; err != nil {
		return domain.WrapError(constants.InternalError, "cancel appointment reminders failed", err)
	}
	return nil
}

func (r *notificationRepository) CancelPendingByTemplate(ctx context.Context, userID uuid.UUID, templateCode string) error {
	if err := r.db.WithContext(ctx).
		Model(&db.NotificationEvent{}).
		Where("user_id = ? AND status = ? AND template_code = ?", userID, constants.NotificationPending, templateCode).
		Update("status", constants.NotificationCancelled).Error; err != nil {
		return domain.WrapError(constants.InternalError, "cancel notification events failed", err)
	}
	return nil
}
