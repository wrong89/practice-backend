package jwt_test

import (
	"practice-backend/internal/lib/jwt"
	"practice-backend/internal/models/user"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewToken(t *testing.T) {
	testCases := []struct {
		name          string
		user          user.User
		duration      time.Duration
		expectedError error
	}{
		{
			name: "Success",
			user: user.User{
				ID:    0,
				Email: "test@test.com",
			},
			duration: time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := jwt.NewToken(tc.user, tc.duration)
			require.ErrorIs(t, err, tc.expectedError)

			require.NotZero(t, token)
		})
	}
}
