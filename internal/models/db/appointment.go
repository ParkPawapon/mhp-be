package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type Appointment struct {
	ID           uuid.UUID                     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID                     `gorm:"type:uuid;not null;index"`
	CreatorID    *uuid.UUID                    `gorm:"type:uuid"`
	Title        string                        `gorm:"size:255;not null"`
	ApptType     constants.AppointmentCategory `gorm:"type:appt_category;not null"`
	ApptDateTime time.Time                     `gorm:"column:appt_datetime;type:timestamptz;not null"`
	LocationName *string                       `gorm:"size:255"`
	SlipImageURL *string                       `gorm:"type:text"`
	Status       constants.AppointmentStatus   `gorm:"type:appt_status_type;default:PENDING"`
	CreatedAt    time.Time                     `gorm:"autoCreateTime"`
	DeletedAt    gorm.DeletedAt                `gorm:"index"`
}

type NurseVisitNote struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AppointmentID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	NurseID           uuid.UUID      `gorm:"type:uuid;not null;index"`
	VisitDetails      string         `gorm:"type:text;not null"`
	VitalSignsSummary datatypes.JSON `gorm:"type:jsonb"`
	NextActionPlan    *string        `gorm:"type:text"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
}
