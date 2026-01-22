package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type IntakeService interface {
	CreateIntake(ctx context.Context, userID string, req dto.CreateIntakeRequest) (dto.IntakeHistoryResponse, error)
	ListHistory(ctx context.Context, userID string, from, to string) ([]dto.IntakeHistoryResponse, error)
}

type intakeService struct {
	repo   repositories.IntakeRepository
	notify NotificationService
}

func NewIntakeService(repo repositories.IntakeRepository, notify NotificationService) IntakeService {
	return &intakeService{repo: repo, notify: notify}
}

func (s *intakeService) CreateIntake(ctx context.Context, userID string, req dto.CreateIntakeRequest) (dto.IntakeHistoryResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return dto.IntakeHistoryResponse{}, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	targetDate, err := time.Parse("2006-01-02", req.TargetDate)
	if err != nil {
		return dto.IntakeHistoryResponse{}, domain.NewError(constants.ValidationFailed, "invalid target_date")
	}
	targetDate = targetDate.UTC()

	var scheduleID *uuid.UUID
	if req.ScheduleID != nil {
		id, err := uuid.Parse(*req.ScheduleID)
		if err != nil {
			return dto.IntakeHistoryResponse{}, domain.NewError(constants.ValidationFailed, "invalid schedule_id")
		}
		scheduleID = &id
	}

	var takenAt *time.Time
	if req.Status == constants.MedTaken {
		now := time.Now().UTC()
		takenAt = &now
	}

	record := &db.IntakeHistory{
		UserID:     uid,
		ScheduleID: scheduleID,
		TargetDate: targetDate,
		TakenAt:    takenAt,
		Status:     req.Status,
		SkipReason: req.SkipReason,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		return dto.IntakeHistoryResponse{}, err
	}

	if req.Status == constants.MedTaken && scheduleID != nil && s.notify != nil {
		_ = s.notify.CancelMedicineAfterMealReminder(ctx, uid, *scheduleID, targetDate)
	}

	return dto.IntakeHistoryResponse{
		ID:         record.ID.String(),
		UserID:     record.UserID.String(),
		ScheduleID: stringPtr(scheduleID),
		TargetDate: record.TargetDate.Format("2006-01-02"),
		TakenAt:    record.TakenAt,
		Status:     record.Status,
		SkipReason: record.SkipReason,
		CreatedAt:  record.CreatedAt,
	}, nil
}

func (s *intakeService) ListHistory(ctx context.Context, userID string, from, to string) ([]dto.IntakeHistoryResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.NewError(constants.ValidationFailed, "invalid user_id")
	}

	var fromDate time.Time
	var toDate time.Time
	if from != "" {
		if parsed, err := time.Parse("2006-01-02", from); err == nil {
			fromDate = parsed.UTC()
		} else {
			return nil, domain.NewError(constants.ValidationFailed, "invalid from")
		}
	}
	if to != "" {
		if parsed, err := time.Parse("2006-01-02", to); err == nil {
			toDate = parsed.UTC()
		} else {
			return nil, domain.NewError(constants.ValidationFailed, "invalid to")
		}
	}

	items, err := s.repo.ListHistory(ctx, uid, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.IntakeHistoryResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.IntakeHistoryResponse{
			ID:         item.ID.String(),
			UserID:     item.UserID.String(),
			ScheduleID: stringPtr(item.ScheduleID),
			TargetDate: item.TargetDate.Format("2006-01-02"),
			TakenAt:    item.TakenAt,
			Status:     item.Status,
			SkipReason: item.SkipReason,
			CreatedAt:  item.CreatedAt,
		})
	}
	return resp, nil
}
