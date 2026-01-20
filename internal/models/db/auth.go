package db

import (
	"time"

	"github.com/google/uuid"
)

type AuthOtpCode struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PhoneNumber string    `gorm:"size:20;not null;index"`
	OtpCode     string    `gorm:"size:10;not null"`
	RefCode     string    `gorm:"size:10;not null"`
	ExpiredAt   time.Time `gorm:"not null"`
	IsUsed      bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
