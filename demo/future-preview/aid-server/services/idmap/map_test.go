package idmap

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Get(uid string) (string, error) {
	args := m.Called(uid)
	return args.String(0), args.Error(1)
}

func (m *MockDB) Set(uid, aid string) error {
	args := m.Called(uid, aid)
	return args.Error(0)
}

func (m *MockDB) Delete(uid string) error {
	args := m.Called(uid)
	return args.Error(0)
}

func (m *MockDB) IsExist(uid string) (bool, error) {
	args := m.Called(uid)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Connect(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func TestNewIDMap(t *testing.T) {
	db := &MockDB{}
	size := 10

	idMap := NewIDMap(size, db)

	assert.NotNil(t, idMap)
	assert.Equal(t, size, idMap.Limit)
	assert.Equal(t, size/5, idMap.DeleteCount)
	assert.Equal(t, 0, idMap.Size)
	assert.Equal(t, db, idMap.DB)
}

func TestIDMap_Get(t *testing.T) {
	db := &MockDB{}
	idMap := NewIDMap(10, db)

	uid := "user123"
	aid := "aid456"

	// Test cache hit
	idMap.Cache[uid] = aid
	result, err := idMap.Get(uid)
	assert.NoError(t, err)
	assert.Equal(t, aid, result)

	// Test cache miss and successful database retrieval
	delete(idMap.Cache, uid)
	db.On("Get", uid).Return(aid, nil)
	result, err = idMap.Get(uid)
	assert.NoError(t, err)
	assert.Equal(t, aid, result)
	assert.Equal(t, 1, idMap.Size)
	assert.Contains(t, idMap.Cache, uid)
	assert.Contains(t, idMap.Keys, uid)
}

func TestIDMap_GetDBError(t *testing.T) {
	db := &MockDB{}
	idMap := NewIDMap(10, db)

	uid := "user123"
	errDB := errors.New("database error")
	db.On("Get", uid).Return("", errDB)
	result, err := idMap.Get(uid)
	assert.Error(t, err)
	assert.Equal(t, errDB, err)
	assert.Empty(t, result)
}

func TestIDMap_Set(t *testing.T) {
	db := &MockDB{}
	idMap := NewIDMap(10, db)

	uid := "user123"
	aid := "aid456"

	db.On("Set", uid, aid).Return(nil)
	err := idMap.Set(uid, aid)
	assert.NoError(t, err)
}

func TestIDMap_remove(t *testing.T) {
	db := &MockDB{}
	idMap := NewIDMap(10, db)

	idMap.Keys = []string{"user1", "user2", "user3", "user4", "user5"}
	idMap.Size = len(idMap.Keys)

	idMap.remove(&FIFO{})

	assert.Equal(t, 2, idMap.DeleteCount)
	assert.Equal(t, 3, idMap.Size)
}
