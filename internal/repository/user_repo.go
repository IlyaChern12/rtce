package repository

import (
	"context"
	"database/sql"

	"github.com/IlyaChern12/rtce/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

// создание юзера
func (ur *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`

	_, err := ur.DB.ExecContext(ctx, query, user.Email, user.PasswordHash)
	return err
}

// получение юзера по почте
func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email=$1`

	err := ur.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}