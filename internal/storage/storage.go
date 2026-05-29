package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
)

var (
	ErrUserExists   = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passHash []byte) (uid uuid.UUID, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID uuid.UUID) (models.App, error)
}
