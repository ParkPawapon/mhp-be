package dto

import (
	"time"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type MeResponse struct {
	ID      string          `json:"id"`
	Role    constants.Role  `json:"role"`
	Profile ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	FirstName             string  `json:"first_name"`
	LastName              string  `json:"last_name"`
	HN                    *string `json:"hn,omitempty"`
	CitizenID             *string `json:"citizen_id,omitempty"`
	DateOfBirth           string  `json:"date_of_birth"`
	Gender                *constants.GenderType `json:"gender,omitempty"`
	BloodType             *string `json:"blood_type,omitempty"`
	AddressText           *string `json:"address_text,omitempty"`
	GPSLat                *float64 `json:"gps_lat,omitempty"`
	GPSLong               *float64 `json:"gps_long,omitempty"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
	AvatarURL             *string `json:"avatar_url,omitempty"`
}

type UpdateProfileRequest struct {
	HN                    *string `json:"hn"`
	CitizenID             *string `json:"citizen_id"`
	FirstName             *string `json:"first_name"`
	LastName              *string `json:"last_name"`
	DateOfBirth           *time.Time `json:"date_of_birth"`
	Gender                *constants.GenderType `json:"gender"`
	BloodType             *string `json:"blood_type"`
	AddressText           *string `json:"address_text"`
	GPSLat                *float64 `json:"gps_lat"`
	GPSLong               *float64 `json:"gps_long"`
	EmergencyContactName  *string `json:"emergency_contact_name"`
	EmergencyContactPhone *string `json:"emergency_contact_phone"`
	AvatarURL             *string `json:"avatar_url"`
}

type DeviceTokenRequest struct {
	DeviceToken string `json:"device_token" validate:"required"`
	Platform    string `json:"platform" validate:"required"`
}
