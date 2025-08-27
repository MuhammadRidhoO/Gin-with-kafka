package service

import (
	"context"
	"fmt"
	"go-kafka/internal/dto/request"
	"go-kafka/internal/dto/response"
	"go-kafka/internal/kafka"
	"go-kafka/internal/models"
	"go-kafka/internal/repositories"

	"github.com/gin-gonic/gin"
)

type UserUsecase interface {
	Create(ctx context.Context, u *request.Request_User) (response.Response_User, error)
	List(ctx *gin.Context) ([]models.User, error)
	Get(ctx *gin.Context, id uint) (*models.User, error)
}

type userUsecase struct {
	repo repositories.UserRepository
	pub  kafka.Producer
}

func NewUserUsecase(repo repositories.UserRepository, pub kafka.Producer) UserUsecase {
	return &userUsecase{repo: repo, pub: pub}
}

func (uc *userUsecase) Create(ctx context.Context, u *request.Request_User) (response.Response_User, error) {
	exists, err := uc.repo.IsEmailExists(ctx, u.Email)
	if err != nil {
		return response.Response_User{}, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return response.Response_User{}, fmt.Errorf("email already registered")
	}

	userModel := models.User{
		Name:  u.Name,
		Email: u.Email,
	}

	savedUser, err := uc.repo.Create(ctx, &userModel)
	if err != nil {
		return response.Response_User{}, fmt.Errorf("failed to create user: %w", err)
	}

	responseUser := response.Response_User{
		ID:    savedUser.ID,
		Name:  savedUser.Name,
		Email: savedUser.Email,
	}

	if err := uc.pub.Publish(ctx, "Request Created User", responseUser); err != nil {
		return responseUser, err
	}

	return responseUser, nil
}

func (uc *userUsecase) List(ctx *gin.Context) ([]models.User, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *userUsecase) Get(ctx *gin.Context, id uint) (*models.User, error) {
	return uc.repo.FindByID(ctx, id)
}
