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
		token, err := validateToken(w, r)
		if err != nil {
			errDTO := NewErrorDTO(auth.ErrInvalidToken)
			http.Error(w, errDTO.String(), http.StatusUnauthorized)
			return
		}

		fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)
		next.ServeHTTP(w, r)
	})
}

func validateToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, auth.ErrInvalidToken
	}

	tokenString := strings.Split(authHeader, " ")[1]

	if tokenString == "" {
		return nil, auth.ErrInvalidToken
	}

	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	return token, nil
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

func AdminMiddleware(authService Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := validateToken(w, r)
			if err != nil {
				errDTO := NewErrorDTO(auth.ErrInvalidToken)
				http.Error(w, errDTO.String(), http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			userID, ok := claims["uid"].(float64)
			if !ok {
				errDTO := NewErrorDTO(ErrInvalidOrEmptyID)
				http.Error(w, errDTO.String(), http.StatusUnauthorized)
				return
			}

			isAdmin, err := authService.IsAdmin(r.Context(), int(userID))
			if err != nil {
				errDTO := NewErrorDTO(err)
				http.Error(w, errDTO.String(), http.StatusUnauthorized)
				return
			}

			if !isAdmin {
				errDTO := NewErrorDTO(errors.New("user is not admin"))
				http.Error(w, errDTO.String(), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}
