package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/kartikx04/chat/models"
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
		Id:        uuid.New(),
		AuthOId:   authOId,
		Email:     email,
		Username:  username,
		Picture:   picture,
		Role:      "user",
		CreatedAt: time.Now(),
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
