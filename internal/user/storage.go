package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"rest-api-tutorial/pkg/logging"
)

type Storage struct {
	client *pgxpool.Pool
	logger *logging.Logger
}

/*type Repository interface {
	Create(ctx context.Context, user User) (string, error)
	FindAll(ctx context.Context) (u []User, err error)
	FindOne(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}*/

func NewUserStorage(pool *pgxpool.Pool, logger *logging.Logger) *Storage {
	return &Storage{
		client: pool,
		logger: logger,
	}
}

func (s *Storage) Create(ctx context.Context, user User) error {
	q := `
        INSERT INTO users (id, name, email, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := s.client.Exec(
		ctx,
		q,
		user.ID,
		user.Name,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (s *Storage) FindOne(ctx context.Context, id string) (*User, error) {
	q := `
        SELECT id, name, email, created_at, updated_at 
        FROM users 
        WHERE id = $1
    `

	var user User
	err := s.client.QueryRow(ctx, q, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warnf("User not found: %s", id)
			return nil, fmt.Errorf("user not found")
		}
		s.logger.Errorf("Failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *Storage) PartialUpdate(ctx context.Context, id string, input UserUpdate) error {
	q := `
        UPDATE users 
        SET 
            name = COALESCE($2, name),
            email = COALESCE($3, email),
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := s.client.Exec(ctx, q, id, input.Name, input.Email)
	return err
}

func (s *Storage) Update(ctx context.Context, id string, input User) error {
	q := `
        UPDATE users 
        SET 
            name = COALESCE($2, name),
            email = COALESCE($3, email),
            updated_at = NOW()
        WHERE id = $1
    `

	_, err := s.client.Exec(ctx, q, id, input.Name, input.Email)

	if err != nil {
		s.logger.Errorf("Failed to update user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM users WHERE id = $1`

	_, err := s.client.Exec(ctx, q, id)
	if err != nil {
		s.logger.Errorf("Failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *Storage) FindAll(ctx context.Context) ([]User, error) {
	q := `SELECT id, name, email, created_at, updated_at FROM users`
	rows, err := s.client.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
