package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicineMaster struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TradeName       string    `gorm:"size:255;not null"`
	GenericName     *string   `gorm:"size:255"`
	DosageUnit      string    `gorm:"size:50;not null"`
	DefaultImageURL *string   `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}

func (MedicineMaster) TableName() string {
	return "medicines_master"
}

type PatientMedicine struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index"`
	MedicineMasterID *uuid.UUID     `gorm:"type:uuid"`
	CustomName       *string        `gorm:"size:255"`
	DosageAmount     string         `gorm:"size:100;not null"`
	Instruction      *string        `gorm:"type:text"`
	Indication       *string        `gorm:"type:text"`
	MyDrugImageURL   *string        `gorm:"type:text"`
	IsActive         bool           `gorm:"default:true"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type MedicineSchedule struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PatientMedicineID uuid.UUID `gorm:"type:uuid;not null;index"`
	TimeSlot          time.Time `gorm:"type:time;not null"`
	MealTiming        *string   `gorm:"size:50"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
}
