package http

import (
	"encoding/json"
	"errors"
	"net/mail"
	"practice-backend/internal/validation"
	"time"
)

type HelloWorldDTO struct {
	Title string `json:"title"`
}

var (
	ErrLoginIsEmpty      = errors.New("login is empty")
	ErrPasswordIsEmpty   = errors.New("password is empty")
	ErrNameIsEmpty       = errors.New("name is empty")
	ErrSurnameIsEmpty    = errors.New("surname is empty")
	ErrPatronymicIsEmpty = errors.New("patronymic is empty")
)

type RegisterUserDTO struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

func (r *RegisterUserDTO) Validate() error {
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return err
	}
	if r.Login == "" {
		return ErrLoginIsEmpty
	}
	if r.Name == "" {
		return ErrNameIsEmpty
	}
	if r.Password == "" {
		return ErrPasswordIsEmpty
	}
	if r.Patronymic == "" {
		return ErrPatronymicIsEmpty
	}
	if err := validation.ValidatePhone(r.Phone); err != nil {
		return err
	}
	if r.Surname == "" {
		return ErrSurnameIsEmpty
	}

	return nil
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
}

func (e ErrorDTO) String() string {
	b, _ := json.MarshalIndent(e, "", "    ")

	return string(b)
}
