package dto

import (
	"time"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type CreateAppointmentRequest struct {
	Title        string                        `json:"title" validate:"required"`
	ApptType     constants.AppointmentCategory `json:"appt_type" validate:"required"`
	ApptDateTime string                        `json:"appt_datetime" validate:"required"`
	LocationName *string                       `json:"location_name"`
	SlipImageURL *string                       `json:"slip_image_url"`
}

type AppointmentResponse struct {
	ID           string                        `json:"id"`
	UserID       string                        `json:"user_id"`
	Title        string                        `json:"title"`
	ApptType     constants.AppointmentCategory `json:"appt_type"`
	ApptDateTime time.Time                     `json:"appt_datetime"`
	LocationName *string                       `json:"location_name,omitempty"`
	SlipImageURL *string                       `json:"slip_image_url,omitempty"`
	Status       constants.AppointmentStatus   `json:"status"`
	CreatedAt    time.Time                     `json:"created_at"`
}

type UpdateAppointmentStatusRequest struct {
	Status constants.AppointmentStatus `json:"status" validate:"required"`
}

type CreateNurseVisitNoteRequest struct {
	VisitDetails      string  `json:"visit_details" validate:"required"`
	VitalSignsSummary any     `json:"vital_signs_summary"`
	NextActionPlan    *string `json:"next_action_plan"`
}
