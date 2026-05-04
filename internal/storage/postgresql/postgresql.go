package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgresql", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO customers(email, pass_hash) VALUES($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, storage.ErrUserExists
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM customers WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// App returns app by id.
func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.Id, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM customers WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
