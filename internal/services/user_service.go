package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
	"github.com/ParkPawapon/mhp-be/internal/models/dto"
	"github.com/ParkPawapon/mhp-be/internal/repositories"
	"github.com/ParkPawapon/mhp-be/internal/utils"
)

type UserService interface {
	GetMe(ctx context.Context, actorID uuid.UUID, role constants.Role) (dto.MeResponse, error)
	UpdateProfile(ctx context.Context, actorID uuid.UUID, req dto.UpdateProfileRequest) error
	SaveDeviceToken(ctx context.Context, actorID uuid.UUID, req dto.DeviceTokenRequest) error
	UpdatePreferences(ctx context.Context, actorID uuid.UUID, req dto.UpdatePreferencesRequest) (dto.PreferencesResponse, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	profileRepo repositories.ProfileRepository
	deviceRepo  repositories.DeviceTokenRepository
	prefRepo    repositories.PreferenceRepository
	notify      NotificationService
}

func NewUserService(userRepo repositories.UserRepository, profileRepo repositories.ProfileRepository, deviceRepo repositories.DeviceTokenRepository, prefRepo repositories.PreferenceRepository, notify NotificationService) UserService {
	return &userService{userRepo: userRepo, profileRepo: profileRepo, deviceRepo: deviceRepo, prefRepo: prefRepo, notify: notify}
}

func (s *userService) GetMe(ctx context.Context, actorID uuid.UUID, role constants.Role) (dto.MeResponse, error) {
	user, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return dto.MeResponse{}, err
	}

	var profile *db.UserProfile
	if p, err := s.profileRepo.FindByUserID(ctx, actorID); err == nil {
		profile = p
	} else {
		if appErr, ok := domain.AsAppError(err); !ok || appErr.Code != constants.UserNotFound {
			return dto.MeResponse{}, err
		}
	}

	resp := dto.MeResponse{
		ID:   user.ID.String(),
		Role: user.Role,
		Profile: dto.ProfileResponse{
			FirstName: "",
			LastName:  "",
		},
	}

	if profile != nil {
		resp.Profile.FirstName = profile.FirstName
		resp.Profile.LastName = profile.LastName
		resp.Profile.HN = profile.HN
		resp.Profile.CitizenID = profile.CitizenID
		resp.Profile.DateOfBirth = profile.DateOfBirth.Format("2006-01-02")
		resp.Profile.Gender = profile.Gender
		resp.Profile.BloodType = profile.BloodType
		resp.Profile.AddressText = profile.AddressText
		resp.Profile.GPSLat = profile.GPSLat
		resp.Profile.GPSLong = profile.GPSLong
		resp.Profile.EmergencyContactName = profile.EmergencyContactName
		resp.Profile.EmergencyContactPhone = profile.EmergencyContactPhone
		resp.Profile.AvatarURL = profile.AvatarURL
	}

	if role != constants.RoleAdmin {
		if resp.Profile.CitizenID != nil {
			masked := utils.MaskCitizenID(*resp.Profile.CitizenID)
			resp.Profile.CitizenID = &masked
		}
	}
	if role == constants.RoleCaregiver {
		resp.Profile.CitizenID = nil
	}

	return resp, nil
}

func (s *userService) UpdateProfile(ctx context.Context, actorID uuid.UUID, req dto.UpdateProfileRequest) error {
	profile, err := s.profileRepo.FindByUserID(ctx, actorID)
	isNew := false
	if err != nil {
		appErr, ok := domain.AsAppError(err)
		if !ok || appErr.Code != constants.UserNotFound {
			return err
		}
		profile = &db.UserProfile{UserID: actorID}
		isNew = true
	}

	if req.FirstName != nil {
		profile.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		profile.LastName = *req.LastName
	}
	if req.HN != nil {
		profile.HN = req.HN
	}
	if req.CitizenID != nil {
		profile.CitizenID = req.CitizenID
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth.UTC()
	}
	if req.Gender != nil {
		profile.Gender = req.Gender
	}
	if req.BloodType != nil {
		profile.BloodType = req.BloodType
	}
	if req.AddressText != nil {
		profile.AddressText = req.AddressText
	}
	if req.GPSLat != nil {
		profile.GPSLat = req.GPSLat
	}
	if req.GPSLong != nil {
		profile.GPSLong = req.GPSLong
	}
	if req.EmergencyContactName != nil {
		profile.EmergencyContactName = req.EmergencyContactName
	}
	if req.EmergencyContactPhone != nil {
		profile.EmergencyContactPhone = req.EmergencyContactPhone
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = req.AvatarURL
	}

	if isNew {
		if profile.DateOfBirth.IsZero() || profile.FirstName == "" || profile.LastName == "" {
			return domain.NewError(constants.ValidationFailed, "first_name, last_name, and date_of_birth required")
		}
	}

	return s.profileRepo.Upsert(ctx, profile)
}

func (s *userService) SaveDeviceToken(ctx context.Context, actorID uuid.UUID, req dto.DeviceTokenRequest) error {
	platform := strings.ToLower(strings.TrimSpace(req.Platform))
	if platform == "" {
		return domain.NewError(constants.ValidationFailed, "platform required")
	}
	if strings.TrimSpace(req.DeviceToken) == "" {
		return domain.NewError(constants.ValidationFailed, "device_token required")
	}

	if s.deviceRepo == nil {
		return domain.NewError(constants.InternalError, "device token repository not configured")
	}

	record := &db.DeviceToken{
		UserID:   actorID,
		Platform: platform,
		Token:    req.DeviceToken,
		IsActive: true,
	}
	return s.deviceRepo.Save(ctx, record)
}

func (s *userService) UpdatePreferences(ctx context.Context, actorID uuid.UUID, req dto.UpdatePreferencesRequest) (dto.PreferencesResponse, error) {
	if req.WeeklyReminderEnabled == nil {
		return dto.PreferencesResponse{}, domain.NewError(constants.ValidationFailed, "weekly_reminder_enabled required")
	}
	if s.prefRepo == nil {
		return dto.PreferencesResponse{}, domain.NewError(constants.InternalError, "preferences repository not configured")
	}

	pref := &db.UserPreference{UserID: actorID, WeeklyReminderEnabled: *req.WeeklyReminderEnabled}
	if err := s.prefRepo.Upsert(ctx, pref); err != nil {
		return dto.PreferencesResponse{}, err
	}

	if s.notify != nil && !*req.WeeklyReminderEnabled {
		_ = s.notify.CancelWeeklyReminders(ctx, actorID)
	}

	return dto.PreferencesResponse{WeeklyReminderEnabled: pref.WeeklyReminderEnabled}, nil
}
