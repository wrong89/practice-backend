package entry

import (
	"context"
	"time"
)

type Entry struct {
	ID            int
	Course        string
	Date          time.Time
	UserID        int
	PaymentMethod string
	Status        string
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

func (e *Entry) UpdateStatus(status string) *Entry {
	e.Status = status
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
	GetEntries(ctx context.Context) ([]Entry, error)
	DeleteEntry(ctx context.Context, id int) error
}
