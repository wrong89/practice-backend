package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePhone(t *testing.T) {
	testCases := []struct {
		title       string
		phone       string
		expectedErr string
	}{
		{
			title:       "happy: valid phone",
			phone:       "13052858783",
			expectedErr: "",
		},
		{
			title:       "sad: empty phone",
			phone:       "",
			expectedErr: "phone number is empty",
		},
		{
			title:       "sad: invalid phone",
			phone:       "83746",
			expectedErr: "invalid phone number",
		},
		{
			title:       "sad: invalid phone",
			phone:       "289374628479356238947563789",
			expectedErr: "invalid phone number",
		},
	}

	for _, tc := range testCases {
		err := ValidatePhone(tc.phone)
		if tc.expectedErr != "" {
			assert.Contains(t, err.Error(), tc.expectedErr, tc.title)
		} else {
			assert.Nil(t, err, tc.title)
		}
	}
}
