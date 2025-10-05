package user

import "context"

type User struct {
	ID    int
	Login string
	// PassHash
	Password   string
	Name       string
	Surname    string
	Patronymic string
	Phone      string
	Email      string
	IsAdmin    bool
}

func NewUser(
	login,
	password,
	name,
	surname,
	patronymic,
	phone,
	email string,
	isAdmin bool,
) *User {
	return &User{
		Login:      login,
		Password:   password,
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Phone:      phone,
		Email:      email,
	}
}

// Actions with user
type UserRepo interface {
	CreateUser(
		ctx context.Context,
		login string,
		password string,
		name string,
		surname string,
		patronymic string,
		phone string,
		email string,
		isAdmin bool,
	) (User, error)
	GetUserByID(ctx context.Context, id int) (User, error)
	GetUserByLogin(ctx context.Context, login string) (User, error)
}
