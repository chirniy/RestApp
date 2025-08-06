package domain

import (
	"context"
	"github.com/google/uuid"
)

type User struct {
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Age       int       `json:"age"`
	ID        uuid.UUID `json:"id"`
}

type PartialUser struct {
	Firstname *string `json:"firstname,omitempty"`
	Lastname  *string `json:"lastname,omitempty"`
	Age       *int    `json:"age,omitempty"`
}

type UserInput struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, u UserInput) error
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	Update(ctx context.Context, id uuid.UUID, u PartialUser) error
	Delete(ctx context.Context, id uuid.UUID) error
}
