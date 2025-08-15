package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"serversTest2/internal/domain"
	"strings"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, u domain.UserInput) error {
	id := uuid.New()
	_, err := r.db.Exec(`INSERT INTO users (id, firstname, lastname, age, email, password) VALUES ($1, $2, $3, $4, $5, $6)`, id, u.Firstname, u.Lastname, u.Age, u.Email, u.Password)
	return err
}

func (r *PostgresUserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.Query(`SELECT id, firstname, lastname, age, email, password FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Age, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, firstname, lastname, age, email, password  FROM users WHERE id = $1", id).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Age, &user.Email, &user.Password)
	return user, err
}

func (r *PostgresUserRepo) Update(ctx context.Context, id uuid.UUID, u domain.PartialUser) error {
	// Динамическая сборка запроса
	fields := []string{}
	args := []interface{}{}
	argId := 1

	if u.Firstname != nil {
		fields = append(fields, fmt.Sprintf("firstname = $%d", argId))
		args = append(args, *u.Firstname)
		argId++
	}
	if u.Lastname != nil {
		fields = append(fields, fmt.Sprintf("lastname = $%d", argId))
		args = append(args, *u.Lastname)
		argId++
	}
	if u.Age != nil {
		fields = append(fields, fmt.Sprintf("age = $%d", argId))
		args = append(args, *u.Age)
		argId++
	}
	if u.Email != nil {
		fields = append(fields, fmt.Sprintf("email = $%d", argId))
		args = append(args, *u.Email)
		argId++
	}
	if u.Password != nil {
		fields = append(fields, fmt.Sprintf("password = $%d", argId))
		args = append(args, *u.Password)
		argId++
	}

	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Добавляем id в конец args
	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(fields, ", "), argId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *PostgresUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, firstname, lastname, age, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Age, &user.Email, &user.Password)
	return user, err
}
