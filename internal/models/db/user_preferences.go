package db

import (
	"time"

	"github.com/google/uuid"
)

type UserPreference struct {
	UserID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	WeeklyReminderEnabled bool      `gorm:"default:true"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime"`
}

func (UserPreference) TableName() string {
	return "user_preferences"
}
