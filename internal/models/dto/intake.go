package dto

import (
	"time"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type CreateIntakeRequest struct {
	ScheduleID *string                  `json:"schedule_id"`
	TargetDate string                   `json:"target_date" validate:"required"`
	Status     constants.MedIntakeStatus `json:"status" validate:"required"`
	SkipReason *string                  `json:"skip_reason"`
}

type IntakeHistoryResponse struct {
	ID         string                    `json:"id"`
	UserID     string                    `json:"user_id"`
	ScheduleID *string                   `json:"schedule_id,omitempty"`
	TargetDate string                    `json:"target_date"`
	TakenAt    *time.Time                `json:"taken_at,omitempty"`
	Status     constants.MedIntakeStatus `json:"status"`
	SkipReason *string                   `json:"skip_reason,omitempty"`
	CreatedAt  time.Time                 `json:"created_at"`
}
