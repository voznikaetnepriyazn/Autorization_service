package auth2

//business logic, database conversations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/lib/jwt"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/lib/sl"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/storage"
	"golang.org/x/crypto/bcrypt"
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

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
	const op = "auth.Login"

	log := a.log.With(
		slog.String("operation", op),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s; %w", op, err)
	}

	log.Info("user logged in cuccessfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// if user with given username already exist, return error.
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (uuid.UUID, error) {
	const op = "auth.RegisterNewuser"

	log := a.log.With(
		slog.String("operation", op),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //hash with salt
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
	}

	log.Info("user registered")

	return id, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "auth.isadmin"

	log := a.log.With(
		slog.String("operation", op),
		slog.Any("user_id", userID),
	)

	log.Info("checking admin")

	adm, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("failed to check", sl.Err(err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is admin", adm))

	return adm, nil
}
