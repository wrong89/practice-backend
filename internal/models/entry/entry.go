package entry

import (
	"context"
	"time"
)

type Entry struct {
	ID            int       `json:"id"`
	Course        string    `json:"course"`
	Date          time.Time `json:"date"`
	UserID        int       `json:"user_id"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
}

func NewEntry(course string, date time.Time, userID int, paymentMethod string) *Entry {
	return &Entry{
		Course:        course,
		Date:          date,
		UserID:        userID,
		PaymentMethod: paymentMethod,
		Status:        "not processed",
	}
}

func (e *Entry) MarkAsProcessed() *Entry {
	e.Status = "processed"
	return e
}

type EntryRepo interface {
	CreateEntry(
		ctx context.Context,
		course string,
		date time.Time,
		userID int,
		paymentMethod string,
	) (Entry, error)
	GetEntryByID(ctx context.Context, id int) (Entry, error)
	DeleteEntry(ctx context.Context, id int) error
}
