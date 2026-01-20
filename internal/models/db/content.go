package db

import (
	"time"

	"github.com/google/uuid"
)

type HealthContent struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title            string    `gorm:"size:255;not null"`
	BodyContent      *string   `gorm:"type:text"`
	ThumbnailURL     *string   `gorm:"type:text"`
	ExternalVideoURL *string   `gorm:"type:text"`
	Category         *string   `gorm:"size:50"`
	IsPublished      bool      `gorm:"default:false"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}

func (HealthContent) TableName() string {
	return "health_content"
}
