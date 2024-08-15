package timestamp

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	now := time.Now()
	ts := GetTime()

	assert.InDelta(t, now.UnixMilli(), int64(ts), 10, "GetTime should return current timestamp")
}

func TestToTimestamp(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Timestamp
	}{
		{
			name:     "Valid timestamp",
			input:    "1718814158970",
			expected: Timestamp(1718814158970),
		},
		{
			name:     "Empty timestamp",
			input:    "",
			expected: Timestamp(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToTimestamp(tc.input)
			assert.Equal(t, tc.expected, result, "ToTimestamp should convert string to Timestamp")
		})
	}
}

func TestTimestampString(t *testing.T) {
	ts := Timestamp(1687254600000)
	expected := "1687254600000"

	assert.Equal(t, expected, ts.String(), "Timestamp.String() should return formatted string")
}

func TestGetCurTimestamp(t *testing.T) {
	now := time.Now()
	ts := GetTime()

	assert.InDelta(t, now.UnixMilli(), int64(ts), 10, "GetCurTimestamp should return current timestamp")
}
