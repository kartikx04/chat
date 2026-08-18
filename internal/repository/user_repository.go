package repository

import (
	"log/slog"

	"github.com/kartikx04/chat/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) CreateUser(authOId, email, username, picture string) (*models.User, error) {
	user := &models.User{
		AuthOId:  authOId,
		Email:    email,
		Username: username,
		Picture:  picture,
		Role:     "user",
	}

	result := r.db.Create(user)

	if result.Error != nil {
		// r.slog.Error("error creating user", "error", result.Error)
		return nil, result.Error
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserById(id string) (*models.User, error) {
	var user models.User

	result := r.db.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
