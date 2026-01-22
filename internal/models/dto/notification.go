package dto

import "time"

type NotificationUpcomingItem struct {
	ID           string    `json:"id"`
	TemplateCode string    `json:"template_code"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	Status       string    `json:"status"`
}

type UpdatePreferencesRequest struct {
	WeeklyReminderEnabled *bool `json:"weekly_reminder_enabled" validate:"required"`
}

type PreferencesResponse struct {
	WeeklyReminderEnabled bool `json:"weekly_reminder_enabled"`
}
