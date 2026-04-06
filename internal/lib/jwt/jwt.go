package jwt

import (
	"time"

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
