package repositories

import (
	"context"
	"errors"
	"go-kafka/config"
	"go-kafka/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository interface {
	IsEmailExists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, u *models.User) (models.User, error)
	FindAll(ctx *gin.Context) ([]models.User, error)
	FindByID(ctx *gin.Context, id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) Create(ctx context.Context, u *models.User) (models.User, error) {

	if err := r.db.Create(u).Error; err != nil {
		return models.User{}, err
	}

	return *u, nil
}

func (r *userRepository) FindAll(ctx *gin.Context) ([]models.User, error) {
	var users []models.User
	if err := config.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(ctx *gin.Context, id uint) (*models.User, error) {
	var u models.User
	if err := config.DB.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	if u.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &u, nil
}
