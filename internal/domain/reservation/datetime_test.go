package reservation

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected DateTime
		err      error
	}{
		{
			input:    `"29-08-2024 13:00"`,
			expected: DateTime{time.Date(2024, 8, 29, 13, 0, 0, 0, time.UTC)},
			err:      nil,
		},
		{
			input: `"invalid time format"`,
			err:   &time.ParseError{},
		},
		{
			input:    `"null"`,
			expected: DateTime{time.Time{}},
			err:      nil,
		},
	}

	for _, test := range tests {
		var dt DateTime
		err := json.Unmarshal([]byte(test.input), &dt)
		if test.err != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, dt)
		}
	}
}

func TestDateTimeMarshalJSON(t *testing.T) {
	tests := []struct {
		input    DateTime
		expected string
	}{
		{
			input:    DateTime{time.Date(2024, 8, 29, 13, 0, 0, 0, time.UTC)},
			expected: `"29-08-2024 13:00"`,
		},
		{
			input:    DateTime{time.Time{}},
			expected: `"01-01-0001 00:00"`, // Go's zero time
		},
	}

	for _, test := range tests {
		output, err := json.Marshal(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, string(output))
	}
}
