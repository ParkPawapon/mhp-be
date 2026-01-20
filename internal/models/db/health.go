package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type HealthRecord struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	RecordDate time.Time `gorm:"type:date;not null"`
	TimePeriod string    `gorm:"size:20;not null"`
	SystolicBP *int      `gorm:"type:int"`
	DiastolicBP *int     `gorm:"type:int"`
	PulseRate  *int      `gorm:"type:int"`
	WeightKG   *float64  `gorm:"type:decimal(5,2)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type DailyAssessment struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index"`
	LogDate         time.Time      `gorm:"type:date;not null"`
	ExerciseMinutes int            `gorm:"default:0"`
	SleepQuality    *string        `gorm:"size:50"`
	StressLevel     *int           `gorm:"type:int"`
	DietCompliance  *string        `gorm:"size:50"`
	Symptoms        datatypes.JSON `gorm:"type:jsonb"`
	Note            *string        `gorm:"type:text"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
}
