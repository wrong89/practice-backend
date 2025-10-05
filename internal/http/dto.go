package http

import (
	"encoding/json"
	"errors"
	"net/mail"
	"practice-backend/internal/validation"
	"time"
)

var (
	ErrInvalidOrEmptyID = errors.New("invalid or empty id")
	ErrInvalidStatus    = errors.New("invalid status")

	ErrLoginIsEmpty      = errors.New("login is empty")
	ErrPasswordIsEmpty   = errors.New("password is empty")
	ErrNameIsEmpty       = errors.New("name is empty")
	ErrSurnameIsEmpty    = errors.New("surname is empty")
	ErrPatronymicIsEmpty = errors.New("patronymic is empty")

	ErrCourseIsEmpty        = errors.New("course is empty")
	ErrDateIsEmpty          = errors.New("date is empty")
	ErrInvalidDate          = errors.New("date is invalid")
	ErrPaymentMethodIsEmpty = errors.New("payment_method is empty")
	ErrInvalidOrEmptyUserID = errors.New("user_id invalid or empty")
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

type LoginUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *LoginUserDTO) Validate() error {
	if r.Login == "" {
		return ErrLoginIsEmpty
	}
	if r.Password == "" {
		return ErrPasswordIsEmpty
	}

	return nil
}

type CreateEntryDTO struct {
	Course        string `json:"course"`
	Date          string `json:"date"`
	UserID        int    `json:"user_id"`
	PaymentMethod string `json:"payment_method"`
}

func (c *CreateEntryDTO) Validate() error {
	if c.Course == "" {
		return ErrCourseIsEmpty
	}

	if c.Date == "" {
		return ErrDateIsEmpty
	}

	// "2025-10-05" valid
	if _, err := time.Parse(time.DateOnly, c.Date); err != nil {
		return ErrInvalidDate
	}

	if c.PaymentMethod == "" {
		return ErrPaymentMethodIsEmpty
	}

	if c.UserID <= 0 {
		return ErrInvalidOrEmptyUserID
	}

	return nil
}

type UpdateEntryDTO struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func (u *UpdateEntryDTO) Validate() error {
	if u.ID <= 0 {
		return ErrInvalidOrEmptyID
	}
	if !(u.Status == "not processed" || u.Status == "processed") {
		return ErrInvalidStatus
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
