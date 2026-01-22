package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type NotificationTemplate struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code      string         `gorm:"size:50;uniqueIndex;not null"`
	Title     string         `gorm:"size:255;not null"`
	Body      string         `gorm:"type:text;not null"`
	Data      datatypes.JSON `gorm:"type:jsonb"`
	IsActive  bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

type NotificationEvent struct {
	ID           uuid.UUID                    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID                    `gorm:"type:uuid;not null;index"`
	TemplateCode string                       `gorm:"size:50;not null"`
	ScheduledAt  time.Time                    `gorm:"type:timestamptz;not null;index"`
	SentAt       *time.Time                   `gorm:"type:timestamptz"`
	Status       constants.NotificationStatus `gorm:"type:notification_status;not null;default:PENDING"`
	Payload      datatypes.JSON               `gorm:"type:jsonb"`
	CreatedAt    time.Time                    `gorm:"autoCreateTime"`
}

func (NotificationEvent) TableName() string {
	return "notification_events"
}
