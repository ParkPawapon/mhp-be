package db

import (
	"time"

	"github.com/google/uuid"
)

type CaregiverAssignment struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PatientID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CaregiverID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Relationship string    `gorm:"size:50;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
