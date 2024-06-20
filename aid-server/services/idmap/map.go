package idmap

import (
	"aid-server/pkg/ldb"
	"sync"
)

type IDMap struct {
	Cache       map[string]string
	Keys        []string
	Size        int
	Limit       int
	DB          ldb.DB
	DeleteCount int
	mtx         sync.Mutex
}

func NewIDMap(size int, db ldb.DB) *IDMap {
	deleteCount := size / 5
	if deleteCount == 0 {
		panic("size should be bigger than 5")
	}
	return &IDMap{
		Cache:       make(map[string]string, size),
		Keys:        make([]string, 0),
		Size:        0,
		Limit:       size,
		DB:          db,
		DeleteCount: deleteCount,
	}
}

func (i *IDMap) Get(uid string) (string, error) {
	if v, ok := i.Cache[uid]; ok {
		return v, nil
	}
	v, err := i.DB.Get(uid)
	if err != nil {
		return "", err
	}
	i.mtx.Lock()
	defer i.mtx.Unlock()
	i.Cache[uid] = v
	i.Keys = append(i.Keys, uid)
	i.Size++
	if i.Size > i.Limit {
		i.remove(&FIFO{})
	}
	return v, nil
}

func (i *IDMap) Set(uid, aid string) error {
	return i.DB.Set(uid, aid)
}

func (i *IDMap) remove(s CacheStrategy) {
	s.remove(i)
}
