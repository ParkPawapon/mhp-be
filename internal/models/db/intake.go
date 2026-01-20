package db

import (
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type IntakeHistory struct {
	ID         uuid.UUID             `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID             `gorm:"type:uuid;not null;index"`
	ScheduleID *uuid.UUID            `gorm:"type:uuid;index"`
	TargetDate time.Time             `gorm:"type:date;not null"`
	TakenAt    *time.Time            `gorm:"type:timestamptz"`
	Status     constants.MedIntakeStatus `gorm:"type:med_intake_status;not null"`
	SkipReason *string               `gorm:"type:text"`
	CreatedAt  time.Time             `gorm:"autoCreateTime"`
}

func (IntakeHistory) TableName() string {
	return "intake_history"
}
