package db

import (
	"time"

	"github.com/google/uuid"
)

type MedicineCategory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Code      string    `gorm:"size:50;uniqueIndex;not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (MedicineCategory) TableName() string {
	return "medicine_categories"
}

type MedicineCategoryItem struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID        uuid.UUID `gorm:"type:uuid;not null;index"`
	DisplayName       string    `gorm:"size:255;not null"`
	DefaultDosageText *string   `gorm:"size:100"`
	IsActive          bool      `gorm:"default:true"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
}

func (MedicineCategoryItem) TableName() string {
	return "medicine_category_items"
}
