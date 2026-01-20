package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/middleware"
	"github.com/ParkPawapon/mhp-be/internal/services"
	"github.com/ParkPawapon/mhp-be/internal/utils"
)

func bindAndValidateJSON(c *gin.Context, dst any) error {
	if err := c.ShouldBindJSON(dst); err != nil {
		return domain.NewError(constants.ValidationFailed, "invalid json")
	}
	if err := utils.Validate.Struct(dst); err != nil {
		details := utils.ValidationErrors(err)
		return domain.WithDetails(domain.NewError(constants.ValidationFailed, "validation failed"), details)
	}
	return nil
}

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 20

	if v := c.Query("page"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if v := c.Query("page_size"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}
	return page, pageSize
}

func authorizePatientAccess(c *gin.Context, caregiverSvc services.CaregiverService, targetUserID string) (string, error) {
	actorID, _ := middleware.GetActorID(c)
	role, _ := middleware.GetRole(c)

	if targetUserID == "" {
		if role == constants.RoleCaregiver {
			return "", domain.NewError(constants.ValidationFailed, "user_id required")
		}
		targetUserID = actorID.String()
	}

	switch role {
	case constants.RolePatient:
		if targetUserID != actorID.String() {
			return "", domain.NewError(constants.AuthForbidden, "forbidden")
		}
	case constants.RoleCaregiver:
		if caregiverSvc == nil {
			return "", domain.NewError(constants.AuthForbidden, "forbidden")
		}
		pid, err := uuid.Parse(targetUserID)
		if err != nil {
			return "", domain.NewError(constants.ValidationFailed, "invalid user_id")
		}
		assigned, err := caregiverSvc.IsAssigned(c.Request.Context(), actorID, pid)
		if err != nil {
			return "", err
		}
		if !assigned {
			return "", domain.NewError(constants.AuthForbidden, "forbidden")
		}
	}

	return targetUserID, nil
}
