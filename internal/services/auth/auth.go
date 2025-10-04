package auth

import "practice-backend/internal/models/user"

type Auth struct {
	userRepo user.UserRepo
}

func NewAuth(userRepo user.UserRepo) *Auth {
	return &Auth{
		userRepo: userRepo,
	}
}

func (a *Auth) Register() {}
func (a *Auth) Login()    {}
