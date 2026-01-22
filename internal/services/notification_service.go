package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/config"
	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type NotificationService interface {
	ScheduleMedicineReminders(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, mealTiming *string, timeSlot time.Time) error
	CancelMedicineAfterMealReminder(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate time.Time) error
	ScheduleAppointmentReminders(ctx context.Context, appt *db.Appointment) error
	CancelAppointmentReminders(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error
	ListUpcoming(ctx context.Context, userID string, from, to string) ([]dto.NotificationUpcomingItem, error)
	EnsureWeeklyReminders(ctx context.Context) error
	ProcessDue(ctx context.Context) error
	CancelWeeklyReminders(ctx context.Context, userID uuid.UUID) error
}

type notificationService struct {
	cfg      config.NotificationConfig
	db       *gorm.DB
	repo     repositories.NotificationRepository
	prefs    repositories.PreferenceRepository
	sender   NotificationSender
	logger   *zap.Logger
	location *time.Location
	now      func() time.Time
}

func NewNotificationService(cfg config.NotificationConfig, dbConn *gorm.DB, repo repositories.NotificationRepository, prefs repositories.PreferenceRepository, sender NotificationSender, logger *zap.Logger) NotificationService {
	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		location = time.UTC
	}
	return &notificationService{
		cfg:      cfg,
		db:       dbConn,
		repo:     repo,
		prefs:    prefs,
		sender:   sender,
		logger:   logger,
		location: location,
		now:      time.Now,
	}
}

func (s *notificationService) ScheduleMedicineReminders(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, mealTiming *string, timeSlot time.Time) error {
	days := s.cfg.ScheduleDays
	if days <= 0 {
		days = 7
	}

	nowUTC := s.now().UTC()
	nowLocal := nowUTC.In(s.location)
	events := make([]db.NotificationEvent, 0, days*2)

	for i := 0; i < days; i++ {
		date := nowLocal.AddDate(0, 0, i)
		localTime := time.Date(date.Year(), date.Month(), date.Day(), timeSlot.Hour(), timeSlot.Minute(), timeSlot.Second(), 0, s.location)
		dateKey := date.Format("2006-01-02")

		payload := map[string]any{
			"schedule_id": scheduleID.String(),
			"target_date": dateKey,
		}
		payloadBytes, _ := json.Marshal(payload)

		timing := ""
		if mealTiming != nil {
			timing = *mealTiming
		}

		switch timing {
		case constants.MealTimingBeforeMeal:
			beforeTime := localTime.Add(-5 * time.Minute).UTC()
			afterTime := localTime.Add(20 * time.Minute).UTC()
			if beforeTime.After(nowUTC) {
				events = append(events, db.NotificationEvent{
					UserID:       userID,
					TemplateCode: constants.TemplateMedBeforeMeal5Min,
					ScheduledAt:  beforeTime,
					Status:       constants.NotificationPending,
					Payload:      datatypes.JSON(payloadBytes),
				})
			}
			if afterTime.After(nowUTC) {
				events = append(events, db.NotificationEvent{
					UserID:       userID,
					TemplateCode: constants.TemplateMedBeforeMeal20Min,
					ScheduledAt:  afterTime,
					Status:       constants.NotificationPending,
					Payload:      datatypes.JSON(payloadBytes),
				})
			}
		default:
			atTime := localTime.UTC()
			if atTime.After(nowUTC) {
				events = append(events, db.NotificationEvent{
					UserID:       userID,
					TemplateCode: constants.TemplateMedAfterMealNow,
					ScheduledAt:  atTime,
					Status:       constants.NotificationPending,
					Payload:      datatypes.JSON(payloadBytes),
				})
			}
		}
	}

	return s.repo.CreateEvents(ctx, events)
}

func (s *notificationService) CancelMedicineAfterMealReminder(ctx context.Context, userID uuid.UUID, scheduleID uuid.UUID, targetDate time.Time) error {
	dateKey := targetDate.In(s.location).Format("2006-01-02")
	return s.repo.CancelPendingBySchedule(ctx, userID, scheduleID, dateKey)
}

func (s *notificationService) ScheduleAppointmentReminders(ctx context.Context, appt *db.Appointment) error {
	if appt == nil {
		return nil
	}
	if appt.ApptType != constants.ApptHospital {
		return nil
	}

	nowUTC := s.now().UTC()
	payload := map[string]any{"appointment_id": appt.ID.String()}
	payloadBytes, _ := json.Marshal(payload)

	reminders := []struct {
		Code string
		Days int
	}{
		{Code: constants.TemplateAppt5Days, Days: -5},
		{Code: constants.TemplateAppt1Day, Days: -1},
	}

	events := make([]db.NotificationEvent, 0, len(reminders))
	for _, reminder := range reminders {
		scheduled := appt.ApptDateTime.AddDate(0, 0, reminder.Days).UTC()
		if scheduled.After(nowUTC) {
			events = append(events, db.NotificationEvent{
				UserID:       appt.UserID,
				TemplateCode: reminder.Code,
				ScheduledAt:  scheduled,
				Status:       constants.NotificationPending,
				Payload:      datatypes.JSON(payloadBytes),
			})
		}
	}

	return s.repo.CreateEvents(ctx, events)
}

