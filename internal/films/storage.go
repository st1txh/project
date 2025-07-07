package films

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

func NewFilmStorage(pool *pgxpool.Pool, logger *logging.Logger) *Storage {
	return &Storage{
		client: pool,
		logger: logger,
	}
}

func (s *Storage) Create(ctx context.Context, film Film) error {
	q := `
        INSERT INTO films (film_id, title, description, rating, release_date, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := s.client.Exec(
		ctx,
		q,
		film.ID,
		film.Title,
		film.Description,
		film.Rating,
		film.ReleaseDate,
		film.CreatedAt,
		film.UpdatedAt,
	)
	return err
}

func (s *Storage) FindOne(ctx context.Context, id string) ([]Film, error) {
	q := `
       		SELECT films.*
			FROM films
			LEFT JOIN user_film ON films.film_id = user_film.film_id
			WHERE user_film.user_id = $1;
    `

	rows, err := s.client.Query(ctx, q, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warnf("No films found for user: %s", id)
			return nil, nil // или можно вернуть пустой слайс
		}
		s.logger.Errorf("Failed to get films for user: %v", err)
		return nil, fmt.Errorf("failed to get films for user: %w", err)
	}
	defer rows.Close()

	var userFilms []Film
	for rows.Next() {
		var film Film
		if err := rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.Rating,
			&film.ReleaseDate,
			&film.CreatedAt,
			&film.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan film: %w", err)
		}
		userFilms = append(userFilms, film)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return userFilms, nil
}

func (s *Storage) FindAll(ctx context.Context) ([]Film, error) {
	q := `SELECT film_id, title, description, rating, release_date, created_at, updated_at
        FROM films`
	rows, err := s.client.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of films: %w", err)
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		if err := rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.Rating,
			&film.ReleaseDate,
			&film.CreatedAt,
			&film.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		films = append(films, film)
	}
	return films, nil
}

func (s *Storage) FindAllSort(ctx context.Context) ([]Film, error) {
	q := `  SELECT film_id, title, description, rating, release_date, created_at, updated_at
        	FROM films
        	ORDER BY title, rating, release_date 
		  `
	rows, err := s.client.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of films: %w", err)
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		if err := rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.Rating,
			&film.ReleaseDate,
			&film.CreatedAt,
			&film.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		films = append(films, film)
	}
	return films, nil
}

func (s *Storage) PartialUpdate(ctx context.Context, id string, input UpdateFilm) error {
	q := `
        UPDATE films 
        SET 
            title = COALESCE($2, title),
            description = COALESCE($3, description),
            rating = COALESCE($4, rating),
            release_date = COALESCE($5, release_date),
            updated_at = NOW()
        WHERE film_id = $1
    `
	_, err := s.client.Exec(ctx, q, id, input.Title, input.Description, input.Rating, input.ReleaseDate)
	return err
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM films WHERE film_id = $1`

	_, err := s.client.Exec(ctx, q, id)
	if err != nil {
		s.logger.Errorf("Failed to delete film: %v", err)
		return fmt.Errorf("failed to delete film: %w", err)
	}
	return nil
}
