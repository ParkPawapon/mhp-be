package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
	"github.com/ParkPawapon/mhp-be/internal/models/db"
)

type UserRepository interface {
	Create(ctx context.Context, user *db.User) error
	FindByUsername(ctx context.Context, username string) (*db.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*db.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(dbConn *gorm.DB) UserRepository {
	return &userRepository{db: dbConn}
}

func (r *userRepository) Create(ctx context.Context, user *db.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.NewError(constants.UserConflict, "user already exists")
		}
		return domain.WrapError(constants.InternalError, "create user failed", err)
	}
	return nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*db.User, error) {
	var user db.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.UserNotFound, "user not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find user failed", err)
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	var user db.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.NewError(constants.UserNotFound, "user not found")
		}
		return nil, domain.WrapError(constants.InternalError, "find user failed", err)
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if err := r.db.WithContext(ctx).Model(&db.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error; err != nil {
		return domain.WrapError(constants.InternalError, "update password failed", err)
	}
	return nil
}
