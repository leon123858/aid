package mlm

import (
	"aid-server/configs"
	"errors"
	"github.com/emirpasic/gods/maps/hashmap"
	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/google/uuid"
	"sync"
	"time"
)

/*
	multi-level map, this is the simple implementation of multidimensional red black tree
	we hope to find the closest value group by all properties in struct
	but the concept is difficult to implement, so we just use the simple map to store the data
	this implementation is not the best, but it is simple and easy to understand
	it can find same value by all properties in struct quickly
*/

var closeDuration = configs.Configs.Time.LoginCache

const maxSize = 100
const clearSize = 30

type KeyItem struct {
	IP      string
	Browser string
}

type ValueItem struct {
	AID  uuid.UUID
	Time time.Time
}

type MultiLevelMap interface {
	// Set value to map
	Set(key KeyItem, value uuid.UUID) error
	// Get values from map
	Get(key KeyItem) ([]uuid.UUID, error)
}

type MultiLevelMapImp struct {
	// maps to store the data
	hm hashmap.Map
	// mtx to lock the map
	mtx sync.Mutex
	// size of map
	size int
	// key list
	keyListQueue llq.Queue
}

func NewMultiLevelMap() MultiLevelMap {
	m := &MultiLevelMapImp{
		hm:           *hashmap.New(),
		size:         0,
		mtx:          sync.Mutex{},
		keyListQueue: *llq.New(),
	}
	return m
}

func (m *MultiLevelMapImp) Set(key KeyItem, value uuid.UUID) error {
	if value == uuid.Nil {
		return errors.New("value is nil")
	}
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.size > maxSize {
		clearKeys(m)
	}
	// check if key exists
	data, ok := m.hm.Get(key)
	if !ok {
		m.size++
		m.keyListQueue.Enqueue(key)
		m.hm.Put(key, []ValueItem{{
			AID:  value,
			Time: time.Now(),
		}})
		return nil
	}
	// if key exists, concat the value
	arr := data.([]ValueItem)
	// check if value exists
	for _, v := range arr {
		if v.AID == value {
			return nil
		}
	}
	m.hm.Put(key, append(arr, ValueItem{
		AID:  value,
		Time: time.Now(),
	}))
	return nil
}

func (m *MultiLevelMapImp) Get(key KeyItem) ([]uuid.UUID, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	v, ok := m.hm.Get(key)
	if !ok {
		return nil, errors.New("key not found")
	}
	// sort by time
	arr := v.([]ValueItem)
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i].Time.Before(arr[j].Time) {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	// get first 2 item if more than 2
	if len(arr) > 2 {
		arr = arr[:2]
		// check if each item time is in 5 minutes
		if arr[0].Time.Add(closeDuration).Before(arr[1].Time) {
			arr = arr[:1]
		}
	}
	// if value is too old, remove it
	if arr[0].Time.Add(configs.Configs.Jwt.Duration).Before(time.Now()) {
		m.hm.Remove(key)
		m.size--
		return nil, errors.New("value is too old")
	}
	// extract AID
	res := make([]uuid.UUID, 0)
	for _, v := range arr {
		res = append(res, v.AID)
	}

	return res, nil
}
