package ldb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
)

type LDBTestSuite struct {
	suite.Suite
	dir string
	ldb *ldb
}

func (s *LDBTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "test_ldb")
	s.NoError(err)
	s.dir = dir

	s.ldb = &ldb{}
	err = s.ldb.Connect(s.dir + "/test")
	s.NoError(err)
}

func (s *LDBTestSuite) TearDownTest() {
	s.NoError(s.ldb.Close())
	err := os.RemoveAll(s.dir)
	if err != nil {
		return
	}
}

func TestLDBTestSuite(t *testing.T) {
	suite.Run(t, new(LDBTestSuite))
}

func (s *LDBTestSuite) TestLDB_Set() {
	err := s.ldb.Set("key", "value")
	s.NoError(err)
}

func (s *LDBTestSuite) TestLDB_Get() {
	err := s.ldb.Set("key", "value")
	s.NoError(err)

	value, err := s.ldb.Get("key")
	s.NoError(err)
	s.Equal("value", value)
}

func (s *LDBTestSuite) TestLDB_Delete() {
	err := s.ldb.Set("key", "value")
	s.NoError(err)

	err = s.ldb.Delete("key")
	s.NoError(err)

	_, err = s.ldb.Get("key")
	s.Equal(leveldb.ErrNotFound, err)
}

func (s *LDBTestSuite) TestLDB_IsExist() {
	err := s.ldb.Set("key", "value")
	s.NoError(err)

	exist, err := s.ldb.IsExist("key")
	s.NoError(err)
	s.True(exist)
}

func (s *LDBTestSuite) TestLDB_IsExist_NotFound() {
	exist, err := s.ldb.IsExist("key")
	s.NoError(err)
	s.False(exist)
}

func TestNew(t *testing.T) {
	db := New()
	assert.NotNil(t, db)
	assert.IsType(t, &ldb{}, db)
}
