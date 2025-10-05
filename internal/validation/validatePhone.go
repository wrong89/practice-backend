package validation

import (
	"errors"
)

var (
	ErrPhoneIsEmpty = errors.New("phone number is empty")
	ErrInvalidPhone = errors.New("invalid phone number")
)

func ValidatePhone(phone string) error {
	if phone == "" {
		return ErrPhoneIsEmpty
	}

	if len(phone) != 11 {
		return ErrInvalidPhone
	}

	return nil
}
