package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string         `gorm:"size:50;uniqueIndex;not null"`
	PasswordHash string         `gorm:"size:255;not null"`
	Role         constants.Role `gorm:"type:role_type;not null;default:PATIENT"`
	IsActive     bool           `gorm:"default:true"`
	IsVerified   bool           `gorm:"default:false"`
	LineUserID   *string        `gorm:"size:100;uniqueIndex"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Profile      UserProfile    `gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	UserID                uuid.UUID             `gorm:"type:uuid;primaryKey"`
	HN                    *string               `gorm:"size:20;uniqueIndex"`
	CitizenID             *string               `gorm:"size:13;uniqueIndex"`
	FirstName             string                `gorm:"size:100;not null"`
	LastName              string                `gorm:"size:100;not null"`
	DateOfBirth           time.Time             `gorm:"type:date;not null"`
	Gender                *constants.GenderType `gorm:"type:gender_type"`
	BloodType             *string               `gorm:"size:5"`
	AddressText           *string               `gorm:"type:text"`
	GPSLat                *float64              `gorm:"type:decimal(10,8)"`
	GPSLong               *float64              `gorm:"type:decimal(11,8)"`
	EmergencyContactName  *string               `gorm:"size:100"`
	EmergencyContactPhone *string               `gorm:"size:15"`
	AvatarURL             *string               `gorm:"type:text"`
}
