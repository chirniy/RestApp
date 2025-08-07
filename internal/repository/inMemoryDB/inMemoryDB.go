package inMemoryDB

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"serversTest2/internal/domain"
)

type InMemoryRepo struct {
	db map[uuid.UUID]domain.User
}

func NewInMemoryRepo(db map[uuid.UUID]domain.User) *InMemoryRepo {
	return &InMemoryRepo{db: db}
}

func (r *InMemoryRepo) CreateUser(ctx context.Context, u domain.UserInput) error {
	id := uuid.New()
	var user domain.User
	user.ID = id
	user.Firstname = u.Firstname
	user.Lastname = u.Lastname
	user.Age = u.Age
	user.Email = u.Email
	user.Password = u.Password
	r.db[id] = user
	return nil
}

func (r *InMemoryRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	for _, user := range r.db {
		users = append(users, user)
	}
	return users, nil
}

func (r *InMemoryRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, ok := r.db[id]
	if !ok {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryRepo) Update(ctx context.Context, id uuid.UUID, u domain.PartialUser) error {
	user := r.db[id]
	if u.Firstname != nil {
		user.Firstname = *u.Firstname
	}
	if u.Lastname != nil {
		user.Lastname = *u.Lastname
	}
	if u.Age != nil {
		user.Age = *u.Age
	}
	if u.Email != nil {
		user.Email = *u.Email
	}
	if u.Password != nil {
		user.Password = *u.Password
	}
	r.db[id] = user
	return nil
}

func (r *InMemoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, ok := r.db[id]
	if !ok {
		return errors.New("user not found")
	}
	delete(r.db, id)
	return nil
}
