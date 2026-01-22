package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type MedicineRepository interface {
	ListMaster(ctx context.Context, page, pageSize int) ([]db.MedicineMaster, int64, error)
	GetMasterByID(ctx context.Context, id uuid.UUID) (*db.MedicineMaster, error)
	CreatePatientMedicine(ctx context.Context, med *db.PatientMedicine) error
	ListPatientMedicines(ctx context.Context, userID uuid.UUID) ([]db.PatientMedicine, error)
	GetPatientMedicineByID(ctx context.Context, id uuid.UUID) (*db.PatientMedicine, error)
	UpdatePatientMedicine(ctx context.Context, id uuid.UUID, updates map[string]any) error
	DeletePatientMedicine(ctx context.Context, id uuid.UUID) error
	CreateSchedule(ctx context.Context, schedule *db.MedicineSchedule) error
	GetScheduleByID(ctx context.Context, id uuid.UUID) (*db.MedicineSchedule, error)
	DeleteSchedule(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context) ([]db.MedicineCategory, error)
	ListCategoryItems(ctx context.Context, categoryID uuid.UUID) ([]db.MedicineCategoryItem, error)
	GetCategoryItemByID(ctx context.Context, id uuid.UUID) (*db.MedicineCategoryItem, error)
}

type medicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(dbConn *gorm.DB) MedicineRepository {
	return &medicineRepository{db: dbConn}
}

func (r *medicineRepository) ListMaster(ctx context.Context, page, pageSize int) ([]db.MedicineMaster, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&db.MedicineMaster{}).Count(&total).Error; err != nil {
		return nil, 0, domain.WrapError(constants.InternalError, "count medicine master failed", err)
	}
	var items []db.MedicineMaster
	if err := r.db.WithContext(ctx).
		Order("created_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, domain.WrapError(constants.InternalError, "list medicine master failed", err)
	}
	return items, total, nil
}

func (r *medicineRepository) GetMasterByID(ctx context.Context, id uuid.UUID) (*db.MedicineMaster, error) {
	var item db.MedicineMaster
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.MedNotFound, "medicine master not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find medicine master failed", err)
	}
	return &item, nil
}

func (r *medicineRepository) CreatePatientMedicine(ctx context.Context, med *db.PatientMedicine) error {
	if err := r.db.WithContext(ctx).Create(med).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create patient medicine failed", err)
	}
	return nil
}

func (r *medicineRepository) ListPatientMedicines(ctx context.Context, userID uuid.UUID) ([]db.PatientMedicine, error) {
	var items []db.PatientMedicine
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list patient medicines failed", err)
	}
	return items, nil
}

func (r *medicineRepository) GetPatientMedicineByID(ctx context.Context, id uuid.UUID) (*db.PatientMedicine, error) {
	var item db.PatientMedicine
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.MedNotFound, "patient medicine not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find patient medicine failed", err)
	}
	return &item, nil
}

func (r *medicineRepository) UpdatePatientMedicine(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&db.PatientMedicine{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return domain.WrapError(constants.InternalError, "update patient medicine failed", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.NewError(constants.MedNotFound, "patient medicine not found")
	}
	return nil
}

func (r *medicineRepository) DeletePatientMedicine(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&db.PatientMedicine{}, "id = ?", id)
	if result.Error != nil {
		return domain.WrapError(constants.InternalError, "delete patient medicine failed", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.NewError(constants.MedNotFound, "patient medicine not found")
	}
	return nil
}

func (r *medicineRepository) CreateSchedule(ctx context.Context, schedule *db.MedicineSchedule) error {
	if err := r.db.WithContext(ctx).Create(schedule).Error; err != nil {
		return domain.WrapError(constants.InternalError, "create medicine schedule failed", err)
	}
	return nil
}

func (r *medicineRepository) GetScheduleByID(ctx context.Context, id uuid.UUID) (*db.MedicineSchedule, error) {
	var item db.MedicineSchedule
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.MedNotFound, "medicine schedule not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find medicine schedule failed", err)
	}
	return &item, nil
}

func (r *medicineRepository) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&db.MedicineSchedule{}, "id = ?", id)
	if result.Error != nil {
		return domain.WrapError(constants.InternalError, "delete medicine schedule failed", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.NewError(constants.MedNotFound, "medicine schedule not found")
	}
	return nil
}

func (r *medicineRepository) ListCategories(ctx context.Context) ([]db.MedicineCategory, error) {
	var items []db.MedicineCategory
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("created_at asc").Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list medicine categories failed", err)
	}
	return items, nil
}

func (r *medicineRepository) ListCategoryItems(ctx context.Context, categoryID uuid.UUID) ([]db.MedicineCategoryItem, error) {
	var items []db.MedicineCategoryItem
	if err := r.db.WithContext(ctx).
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("created_at asc").
		Find(&items).Error; err != nil {
		return nil, domain.WrapError(constants.InternalError, "list medicine category items failed", err)
	}
	return items, nil
}

func (r *medicineRepository) GetCategoryItemByID(ctx context.Context, id uuid.UUID) (*db.MedicineCategoryItem, error) {
	var item db.MedicineCategoryItem
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.MedNotFound, "medicine category item not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find medicine category item failed", err)
	}
	return &item, nil
}
