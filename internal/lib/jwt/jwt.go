package jwtt

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	// 1. ИСПРАВЛЕНО: Меняем ES256 на HS256, чтобы подписывать обычной строкой (app.Secret)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id.String() // Убедись, что Id переводится в строку, так как uuid.UUID в JSON может сериализоваться криво
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.Id

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken принимает токен и секретное слово приложения, возвращает UUID пользователя и ошибку
func ValidateToken(tokenString string, secret string) (uuid.UUID, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return uuid.Nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	claims := token.Claims.(jwt.MapClaims)

	uidStr, ok := claims["uid"].(string)
	if !ok {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}
