package dto

import "time"

type CreateHealthRecordRequest struct {
	RecordDate  string   `json:"record_date" validate:"required"`
	TimePeriod  string   `json:"time_period" validate:"required"`
	SystolicBP  *int     `json:"systolic_bp"`
	DiastolicBP *int     `json:"diastolic_bp"`
	PulseRate   *int     `json:"pulse_rate"`
	WeightKG    *float64 `json:"weight_kg"`
}

type HealthRecordResponse struct {
	ID         string    `json:"id"`
	RecordDate string    `json:"record_date"`
	TimePeriod string    `json:"time_period"`
	SystolicBP *int      `json:"systolic_bp,omitempty"`
	DiastolicBP *int     `json:"diastolic_bp,omitempty"`
	PulseRate  *int      `json:"pulse_rate,omitempty"`
	WeightKG   *float64  `json:"weight_kg,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateDailyAssessmentRequest struct {
	LogDate         string  `json:"log_date" validate:"required"`
	ExerciseMinutes *int    `json:"exercise_minutes"`
	SleepQuality    *string `json:"sleep_quality"`
	StressLevel     *int    `json:"stress_level"`
	DietCompliance  *string `json:"diet_compliance"`
	Symptoms        any     `json:"symptoms"`
	Note            *string `json:"note"`
}

type DailyAssessmentResponse struct {
	ID              string    `json:"id"`
	LogDate         string    `json:"log_date"`
	ExerciseMinutes int       `json:"exercise_minutes"`
	SleepQuality    *string   `json:"sleep_quality,omitempty"`
	StressLevel     *int      `json:"stress_level,omitempty"`
	DietCompliance  *string   `json:"diet_compliance,omitempty"`
	Symptoms        any       `json:"symptoms,omitempty"`
	Note            *string   `json:"note,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
