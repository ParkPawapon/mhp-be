package db

import (
	"time"

	"github.com/google/uuid"
)

type SupportChatRequest struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Message       string    `gorm:"type:text;not null"`
	Category      string    `gorm:"size:20;not null"`
	AttachmentURL *string   `gorm:"type:text"`
	Status        string    `gorm:"size:20;default:OPEN"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (SupportChatRequest) TableName() string {
	return "support_chat_requests"
}
