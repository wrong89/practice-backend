package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"practice-backend/internal/services/auth"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errDTO := NewErrorDTO(auth.ErrInvalidToken)
			http.Error(w, errDTO.String(), http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		if tokenString == "" {
			errDTO := NewErrorDTO(auth.ErrInvalidToken)
			http.Error(w, errDTO.String(), http.StatusUnauthorized)
			return
		}

		token, err := verifyToken(tokenString)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) {
				errDTO := NewErrorDTO(auth.ErrInvalidToken)

				http.Error(w, errDTO.String(), http.StatusUnauthorized)
				return
			}

			errDTO := NewErrorDTO(err)
			http.Error(w, errDTO.String(), http.StatusUnauthorized)
			return
		}

		fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)
		next.ServeHTTP(w, r)
	})
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte("TEST_SECRET"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token, nil
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
