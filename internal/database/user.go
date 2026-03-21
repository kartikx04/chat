package database

import (
	"errors"

	"github.com/kartikx04/chat/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(authOId, email, username, picture string) (*models.User, error) {
	user := &models.User{
		AuthOId:  authOId,
		Email:    email,
		Username: username,
		Picture:  picture,
		Role:     "user",
	}

	result := r.db.Create(user) // ← GORM handles INSERT

	if result.Error != nil {
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

func (r *UserRepository) GetUserByAuthOId(authOId string) (*models.User, error) {
	var user models.User

	result := r.db.Where("auth_o_id = ?", authOId).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetOrCreateUser(authID, email, username, picture string) (*models.User, error) {
	var user models.User

	err := r.db.Where("auth_o_id = ?", authID).First(&user).Error
	if err == nil {
		return &user, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = models.User{
			AuthOId:  authID,
			Email:    email,
			Username: username,
			Picture:  picture,
			Role:     "user",
		}

		if err := r.db.Create(&user).Error; err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, err
}
