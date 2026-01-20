package constants

type GenderType string

type MedIntakeStatus string

type AppointmentCategory string

type AppointmentStatus string

const (
	GenderMale   GenderType = "MALE"
	GenderFemale GenderType = "FEMALE"
	GenderOther  GenderType = "OTHER"
)

const (
	MedTaken   MedIntakeStatus = "TAKEN"
	MedMissed  MedIntakeStatus = "MISSED"
	MedSkipped MedIntakeStatus = "SKIPPED"
)

const (
	ApptHospital  AppointmentCategory = "HOSPITAL"
	ApptHomeVisit AppointmentCategory = "HOME_VISIT"
)

const (
	ApptPending   AppointmentStatus = "PENDING"
	ApptConfirmed AppointmentStatus = "CONFIRMED"
	ApptCompleted AppointmentStatus = "COMPLETED"
	ApptCancelled AppointmentStatus = "CANCELLED"
)
