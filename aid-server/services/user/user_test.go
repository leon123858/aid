package user

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Connect(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) IsExist(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockDB) Set(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockDB) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

var defaultDeviceFingerPrint = DeviceFingerPrint{
	IP:   "127.0.0.1",
	OS:   "Windows",
	Brow: "Chrome",
}

func TestCreateUser(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	db.On("IsExist", aid.String()).Return(false, nil)
	db.On("Set", aid.String(), mock.AnythingOfType("string")).Return(nil)
	u, err := CreateUser(aid, db)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, aid, u.GetAID())
	assert.Equal(t, &Space{}, u.GetSpace())
	assert.Equal(t, &Time{}, u.GetTime())
	db.AssertExpectations(t)
}

func TestUser_SetSpace(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	u := &User{
		ID: aid,
		DB: db,
		Data: Data{
			Space: Space{},
			Time:  Time{},
		},
	}
	newSpace := Space{
		DeviceFingerPrint: defaultDeviceFingerPrint,
	}
	db.On("Set", aid.String(), mock.AnythingOfType("string")).Return(nil)
	err := u.SetSpace(newSpace)
	assert.NoError(t, err)
	assert.Equal(t, &newSpace, u.GetSpace())
	db.AssertExpectations(t)
}

func TestUser_SetSpace_Error(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	u := &User{
		ID: aid,
		DB: db,
		Data: Data{
			Space: Space{},
			Time:  Time{},
		},
	}
	newSpace := Space{
		DeviceFingerPrint: defaultDeviceFingerPrint,
	}
	db.On("Set", aid.String(), mock.AnythingOfType("string")).Return(assert.AnError)
	err := u.SetSpace(newSpace)
	assert.Error(t, err)
	db.AssertExpectations(t)
}

func TestUser_SetTime(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	u := &User{
		ID: aid,
		DB: db,
		Data: Data{
			Space: Space{},
			Time:  Time{},
		},
	}
	newTime := Time{
		preLoginTime: time.Now(),
	}
	db.On("Set", aid.String(), mock.AnythingOfType("string")).Return(nil)
	err := u.SetTime(newTime)
	assert.NoError(t, err)
	assert.Equal(t, &newTime, u.GetTime())
	db.AssertExpectations(t)
}

func TestUser_SetAll(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	u := &User{
		ID: aid,
		DB: db,
		Data: Data{
			Space: Space{},
			Time:  Time{},
		},
	}
	newData := Data{
		Space: Space{
			DeviceFingerPrint: defaultDeviceFingerPrint,
		},
		Time: Time{
			preLoginTime: time.Now(),
		},
	}
	db.On("Set", aid.String(), mock.AnythingOfType("string")).Return(nil)
	err := u.SetAll(newData)
	assert.NoError(t, err)
	assert.Equal(t, &newData, &u.Data)
	assert.Equal(t, &newData.Space, u.GetSpace())
	assert.Equal(t, &newData.Time, u.GetTime())
	db.AssertExpectations(t)
}

func TestNewDB(t *testing.T) {
	db, err := NewDB("test")
	assert.NoError(t, err)
	err = db.Close()
	assert.NoError(t, err)
	err = os.RemoveAll("test.ldb")
	assert.Nil(t, err)
}
