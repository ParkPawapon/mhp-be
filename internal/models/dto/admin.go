package dto

import "time"

type StaffLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PatientSummaryResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	HN        *string   `json:"hn,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PatientDetailResponse struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	HN        *string `json:"hn,omitempty"`
	CitizenID *string `json:"citizen_id,omitempty"`
}
