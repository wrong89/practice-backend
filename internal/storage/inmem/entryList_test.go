package inmem

import (
	"context"
	"practice-backend/internal/models/entry"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetEntry(t *testing.T) {
	l := NewEntryList()

	initialEntry := entry.NewEntry(
		"some",
		time.Now(),
		0,
		"card",
	)

	l.list.AddData(*initialEntry)

	testCases := []struct {
		title   string
		id      int
		wantErr string
	}{
		{
			title:   "happy: get existing data",
			id:      0,
			wantErr: "",
		},
		{
			title:   "sad: get data by invalid id",
			id:      -1,
			wantErr: "invalid id",
		},
		{
			title:   "sad: get data by not existing id",
			id:      10,
			wantErr: "entry not found",
		},
	}

	for _, tc := range testCases {
		entry, err := l.GetEntryByID(context.Background(), tc.id)
		if tc.wantErr != "" {
			assert.Contains(t, err.Error(), tc.wantErr)
			assert.Empty(t, entry)
		} else {
			assert.NotEmpty(t, entry)
		}
	}
}

func TestDeleteEntry(t *testing.T) {
	l := NewEntryList()

	initialEntry := entry.NewEntry(
		"some",
		time.Now(),
		0,
		"card",
	)
	l.list.AddData(*initialEntry)

	testCases := []struct {
		title   string
		id      int
		wantErr string
	}{
		{
			title:   "happy: delete existing data",
			id:      0,
			wantErr: "",
		},
		{
			title:   "sad: delete data by invalid id",
			id:      -1,
			wantErr: "invalid id",
		},
		{
			title:   "sad: delete not existing data",
			id:      10,
			wantErr: "entry not found",
		},
	}

	for _, tc := range testCases {
		err := l.DeleteEntry(context.Background(), tc.id)
		if tc.wantErr != "" {
			assert.Contains(t, err.Error(), tc.wantErr)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestUpdateEntryStatus(t *testing.T) {
	l := NewEntryList()

	initialEntry := entry.NewEntry(
		"some",
		time.Now(),
		0,
		"card",
	)
	l.list.AddData(*initialEntry)

	testCases := []struct {
		title   string
		id      int
		status  string
		wantErr string
	}{
		{
			title:   "happy: mark existing data",
			id:      0,
			status:  "processed",
			wantErr: "",
		},
		{
			title:   "sad: delete data by invalid id",
			id:      -1,
			status:  "processed",
			wantErr: "invalid id",
		},
		{
			title:   "sad: delete not existing data",
			id:      10,
			status:  "processed",
			wantErr: "entry not found",
		},
	}

	for _, tc := range testCases {
		markedEntry, err := l.UpdateEntryStatus(t.Context(), tc.id, tc.status)
		if tc.wantErr != "" {
			assert.Contains(t, err.Error(), tc.wantErr, tc.title)
			continue
		}
		assert.Nil(t, err)

		if assert.Equal(t, "processed", markedEntry.Status) {
			entry, _ := l.GetEntryByID(context.TODO(), tc.id)
			assert.Equal(t, "processed", entry.Status, tc.title)
		}
	}
}

func TestGetEntries(t *testing.T) {
	l := NewEntryList()

	dataCount := 5
	initialTime := time.Now()

	initialEntry := entry.NewEntry(
		"some",
		initialTime,
		0,
		"card",
	)

	for range dataCount {
		l.list.AddData(*initialEntry)
	}

	entries, err := l.GetEntries(context.TODO())
	assert.Nil(t, err)

	for _, entry := range entries {
		assert.Equal(t, *initialEntry, entry)
	}
}
