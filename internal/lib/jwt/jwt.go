package jwt

import (
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims) //token convert to map
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.Id

	tokenString, err := token.SignedString([]byte(app.Secret)) //secret must not move to logs
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

// ValidateToken принимает токен и секретное слово приложения, возвращает UUID пользователя и ошибку
func ValidateToken(tokenString string, secret string) (uuid.UUID, error) {
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
