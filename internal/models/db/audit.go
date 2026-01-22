package db

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ActorID      *uuid.UUID `gorm:"type:uuid"`
	TargetUserID *uuid.UUID `gorm:"type:uuid"`
	ActionType   string     `gorm:"size:50;not null"`
	IPAddress    *string    `gorm:"size:45"`
	UserAgent    *string    `gorm:"type:text"`
	Timestamp    time.Time  `gorm:"column:timestamp;autoCreateTime"`
}
