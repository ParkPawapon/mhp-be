package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type CaregiverRepository interface {
	CreateAssignment(ctx context.Context, assignment *db.CaregiverAssignment) error
	ListAssignmentsByPatient(ctx context.Context, patientID uuid.UUID) ([]db.CaregiverAssignment, error)
	IsAssigned(ctx context.Context, caregiverID, patientID uuid.UUID) (bool, error)
}

type caregiverRepository struct {
	db *gorm.DB
}

func NewCaregiverRepository(dbConn *gorm.DB) CaregiverRepository {
	return &caregiverRepository{db: dbConn}
}

func (r *caregiverRepository) CreateAssignment(ctx context.Context, assignment *db.CaregiverAssignment) error {
	if err := r.db.WithContext(ctx).Create(assignment).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create caregiver assignment failed", err)
	}
	return nil
}

func (r *caregiverRepository) ListAssignmentsByPatient(ctx context.Context, patientID uuid.UUID) ([]db.CaregiverAssignment, error) {
	var items []db.CaregiverAssignment
	if err := r.db.WithContext(ctx).Where("patient_id = ?", patientID).Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list caregiver assignments failed", err)
	}
	return items, nil
}

func (r *caregiverRepository) IsAssigned(ctx context.Context, caregiverID, patientID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&db.CaregiverAssignment{}).
		Where("caregiver_id = ? AND patient_id = ?", caregiverID, patientID).
		Count(&count).Error; err != nil {
		return false, domain.WrapError(constants.InternalError, "check caregiver assignment failed", err)
	}
	return count > 0, nil
}
