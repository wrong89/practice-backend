package entry

import (
	"context"
	"time"
)

type Entry struct {
	ID        int       `json:"id"`
	Course    string    `json:"course"`
	Date      time.Time `json:"date"`
	UserID    int       `json:"user_id"`
	PaymentID int       `json:"payment_id"`
	Status    string    `json:"status"`
}

func NewEntry(course string, date time.Time, userID, paymentID int) *Entry {
	return &Entry{
		Course:    course,
		Date:      date,
		UserID:    userID,
		PaymentID: paymentID,
		Status:    "not processed",
	}
}

type EntryRepo interface {
	CreateEntry(
		ctx context.Context,
		course string,
		date time.Time,
		UserID int,
		PaymentID int,
	) (Entry, error)
	GetEntryByID(ctx context.Context, id int) (Entry, error)
	DeleteEntry(ctx context.Context, id int) error
}
