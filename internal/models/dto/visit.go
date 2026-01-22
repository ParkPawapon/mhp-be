package dto

import "time"

type VisitHistoryItem struct {
	AppointmentID     string    `json:"appointment_id"`
	VisitNoteID       string    `json:"visit_note_id"`
	ApptDateTime      time.Time `json:"appt_datetime"`
	Title             string    `json:"title"`
	LocationName      *string   `json:"location_name,omitempty"`
	NurseID           string    `json:"nurse_id"`
	VisitDetails      string    `json:"visit_details"`
	VitalSignsSummary any       `json:"vital_signs_summary"`
	NextActionPlan    *string   `json:"next_action_plan,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}
