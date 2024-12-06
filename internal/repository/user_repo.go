package repository

import (
	"context"
	"database/sql"
	"golang-tes/internal/domain"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, name, email, password, role) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}

func (r *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, password, role FROM users WHERE email = ?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, password, role FROM users WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mysqlUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET name = ?, email = ?, password = ?, role = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role, user.ID)
	return err
}
