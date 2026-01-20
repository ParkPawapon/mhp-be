package dto

import "time"

type MedicineMasterResponse struct {
	ID          string  `json:"id"`
	TradeName   string  `json:"trade_name"`
	GenericName *string `json:"generic_name,omitempty"`
	DosageUnit  string  `json:"dosage_unit"`
	ImageURL    *string `json:"image_url,omitempty"`
}

type CreatePatientMedicineRequest struct {
	MedicineMasterID *string `json:"medicine_master_id"`
	CustomName       *string `json:"custom_name"`
	DosageAmount     string  `json:"dosage_amount" validate:"required"`
	Instruction      *string `json:"instruction"`
	Indication       *string `json:"indication"`
	MyDrugImageURL   *string `json:"my_drug_image_url"`
}

type PatientMedicineResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	MedicineMasterID *string   `json:"medicine_master_id,omitempty"`
	CustomName       *string   `json:"custom_name,omitempty"`
	DosageAmount     string    `json:"dosage_amount"`
	Instruction      *string   `json:"instruction,omitempty"`
	Indication       *string   `json:"indication,omitempty"`
	MyDrugImageURL   *string   `json:"my_drug_image_url,omitempty"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

type UpdatePatientMedicineRequest struct {
	CustomName     *string `json:"custom_name"`
	DosageAmount   *string `json:"dosage_amount"`
	Instruction    *string `json:"instruction"`
	Indication     *string `json:"indication"`
	MyDrugImageURL *string `json:"my_drug_image_url"`
	IsActive       *bool   `json:"is_active"`
}

type CreateMedicineScheduleRequest struct {
	TimeSlot   string  `json:"time_slot" validate:"required"`
	MealTiming *string `json:"meal_timing"`
}

type MedicineScheduleResponse struct {
	ID                string    `json:"id"`
	PatientMedicineID string    `json:"patient_medicine_id"`
	TimeSlot          string    `json:"time_slot"`
	MealTiming        *string   `json:"meal_timing,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}
