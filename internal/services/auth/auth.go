package auth2

//business logic, database conversations

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passHash []byte) (uid uuid.UUID, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID uuid.UUID) (models.App, error)
}

// InitAuth returns a new instanse of Auth service
func InitAuth(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login cheks if user with given credentials exists in the system
//
// If user exists, but password is incorrect, returns error.
// if user doesn't exist, return error.
func (a *Auth) Login(ctx context.Context, email string, password string, appID uuid.UUID) (string, error) {
	panic("not implemented")
}

// RegisterNewUser registers new user in the system and returns user ID.
// if user with given username already exist, return error.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (uuid.UUID, error) {
	panic("not implemented")
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID uuid.UUID) (string, error) {
	panic("not implemented")
}
