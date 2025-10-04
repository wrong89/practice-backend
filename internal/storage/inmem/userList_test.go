package inmem

import (
	"context"
	"practice-backend/internal/models/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	l := NewUserList()

	initialUser := user.NewUser(
		"login123",
		"pass123",
		"Ivan",
		"Ivanov",
		"Ivanovich",
		"+79991234567",
		"ivan@example.com",
	)
	l.list.AddData(*initialUser)

	testCases := []struct {
		title   string
		id      int
		wantErr string
	}{
		{
			title:   "happy: get existing user",
			id:      0,
			wantErr: "",
		},
		{
			title:   "sad: get user by invalid id",
			id:      -1,
			wantErr: "invalid id",
		},
		{
			title:   "sad: get user by not existing id",
			id:      10,
			wantErr: "user not found",
		},
	}

	for _, tc := range testCases {
		u, err := l.GetUserByID(context.Background(), tc.id)
		if tc.wantErr != "" {
			assert.Contains(t, err.Error(), tc.wantErr, tc.title)
			assert.Empty(t, u, tc.title)
		} else {
			assert.NotEmpty(t, u, tc.title)
		}
	}
}

func TestGetUserByLogin(t *testing.T) {
	l := NewUserList()

	initialLogin := "login123"

	l.CreateUser(
		t.Context(),
		initialLogin,
		"pass123",
		"Ivan",
		"Ivanov",
		"Ivanovich",
		"+79991234567",
		"ivan@example.com",
	)

	testCases := []struct {
		title   string
		login   string
		wantErr string
	}{
		{
			title:   "happy: get existing user",
			login:   initialLogin,
			wantErr: "",
		},
		{
			title:   "sad: get user by empty login",
			login:   "",
			wantErr: "user not found",
		},
	}

	for _, tc := range testCases {
		u, err := l.GetUserByLogin(context.Background(), tc.login)
		if tc.wantErr != "" {
			assert.Contains(t, err.Error(), tc.wantErr, tc.title)
			assert.Empty(t, u, tc.title)
		} else {
			assert.NotEmpty(t, u, tc.title)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	l := NewUserList()

	initialUser := user.NewUser(
		"login123",
		"pass123",
		"Ivan",
		"Ivanov",
		"Ivanovich",
		"+79991234567",
		"ivan@example.com",
	)
	l.list.AddData(*initialUser)

	testCases := []struct {
		title   string
		id      int
		wantErr string
	}{
		{
			title:   "happy: delete existing user",
			id:      0,
			wantErr: "",
		},
		{
			title:   "sad: delete user by invalid id",
			id:      -1,
			wantErr: "invalid id",
		},
		{
			title:   "sad: delete not existing user",
			id:      10,
			wantErr: "user not found",
		},
	}

	for _, tc := range testCases {
		err := l.DeleteUser(context.Background(), tc.id)
		if tc.wantErr != "" {
			assert.Contains(t, err.Error(), tc.wantErr, tc.title)
		} else {
			assert.Nil(t, err, tc.title)
		}
	}
}
