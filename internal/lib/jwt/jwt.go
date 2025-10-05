package jwt

import (
	"practice-backend/internal/models/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user user.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["surname"] = user.Surname
	claims["patronymic"] = user.Patronymic
	claims["phone"] = user.Phone
	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte("TEST_SECRET"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
