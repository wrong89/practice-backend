package inmem

import (
	"context"
	"errors"
	"practice-backend/internal/models/entry"
	"practice-backend/internal/storage/inmem/ilist"
	"sync"
	"time"
)

var (
	ErrEntryNotFound = errors.New("entry not found")
)

// Concurrent-Use
type EntryList struct {
	list ilist.List[entry.Entry]
	mtx  *sync.Mutex
}

func NewEntryList() EntryList {
	return EntryList{
		list: ilist.NewList[entry.Entry](),
		mtx:  new(sync.Mutex),
	}
}

func (el *EntryList) GetEntryByID(ctx context.Context, id int) (entry.Entry, error) {
	el.mtx.Lock()
	defer el.mtx.Unlock()

	e, err := el.list.GetData(id)
	if err != nil {
		if errors.Is(err, ilist.ErrDataNotFound) {
			return entry.Entry{}, ErrEntryNotFound
		}
		return entry.Entry{}, err
	}
	return e, nil
}

func (el *EntryList) MarkAsProcessed(ctx context.Context, id int) (entry.Entry, error) {
	e, err := el.GetEntryByID(ctx, id)
	if err != nil {
		return entry.Entry{}, err
	}

	el.mtx.Lock()
	defer el.mtx.Unlock()

	e.MarkAsProcessed()

	return el.list.UpdateData(e.ID, e)
}

func (el *EntryList) CreateEntry(
	ctx context.Context,
	course string,
	date time.Time,
	userID int,
	paymentMethod string,
) (entry.Entry, error) {
	newEntry := entry.NewEntry(course, date, userID, paymentMethod)

	el.mtx.Lock()
	defer el.mtx.Unlock()

	newEntry.ID = el.list.GetLastID() + 1
	e, err := el.list.AddData(*newEntry)
	if err != nil {
		return entry.Entry{}, err
	}

	return e, nil
}

func (el *EntryList) DeleteEntry(ctx context.Context, id int) error {
	el.mtx.Lock()
	defer el.mtx.Unlock()

	e, err := el.list.GetData(id)
	if err != nil {
		return err
	}

	if err := el.list.DeleteData(e.ID); err != nil {
		if errors.Is(err, ilist.ErrDataNotFound) {
			return ErrEntryNotFound
		}
		return err
	}

	return nil
}
