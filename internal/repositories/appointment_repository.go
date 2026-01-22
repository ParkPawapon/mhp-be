package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type VisitHistoryRow struct {
	AppointmentID     uuid.UUID
	VisitNoteID       uuid.UUID
	ApptDateTime      time.Time
	Title             string
	LocationName      *string
	NurseID           uuid.UUID
	VisitDetails      string
	VitalSignsSummary []byte
	NextActionPlan    *string
	CreatedAt         time.Time
}

type AppointmentRepository interface {
	ListAppointments(ctx context.Context, userID uuid.UUID) ([]db.Appointment, error)
	CreateAppointment(ctx context.Context, appt *db.Appointment) error
	FindByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status constants.AppointmentStatus) error
	DeleteAppointment(ctx context.Context, id uuid.UUID) error
	CreateNurseVisitNote(ctx context.Context, note *db.NurseVisitNote) error
	ListVisitHistory(ctx context.Context, userID uuid.UUID) ([]VisitHistoryRow, error)
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(dbConn *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: dbConn}
}

func (r *appointmentRepository) ListAppointments(ctx context.Context, userID uuid.UUID) ([]db.Appointment, error) {
	var items []db.Appointment
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("appt_datetime desc").
		Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list appointments failed", err)
	}
	return items, nil
}

func (r *appointmentRepository) CreateAppointment(ctx context.Context, appt *db.Appointment) error {
	if err := r.db.WithContext(ctx).Create(appt).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create appointment failed", err)
	}
	return nil
}

func (r *appointmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error) {
	var appt db.Appointment
	if err := r.db.WithContext(ctx).First(&appt, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.ApptNotFound, "appointment not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find appointment failed", err)
	}
	return &appt, nil
}

func (r *appointmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status constants.AppointmentStatus) error {
	result := r.db.WithContext(ctx).Model(&db.Appointment{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return domain.WrapError(constants.InternalError, "update appointment status failed", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.NewError(constants.ApptNotFound, "appointment not found")
	}
	return nil
}

func (r *appointmentRepository) DeleteAppointment(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&db.Appointment{}, "id = ?", id)
	if result.Error != nil {
		return domain.WrapError(constants.InternalError, "delete appointment failed", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.NewError(constants.ApptNotFound, "appointment not found")
	}
	return nil
}

func (r *appointmentRepository) CreateNurseVisitNote(ctx context.Context, note *db.NurseVisitNote) error {
	if err := r.db.WithContext(ctx).Create(note).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create visit note failed", err)
	}
	return nil
}

func (r *appointmentRepository) ListVisitHistory(ctx context.Context, userID uuid.UUID) ([]VisitHistoryRow, error) {
	var rows []VisitHistoryRow
	if err := r.db.WithContext(ctx).
		Table("nurse_visit_notes").
		Select("nurse_visit_notes.id as visit_note_id, nurse_visit_notes.nurse_id, nurse_visit_notes.visit_details, nurse_visit_notes.vital_signs_summary, nurse_visit_notes.next_action_plan, nurse_visit_notes.created_at, appointments.id as appointment_id, appointments.appt_datetime, appointments.title, appointments.location_name").
		Joins("join appointments on appointments.id = nurse_visit_notes.appointment_id").
		Where("appointments.user_id = ? AND appointments.appt_type = ? AND appointments.deleted_at IS NULL", userID, constants.ApptHomeVisit).
		Order("appointments.appt_datetime desc, nurse_visit_notes.created_at desc").
		Scan(&rows).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list visit history failed", err)
	}
	return rows, nil
}
