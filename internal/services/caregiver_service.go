package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
)

type CaregiverService interface {
	CreateAssignment(ctx context.Context, req dto.CreateCaregiverAssignmentRequest) (dto.CaregiverAssignmentResponse, error)
	ListAssignments(ctx context.Context, patientID uuid.UUID) ([]dto.CaregiverAssignmentResponse, error)
	IsAssigned(ctx context.Context, caregiverID, patientID uuid.UUID) (bool, error)
}

type caregiverService struct {
	repo repositories.CaregiverRepository
}

func NewCaregiverService(repo repositories.CaregiverRepository) CaregiverService {
	return &caregiverService{repo: repo}
}

func (s *caregiverService) CreateAssignment(ctx context.Context, req dto.CreateCaregiverAssignmentRequest) (dto.CaregiverAssignmentResponse, error) {
	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		return dto.CaregiverAssignmentResponse{}, domain.NewError(constants.ValidationFailed, "invalid patient_id")
	}
	caregiverID, err := uuid.Parse(req.CaregiverID)
	if err != nil {
		return dto.CaregiverAssignmentResponse{}, domain.NewError(constants.ValidationFailed, "invalid caregiver_id")
	}

	assignment := &db.CaregiverAssignment{
		PatientID:    patientID,
		CaregiverID:  caregiverID,
		Relationship: req.Relationship,
	}
	if err := s.repo.CreateAssignment(ctx, assignment); err != nil {
		return dto.CaregiverAssignmentResponse{}, err
	}

	return dto.CaregiverAssignmentResponse{
		ID:           assignment.ID.String(),
		PatientID:    assignment.PatientID.String(),
		CaregiverID:  assignment.CaregiverID.String(),
		Relationship: assignment.Relationship,
	}, nil
}

func (s *caregiverService) ListAssignments(ctx context.Context, patientID uuid.UUID) ([]dto.CaregiverAssignmentResponse, error) {
	items, err := s.repo.ListAssignmentsByPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.CaregiverAssignmentResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, dto.CaregiverAssignmentResponse{
			ID:           item.ID.String(),
			PatientID:    item.PatientID.String(),
			CaregiverID:  item.CaregiverID.String(),
			Relationship: item.Relationship,
		})
	}
	return resp, nil
}

func (s *caregiverService) IsAssigned(ctx context.Context, caregiverID, patientID uuid.UUID) (bool, error) {
	return s.repo.IsAssigned(ctx, caregiverID, patientID)
}
