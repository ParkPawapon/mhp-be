package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type MedicineService interface {
	ListMaster(ctx context.Context, page, pageSize int) ([]dto.MedicineMasterResponse, int64, error)
	CreatePatientMedicine(ctx context.Context, userID string, req dto.CreatePatientMedicineRequest) (dto.PatientMedicineResponse, error)
	ListPatientMedicines(ctx context.Context, userID string) ([]dto.PatientMedicineResponse, error)
	UpdatePatientMedicine(ctx context.Context, id string, req dto.UpdatePatientMedicineRequest) error
	DeletePatientMedicine(ctx context.Context, id string) error
	CreateSchedule(ctx context.Context, patientMedicineID string, req dto.CreateMedicineScheduleRequest) (dto.MedicineScheduleResponse, error)
	DeleteSchedule(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]dto.MedicineCategoryResponse, error)
	ListCategoryItems(ctx context.Context, categoryID string) ([]dto.MedicineCategoryItemResponse, error)
	GetDosageOptions(ctx context.Context) []string
	GetMealTimingOptions(ctx context.Context) []string
}

type medicineService struct {
	repo   repositories.MedicineRepository
	notify NotificationService
}

func NewMedicineService(repo repositories.MedicineRepository, notify NotificationService) MedicineService {
	return &medicineService{repo: repo, notify: notify}
}

func (s *medicineService) ListMaster(ctx context.Context, page, pageSize int) ([]dto.MedicineMasterResponse, int64, error) {
	items, total, err := s.repo.ListMaster(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]dto.MedicineMasterResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.MedicineMasterResponse{
			ID:          item.ID.String(),
			TradeName:   item.TradeName,
			GenericName: item.GenericName,
			DosageUnit:  item.DosageUnit,
			ImageURL:    item.DefaultImageURL,
		})
	}
	return resp, total, nil
}

func (s *medicineService) CreatePatientMedicine(ctx context.Context, userID string, req dto.CreatePatientMedicineRequest) (dto.PatientMedicineResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return dto.PatientMedicineResponse{}, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	var masterID *uuid.UUID
	if req.MedicineMasterID != nil {
		mid, err := uuid.Parse(strings.TrimSpace(*req.MedicineMasterID))
		if err != nil {
			return dto.PatientMedicineResponse{}, domain.NewError(constants.ValidationFailed, "invalid medicine_master_id")
		}
		if _, err := s.repo.GetMasterByID(ctx, mid); err != nil {
			return dto.PatientMedicineResponse{}, err
		}
		masterID = &mid
	}

	var categoryItemID *uuid.UUID
	var categoryItem *db.MedicineCategoryItem
	if req.CategoryItemID != nil {
		cid, err := uuid.Parse(strings.TrimSpace(*req.CategoryItemID))
		if err != nil {
			return dto.PatientMedicineResponse{}, domain.NewError(constants.ValidationFailed, "invalid category_item_id")
		}
		item, err := s.repo.GetCategoryItemByID(ctx, cid)
		if err != nil {
			return dto.PatientMedicineResponse{}, err
		}
		categoryItemID = &cid
		categoryItem = item
	}

	customName := trimOrNil(req.CustomName)
	if customName == nil && categoryItem != nil {
		name := strings.TrimSpace(categoryItem.DisplayName)
		if name != "" {
			customName = &name
		}
	}

	if masterID == nil && categoryItemID == nil && customName == nil {
		return dto.PatientMedicineResponse{}, domain.NewError(constants.ValidationFailed, "medicine_master_id, category_item_id, or custom_name required")
	}

	dosageAmount := strings.TrimSpace(req.DosageAmount)
	if dosageAmount == "" && categoryItem != nil && categoryItem.DefaultDosageText != nil {
		dosageAmount = strings.TrimSpace(*categoryItem.DefaultDosageText)
	}
	if dosageAmount == "" {
		return dto.PatientMedicineResponse{}, domain.NewError(constants.ValidationFailed, "dosage_amount required")
	}

	med := &db.PatientMedicine{
		UserID:           uid,
		MedicineMasterID: masterID,
		CategoryItemID:   categoryItemID,
		CustomName:       customName,
		DosageAmount:     dosageAmount,
		Instruction:      trimOrNil(req.Instruction),
		Indication:       trimOrNil(req.Indication),
		MyDrugImageURL:   trimOrNil(req.MyDrugImageURL),
		IsActive:         true,
	}

	if err := s.repo.CreatePatientMedicine(ctx, med); err != nil {
		return dto.PatientMedicineResponse{}, err
	}

	return dto.PatientMedicineResponse{
		ID:               med.ID.String(),
		UserID:           med.UserID.String(),
		MedicineMasterID: stringPtr(masterID),
		CategoryItemID:   stringPtr(categoryItemID),
		CustomName:       med.CustomName,
		DosageAmount:     med.DosageAmount,
		Instruction:      med.Instruction,
		Indication:       med.Indication,
		MyDrugImageURL:   med.MyDrugImageURL,
		IsActive:         med.IsActive,
		CreatedAt:        med.CreatedAt,
	}, nil
}

