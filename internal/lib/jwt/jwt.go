package jwtt

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/voznikaetnepriyazn/Autorization_service/internal/domain/models"
)

func NewToken(user models.User, app models.App, duration time.Duration, secret string) (string, error) {
	// Использовали HS256 для работы с обычными строковыми секретами
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	// ВАЖНО: приводим Id к строке. Если у тебя в базе Id это int, используй strconv.Itoa(int(user.Id))
	// Если Id это UUID, то user.Id.String()
	claims["uid"] = fmt.Sprintf("%v", user.Id)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.Id

	// Подписываем переданным секретом
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err // Честно возвращаем ошибку наружу
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Защита: проверяем, что метод подписи именно HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	// БЕЗОПАСНОЕ ИЗВЛЕЧЕНИЕ uid:
	// Так как библиотека jwt может восстановить строку как интерфейс,
	// берем через fmt.Sprintf, чтобы гарантированно получить строку и не поймать панику.
	uidRaw, exists := claims["uid"]
	if !exists {
		return uuid.Nil, fmt.Errorf("uid claim not found in token")
	}
	uidStr := fmt.Sprintf("%v", uidRaw)

	// Парсим строку в UUID.
	// ВНИМАНИЕ: Если у тебя в базе ID пользователей — это обычные числа (1, 2, 3),
	// то метод uuid.Parse упадет. Этот метод нужен ТОЛЬКО если в БД юзер имеет тип UUID.
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	return uid, nil
}
