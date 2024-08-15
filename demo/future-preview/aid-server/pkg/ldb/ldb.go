package ldb

import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type DB interface {
	Connect(path string) error
	Close() error
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	IsExist(key string) (bool, error)
}

type ldb struct {
	db *leveldb.DB
}

func NewDB(path string) (DB, error) {
	db := New()
	err := db.Connect(path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func FreeDB(db DB) error {
	return db.Close()
}

func (l *ldb) Connect(path string) error {
	db, err := leveldb.OpenFile(path+".ldb", nil)
	if err != nil {
		return err
	}
	l.db = db
	return nil
}

func (l *ldb) Close() error {
	return l.db.Close()
}

func (l *ldb) Set(key string, value string) error {
	return l.db.Put([]byte(key), []byte(value), nil)
}

func (l *ldb) Get(key string) (string, error) {
	data, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (l *ldb) Delete(key string) error {
	return l.db.Delete([]byte(key), nil)
}

func (l *ldb) IsExist(key string) (bool, error) {
	data, err := l.db.Get([]byte(key), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	if data == nil {
		return false, nil
	}
	return true, nil
}

func New() DB {
	return &ldb{}
}