func (s *notificationService) CancelAppointmentReminders(ctx context.Context, userID uuid.UUID, appointmentID uuid.UUID) error {
	return s.repo.CancelPendingByAppointment(ctx, userID, appointmentID)
}

func (s *notificationService) ListUpcoming(ctx context.Context, userID string, from, to string) ([]dto.NotificationUpcomingItem, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	var fromTime time.Time
	var toTime time.Time
	if from != "" {
		if parsed, err := time.Parse(time.RFC3339, from); err == nil {
			fromTime = parsed.UTC()
		} else {
			return nil, domain.NewError(constants.ValidationFailed, "invalid from")
		}
	}
	if to != "" {
		if parsed, err := time.Parse(time.RFC3339, to); err == nil {
			toTime = parsed.UTC()
		} else {
			return nil, domain.NewError(constants.ValidationFailed, "invalid to")
		}
	}

	events, err := s.repo.ListUpcoming(ctx, uid, fromTime, toTime)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.NotificationUpcomingItem, 0, len(events))
	for _, event := range events {
		resp = append(resp, dto.NotificationUpcomingItem{
			ID:           event.ID.String(),
			TemplateCode: event.TemplateCode,
			ScheduledAt:  event.ScheduledAt,
			Status:       string(event.Status),
		})
	}
	return resp, nil
}

func (s *notificationService) EnsureWeeklyReminders(ctx context.Context) error {
	if s.prefs == nil {
		return nil
	}
	users, err := s.prefs.ListWeeklyReminderUsers(ctx)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	localNow := s.now().In(s.location)
	scheduledLocal := nextWeeklyTime(localNow, s.cfg.WeeklyReminderHour, s.cfg.WeeklyReminderMinute)
	scheduledUTC := scheduledLocal.UTC()

	payloadBytes, _ := json.Marshal(map[string]any{"type": "weekly_health_log"})
	events := make([]db.NotificationEvent, 0, len(users))
	for _, userID := range users {
		events = append(events, db.NotificationEvent{
			UserID:       userID,
			TemplateCode: constants.TemplateWeeklyHealthLog,
			ScheduledAt:  scheduledUTC,
			Status:       constants.NotificationPending,
			Payload:      datatypes.JSON(payloadBytes),
		})
	}

	return s.repo.CreateEvents(ctx, events)
}

func (s *notificationService) ProcessDue(ctx context.Context) error {
	if s.db == nil {
		return nil
	}
	batchSize := s.cfg.JobBatchSize
	if batchSize <= 0 {
		batchSize = 100
	}

	nowUTC := s.now().UTC()
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := s.repo.WithTx(tx)
		events, err := repo.ListDueForUpdate(ctx, nowUTC, batchSize)
		if err != nil {
			return err
		}
		for _, event := range events {
			tpl, err := repo.FindTemplateByCode(ctx, event.TemplateCode)
			if err != nil {
				_ = repo.UpdateEventStatus(ctx, event.ID, constants.NotificationFailed, nil)
				continue
			}

			status := constants.NotificationSent
			sentAt := s.now().UTC()
			if s.sender != nil {
				if err := s.sender.Send(ctx, event, *tpl); err != nil {
					status = constants.NotificationFailed
					if s.logger != nil {
						s.logger.Warn("notification send failed", zap.String("request_id", "job"), zap.String("user_id", event.UserID.String()), zap.String("template_code", event.TemplateCode), zap.Error(err))
					}
				}
			}

			if err := repo.UpdateEventStatus(ctx, event.ID, status, &sentAt); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *notificationService) CancelWeeklyReminders(ctx context.Context, userID uuid.UUID) error {
	return s.repo.CancelPendingByTemplate(ctx, userID, constants.TemplateWeeklyHealthLog)
}

func nextWeeklyTime(now time.Time, hour, minute int) time.Time {
	scheduled := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	weekday := time.Monday
	if now.Weekday() != weekday || scheduled.Before(now) {
		daysUntil := (int(weekday) - int(now.Weekday()) + 7) % 7
		if daysUntil == 0 {
			daysUntil = 7
		}
		scheduled = scheduled.AddDate(0, 0, daysUntil)
	}
	return scheduled
}
