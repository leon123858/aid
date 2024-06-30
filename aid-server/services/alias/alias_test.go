package alias

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	// 使用內存數據庫進行測試
	db, err := NewDB(":memory:")
	require.NoError(t, err)
	defer db.Close()

	t.Run("AddUser", func(t *testing.T) {
		err := db.AddUser("user1", "John Doe", "1234")
		assert.NoError(t, err)
	})

	t.Run("ValidateUser", func(t *testing.T) {
		uids, err := db.ValidateUser("John Doe", "1234")
		assert.NoError(t, err)
		assert.Equal(t, []string{"user1"}, uids)

		uids, err = db.ValidateUser("John Doe", "wrong_pin")
		assert.NoError(t, err)
		assert.Empty(t, uids)
	})

	t.Run("AddLoginRecord", func(t *testing.T) {
		err := db.AddLoginRecord("user1", "127.0.0.1", "Chrome")
		assert.NoError(t, err)
	})

	t.Run("GetUserLoginHistory", func(t *testing.T) {
		record, err := db.GetUserLoginHistory("user1")
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, "127.0.0.1", record.IP)
		assert.Equal(t, "Chrome", record.Browser)
		assert.WithinDuration(t, time.Now(), record.LoginTime, time.Second)

		record, err = db.GetUserLoginHistory("non_existent_user")
		assert.NoError(t, err)
		assert.Nil(t, record)
	})
}

func TestNewDB(t *testing.T) {
	db, err := NewDB(":memory:")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()

	// 測試數據庫初始化
	_, err = db.db.Exec("INSERT INTO Users (Uid, Name, Pin) VALUES (?, ?, ?)", "test", "Test User", "1234")
	assert.NoError(t, err)
}
