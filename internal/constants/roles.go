package constants

type Role string

const (
	RolePatient   Role = "PATIENT"
	RoleNurse     Role = "NURSE"
	RoleAdmin     Role = "ADMIN"
	RoleCaregiver Role = "CAREGIVER"
)

func (r Role) IsValid() bool {
	switch r {
	case RolePatient, RoleNurse, RoleAdmin, RoleCaregiver:
		return true
	default:
		return false
	}
}

func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}
