package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) domain.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := ur.db.WithContext(ctx).Create(user).Error; err != nil {
		ur.logger.ErrorContext(ctx,
			"failed to create user",
			"email", user.Email,
			"error", err,
		)
		return err
	}

	ur.logger.InfoContext(ctx,
		"user created",
		"user_id", user.ID,
		"email", user.Email,
	)

	return nil
}

func (ur *userRepository) Fetch(ctx context.Context) ([]models.User, error) {
	var users []models.User

	if err := ur.db.WithContext(ctx).Find(&users).Error; err != nil {
		ur.logger.ErrorContext(ctx,
			"failed to fetch users",
			"error", err,
		)
		return nil, err
	}

	return users, nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := ur.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.WarnContext(ctx,
				"user not found",
				"email", email,
			)
			return models.User{}, gorm.ErrRecordNotFound
		}

		ur.logger.ErrorContext(ctx,
			"failed to fetch user by email",
			"email", email,
			"error", err,
		)
		return models.User{}, err
	}

	return user, nil
}

func (ur *userRepository) GetByID(ctx context.Context, id string) (models.User, error) {
	var user models.User

	err := ur.db.WithContext(ctx).
		First(&user, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.WarnContext(ctx,
				"user not found",
				"user_id", id,
			)
			return models.User{}, gorm.ErrRecordNotFound
		}

		ur.logger.ErrorContext(ctx,
			"failed to fetch user by id",
			"user_id", id,
			"error", err,
		)
		return models.User{}, err
	}

	return user, nil
}
