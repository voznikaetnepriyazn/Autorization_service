package auth2

//business logic, database conversations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
	jwtt "github.com/voznikaetnepriyazn/Autorization_service/internal/lib/jwt"
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
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) // Добавляем метод для получения пользователя по ID
}

type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
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

// Login checs if user with given credentials exists in the system
//
// If user exists, but password is incorrect, returns error.
// if user doesn't exist, return error.
func (a *Auth) Login(ctx context.Context, email string, password string, appID int32) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		// Твой код обработки ошибки пользователя...
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверка пароля...
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Получаем данные приложения из базы
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	// Генерируем токен.
	// ВАЖНО: Твоя функция NewToken должна принимать app.Id как int32 или внутри конвертировать в int!
	token, err := jwtt.NewToken(user, app, a.tokenTTL)
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
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", sl.Err(err))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
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
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found", sl.Err(err))

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		log.Error("failed to check", sl.Err(err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is admin", adm))

	return adm, nil
}

func (a *Auth) ValidateToken(tokenString string, secret string) (uuid.UUID, error) {
	const op = "auth.ValidateToken"

	if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		tokenString = tokenString[7:]
	}

	// 2. Очищаем от случайных пробелов, переносов строк или лишних кавычек
	tokenString = strings.TrimSpace(tokenString)
	tokenString = strings.Trim(tokenString, "\"") // Убирает лишние кавычки, если фронт передал JSON-строку

	// 3. ДЕБАГ ЛОГ В GO
	segments := strings.Split(tokenString, ".")
	fmt.Printf("[Go Debug] Метод вызвался. Сегментов в JWT: %d, Длина строки: %d\n", len(segments), len(tokenString))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, fmt.Errorf("%s: invalid token: %w", op, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("%s: invalid token claims", op)
	}

	uidStr, ok := claims["uid"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("%s: uid not found or not a string", op)
	}

	// Конвертируем строку из JWT обратно в тип uuid.UUID
	parsedUUID, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: invalid uuid format in token: %w", op, err)
	}

	return parsedUUID, nil
}
