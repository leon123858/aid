package rba

import (
	"aid-server/pkg/timestamp"
	"aid-server/services/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockUser struct {
	user.Time
	user.Space
}

func (m *MockUser) IsExist() bool {
	return true
}

func (m *MockUser) GetAID() uuid.UUID {
	return uuid.New()
}

func (m *MockUser) GetInfo() user.Info {
	return user.Info{}
}

func (m *MockUser) GetSpace() *user.Space {
	return &m.Space
}

func (m *MockUser) GetTime() *user.Time {
	return &m.Time
}

func (m *MockUser) SetRecord(user.Record) error {
	return nil
}

func (m *MockUser) SetInfo(user.Info) error {
	return nil
}

func TestSimpleStrategy_Verify(t *testing.T) {
	// Prepare test data
	userItem := &MockUser{
		Time: user.Time{
			CurEventTime: timestamp.GetTime(),
		},
		Space: user.Space{DeviceFingerPrint: user.DeviceFingerPrint{
			IP:   "127.0.0.1",
			Brow: "Chrome",
		}}}

	tests := []struct {
		name     string
		input    *user.Record
		expected bool
	}{
		{
			name: "Valid input",
			input: &user.Record{
				Time: user.Time{
					CurEventTime: timestamp.GetTime(),
				},
				Space: user.Space{DeviceFingerPrint: user.DeviceFingerPrint{
					IP:   "127.0.0.1",
					Brow: "Chrome",
				}},
			},
			expected: true,
		},
		{
			name: "Invalid time",
			input: &user.Record{
				Time: user.Time{
					CurEventTime: timestamp.GetTime() + 60000,
				},
				Space: user.Space{DeviceFingerPrint: user.DeviceFingerPrint{
					IP:   "127.0.0.1",
					Brow: "Chrome",
				}},
			},
			expected: false,
		},
		{
			name: "Invalid space",
			input: &user.Record{
				Time: user.Time{
					CurEventTime: timestamp.GetTime(),
				},
				Space: user.Space{DeviceFingerPrint: user.DeviceFingerPrint{
					IP:   "127.0.0.1",
					Brow: "Safari",
				}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			algo := SimpleAlgo
			result := algo.Verify(userItem, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
