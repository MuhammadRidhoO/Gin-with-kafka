package service

import (
	"context"
	"go-kafka/internal/kafka"
	"go-kafka/internal/models"
	"go-kafka/internal/repositories"
)

type UserUsecase interface {
	Create(ctx context.Context, u *models.User) error
	List(ctx context.Context) ([]models.User, error)
	Get(ctx context.Context, id uint) (*models.User, error)
}

type userUsecase struct {
	repo repositories.UserRepository
	pub  kafka.Producer
}

func NewUserUsecase(repo repositories.UserRepository, pub kafka.Producer) UserUsecase {
	return &userUsecase{repo: repo, pub: pub}
}

func (uc *userUsecase) Create(ctx context.Context, u *models.User) error {
	if err := uc.repo.Create(ctx, u); err != nil {
		return err
	}

	return uc.pub.Publish(ctx, "user-topic", u)
}

func (uc *userUsecase) List(ctx context.Context) ([]models.User, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *userUsecase) Get(ctx context.Context, id uint) (*models.User, error) {
	return uc.repo.FindByID(ctx, id)
}
