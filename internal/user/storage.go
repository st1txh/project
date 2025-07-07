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

func NewUserStorage(pool *pgxpool.Pool, logger *logging.Logger) *Storage {
	return &Storage{
		client: pool,
		logger: logger,
	}
}

func (s *Storage) Create(ctx context.Context, user User) error {
	q := `
        INSERT INTO users (id, name, email, date_of_birth, gender, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := s.client.Exec(
		ctx,
		q,
		user.ID,
		user.Name,
		user.Email,
		user.DateOfBirth,
		user.Gender,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if len(user.FilmUUID) > 0 {
		qFilm := `
            INSERT INTO user_film (user_id, film_id)
            VALUES ($1, $2)
        `
		for _, filmID := range user.FilmUUID {
			_, err = s.client.Exec(ctx, qFilm, user.ID, filmID)
			if err != nil {
				return fmt.Errorf("failed to insert user-film relation: %w", err)
			}
		}
	}

	return err
}

func (s *Storage) FindOne(ctx context.Context, id string) (*User, error) {
	q := `
        SELECT id, name, email, date_of_birth, gender, created_at, updated_at 
        FROM users 
        WHERE id = $1
    `

	var user User
	err := s.client.QueryRow(ctx, q, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.DateOfBirth,
		&user.Gender,
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

func (s *Storage) PartialUpdate(ctx context.Context, id string, input Update) error {
	q := `
        UPDATE users 
        SET 
            name = COALESCE($2, name),
            email = COALESCE($3, email),
            date_of_birth = coalesce($4, date_of_birth),
            gender = coalesce($5, gender),
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := s.client.Exec(ctx, q, id, input.Name, input.Email, input.DateOfBirth, input.Gender)
	return err
}

func (s *Storage) Update(ctx context.Context, id string, input User) error {
	q := `
        UPDATE users 
        SET 
            name = COALESCE($2, name),
            email = COALESCE($3, email),
            date_of_birth = coalesce($4, date_of_birth),
            gender = coalesce($5, gender),
            updated_at = NOW()
        WHERE id = $1
    `

	_, err := s.client.Exec(ctx, q, id, input.Name, input.Email, input.DateOfBirth, input.Gender)

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
	q := `SELECT id, name, email, date_of_birth, gender, created_at, updated_at FROM users`
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
			&user.DateOfBirth,
			&user.Gender,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
