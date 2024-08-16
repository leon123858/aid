package repository

import (
	"github.com/syndtr/goleveldb/leveldb"
)

var LDB *LevelDBStore

func init() {
	// 這個函數將在程序啟動時自動調用
	// 用於設置數據庫的初始化工作
	ldb, err := NewLevelDBStore("./data/ldb")
	if err != nil {
		panic(err)
	}
	LDB = ldb
}

type LevelDBStore struct {
	db *leveldb.DB
}

// NewLevelDBStore 創建一個新的 LevelDBStore
func NewLevelDBStore(dbPath string) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{db: db}, nil
}

// Close 關閉數據庫連接
func (s *LevelDBStore) Close() error {
	return s.db.Close()
}

// Create 創建一個新的鍵值對
func (s *LevelDBStore) Create(key, value []byte) error {
	return s.db.Put(key, value, nil)
}

// Read 讀取一個鍵的值
func (s *LevelDBStore) Read(key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}

// Update 更新一個鍵的值
func (s *LevelDBStore) Update(key, value []byte) error {
	return s.db.Put(key, value, nil)
}

// Delete 刪除一個鍵值對
func (s *LevelDBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

// List 列出所有的鍵值對
func (s *LevelDBStore) List() (map[string]string, error) {
	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()

	result := make(map[string]string)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		result[string(key)] = string(value)
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return result, nil
}
