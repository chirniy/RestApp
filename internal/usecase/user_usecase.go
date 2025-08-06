package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"serversTest2/internal/domain"
)

type UserUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Create(ctx context.Context, userInput domain.UserInput) error {
	if userInput.Firstname == "" || userInput.Lastname == "" {
		return errors.New("Firstname and Lastname are required")
	}

	if userInput.Age <= 0 {
		return errors.New("Age must be greater than 0")
	}
	return u.repo.CreateUser(ctx, userInput)
}

func (u *UserUsecase) GetAll(ctx context.Context) ([]domain.User, error) {
	return u.repo.GetAll(ctx)
}

func (u *UserUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UserUsecase) Update(ctx context.Context, id uuid.UUID, userPartial domain.PartialUser) error {
	return u.repo.Update(ctx, id, userPartial)
}

func (u *UserUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