func (s *medicineService) ListPatientMedicines(ctx context.Context, userID string) ([]dto.PatientMedicineResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	items, err := s.repo.ListPatientMedicines(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.PatientMedicineResponse, 0, len(items))
	for _, med := range items {
		resp = append(resp, dto.PatientMedicineResponse{
			ID:               med.ID.String(),
			UserID:           med.UserID.String(),
			MedicineMasterID: stringPtr(med.MedicineMasterID),
			CategoryItemID:   stringPtr(med.CategoryItemID),
			CustomName:       med.CustomName,
			DosageAmount:     med.DosageAmount,
			Instruction:      med.Instruction,
			Indication:       med.Indication,
			MyDrugImageURL:   med.MyDrugImageURL,
			IsActive:         med.IsActive,
			CreatedAt:        med.CreatedAt,
		})
	}
	return resp, nil
}

func (s *medicineService) UpdatePatientMedicine(ctx context.Context, id string, req dto.UpdatePatientMedicineRequest) error {
	medID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid id")
	}

	updates := map[string]any{}
	if req.CustomName != nil {
		updates["custom_name"] = strings.TrimSpace(*req.CustomName)
	}
	if req.DosageAmount != nil {
		value := strings.TrimSpace(*req.DosageAmount)
		if value == "" {
			return domain.NewError(constants.ValidationFailed, "dosage_amount required")
		}
		updates["dosage_amount"] = value
	}
	if req.Instruction != nil {
		updates["instruction"] = trimString(req.Instruction)
	}
	if req.Indication != nil {
		updates["indication"] = trimString(req.Indication)
	}
	if req.MyDrugImageURL != nil {
		updates["my_drug_image_url"] = trimString(req.MyDrugImageURL)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if len(updates) == 0 {
		return domain.NewError(constants.ValidationFailed, "no fields to update")
	}

	return s.repo.UpdatePatientMedicine(ctx, medID, updates)
}

func (s *medicineService) DeletePatientMedicine(ctx context.Context, id string) error {
	medID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid id")
	}
	return s.repo.DeletePatientMedicine(ctx, medID)
}

func (s *medicineService) CreateSchedule(ctx context.Context, patientMedicineID string, req dto.CreateMedicineScheduleRequest) (dto.MedicineScheduleResponse, error) {
	medID, err := uuid.Parse(patientMedicineID)
	if err != nil {
		return dto.MedicineScheduleResponse{}, domain.NewError(constants.ValidationFailed, "invalid patient medicine id")
	}

	medicine, err := s.repo.GetPatientMedicineByID(ctx, medID)
	if err != nil {
		return dto.MedicineScheduleResponse{}, err
	}

	timeSlot, err := time.Parse("15:04", strings.TrimSpace(req.TimeSlot))
	if err != nil {
		return dto.MedicineScheduleResponse{}, domain.NewError(constants.ValidationFailed, "invalid time_slot")
	}

	mealTiming := trimOrNil(req.MealTiming)
	if mealTiming != nil && !isAllowed(*mealTiming, constants.MealTimingOptions) {
		return dto.MedicineScheduleResponse{}, domain.NewError(constants.ValidationFailed, "invalid meal_timing")
	}

	schedule := &db.MedicineSchedule{
		PatientMedicineID: medicine.ID,
		TimeSlot:          timeSlot,
		MealTiming:        mealTiming,
	}

	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		return dto.MedicineScheduleResponse{}, err
	}

	if s.notify != nil {
		_ = s.notify.ScheduleMedicineReminders(ctx, medicine.UserID, schedule.ID, schedule.MealTiming, schedule.TimeSlot)
	}

	return dto.MedicineScheduleResponse{
		ID:                schedule.ID.String(),
		PatientMedicineID: schedule.PatientMedicineID.String(),
		TimeSlot:          schedule.TimeSlot.Format("15:04"),
		MealTiming:        schedule.MealTiming,
		CreatedAt:         schedule.CreatedAt,
	}, nil
}

func (s *medicineService) DeleteSchedule(ctx context.Context, id string) error {
	scheduleID, err := uuid.Parse(id)
	if err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid id")
	}
	return s.repo.DeleteSchedule(ctx, scheduleID)
}

func (s *medicineService) ListCategories(ctx context.Context) ([]dto.MedicineCategoryResponse, error) {
	items, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.MedicineCategoryResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.MedicineCategoryResponse{
			ID:       item.ID.String(),
			Name:     item.Name,
			Code:     item.Code,
			IsActive: item.IsActive,
		})
	}
	return resp, nil
}

func (s *medicineService) ListCategoryItems(ctx context.Context, categoryID string) ([]dto.MedicineCategoryItemResponse, error) {
	cid, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, domain.NewError(constants.ValidationFailed, "invalid category id")
	}
	items, err := s.repo.ListCategoryItems(ctx, cid)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.MedicineCategoryItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.MedicineCategoryItemResponse{
			ID:                item.ID.String(),
			CategoryID:        item.CategoryID.String(),
			DisplayName:       item.DisplayName,
			DefaultDosageText: item.DefaultDosageText,
			IsActive:          item.IsActive,
		})
	}
	return resp, nil
}

func (s *medicineService) GetDosageOptions(ctx context.Context) []string {
	return constants.DosageOptions
}

func (s *medicineService) GetMealTimingOptions(ctx context.Context) []string {
	return constants.MealTimingOptions
}

func trimOrNil(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func trimString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
