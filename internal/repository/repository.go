package repository

import (
	"context"
	"github.com/google/uuid"
	"serversTest2/internal/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, u domain.UserInput) error
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, u domain.PartialUser) error
	Delete(ctx context.Context, id uuid.UUID) error
}
