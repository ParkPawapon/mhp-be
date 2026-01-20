package dto

type CreateCaregiverAssignmentRequest struct {
	PatientID    string `json:"patient_id" validate:"required"`
	CaregiverID  string `json:"caregiver_id" validate:"required"`
	Relationship string `json:"relationship" validate:"required"`
}

type CaregiverAssignmentResponse struct {
	ID           string `json:"id"`
	PatientID    string `json:"patient_id"`
	CaregiverID  string `json:"caregiver_id"`
	Relationship string `json:"relationship"`
}
