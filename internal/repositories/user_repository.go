package repositories

import (
	"context"
	"errors"
	"go-kafka/config"
	"go-kafka/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindAll(ctx context.Context) ([]models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := config.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	if err := config.DB.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &u, nil
}
