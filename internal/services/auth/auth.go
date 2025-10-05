package auth

import (
	"context"
	"errors"
	"practice-backend/internal/lib/jwt"
	"practice-backend/internal/models/user"
	"practice-backend/internal/storage/inmem"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Auth struct {
	userRepo user.UserRepo
}

func NewAuth(userRepo user.UserRepo) *Auth {
	return &Auth{
		userRepo: userRepo,
	}
}

func (a *Auth) Register(
	ctx context.Context,
	user user.User,
) (userID int, err error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	newUser, err := a.userRepo.CreateUser(ctx,
		user.Login,
		string(passHash),
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Phone,
		user.Email,
		user.IsAdmin,
	)
	if err != nil {
		return -1, err
	}

	return newUser.ID, nil
}

func (a *Auth) Login(
	ctx context.Context,
	login string,
	password string,
) (token string, err error) {
	user, err := a.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return jwt.NewToken(user, time.Hour*24)
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int,
) (bool, error) {
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, inmem.ErrUserNotFound) {
			return false, ErrInvalidCredentials
		}
		return false, ErrInvalidCredentials
	}

	if user.IsAdmin {
		return true, nil
	}
	return false, nil
}

func (a *Auth) CreateAdminUser(ctx context.Context, login string, password string) error {
	var admin user.User
	admin.IsAdmin = true
	admin.Login = login
	admin.Password = password

	_, err := a.Register(ctx, admin)

	return err
}
