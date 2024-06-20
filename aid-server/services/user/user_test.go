package user

import (
	"aid-server/pkg/ldb"
	"aid-server/pkg/timestamp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
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
	assert.Equal(t, Space{
		DeviceFingerPrint: DeviceFingerPrint{
			IP:   "",
			Brow: "",
		},
	}, *u.GetSpace())
	assert.Equal(t, Time{
		CurEventTime: 0,
	}, *u.GetTime())
	db.AssertExpectations(t)
}

func TestUser_SetAll(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	u := &User{
		ID: aid,
		DB: db,
		Data: Data{
			Records: make([]Record, 0),
			Info: Info{
				PublicKey: "",
				AID:       aid.String(),
			},
		},
	}
	newData := Record{
		Space: Space{
			DeviceFingerPrint: defaultDeviceFingerPrint,
		},
		Time: Time{
			CurEventTime: timestamp.GetTime(),
		},
	}
	db.On("Set", aid.String(), mock.AnythingOfType("string")).Return(nil)
	err := u.SetRecord(newData)
	assert.NoError(t, err)
	assert.Equal(t, &newData, &u.Data.Records[0])
	assert.Equal(t, &newData.Space, u.GetSpace())
	assert.Equal(t, &newData.Time, u.GetTime())
	db.AssertExpectations(t)
}

func TestNewDB(t *testing.T) {
	db, err := ldb.NewDB("test")
	assert.NoError(t, err)
	err = ldb.FreeDB(db)
	assert.NoError(t, err)
	err = os.RemoveAll("test.ldb")
	assert.Nil(t, err)
}

func TestUser_IsExist(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	db.On("IsExist", aid.String()).Return(true, nil)
	db.On("Get", aid.String()).Return(`{"Space":{},"Time":{}}`, nil)
	u, err := CreateUser(aid, db)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.True(t, u.IsExist())
	db.AssertExpectations(t)
}

func TestUser_IsNotExist(t *testing.T) {
	db := new(MockDB)
	aid := uuid.New()
	db.On("IsExist", aid.String()).Return(false, nil)
	db.On("Set", aid.String(), mock.AnythingOfType("string")).Return(nil)
	u, err := CreateUser(aid, db)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.False(t, u.IsExist())
	db.AssertExpectations(t)
}
