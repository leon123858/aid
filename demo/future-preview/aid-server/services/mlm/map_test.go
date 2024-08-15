package mlm

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var list = []uuid.UUID{
	uuid.New(), uuid.New(), uuid.New(), uuid.New(),
}

func TestMultiLevelMap(t *testing.T) {
	t.Run("NewMultiLevelMap", func(t *testing.T) {
		m := NewMultiLevelMap()
		assert.NotNil(t, m, "NewMultiLevelMap() should not return nil")
		_, ok := m.(*MultiLevelMapImp)
		assert.True(t, ok, "NewMultiLevelMap() should return a *MultiLevelMapImp")
	})

	t.Run("Set Same Value", func(t *testing.T) {
		m := NewMultiLevelMap()
		testCases := []struct {
			key   KeyItem
			value uuid.UUID
		}{
			{KeyItem{IP: "192.168.1.1", Browser: "Chrome"}, list[0]},
			{KeyItem{IP: "192.168.1.2", Browser: "Firefox"}, list[1]},
			{KeyItem{IP: "192.168.1.1", Browser: "Safari"}, list[2]},
			{KeyItem{IP: "192.168.1.1", Browser: "Safari"}, list[3]},
		}

		for _, tc := range testCases {
			err := m.Set(tc.key, tc.value)
			assert.NoError(t, err, "Set() should not return an error")

			results, err := m.Get(tc.key)
			assert.NoError(t, err, "Get() should not return an error")
			assert.Contains(t, results, tc.value, "Get() should return the set value")
		}

		// Test duplicate value
		results, err := m.Get(testCases[3].key)
		assert.NoError(t, err)
		assert.ElementsMatch(t, results, []uuid.UUID{list[2], list[3]}, "Get() should return all values for duplicate value")
	})

	t.Run("Set and Get", func(t *testing.T) {
		m := NewMultiLevelMap()
		testCases := []struct {
			key   KeyItem
			value uuid.UUID
		}{
			{KeyItem{IP: "192.168.1.1", Browser: "Chrome"}, list[0]},
			{KeyItem{IP: "192.168.1.2", Browser: "Firefox"}, list[1]},
			{KeyItem{IP: "192.168.1.1", Browser: "Safari"}, list[2]},
		}

		for _, tc := range testCases {
			err := m.Set(tc.key, tc.value)
			assert.NoError(t, err, "Set() should not return an error")

			results, err := m.Get(tc.key)
			assert.NoError(t, err, "Get() should not return an error")
			assert.Contains(t, results, tc.value, "Get() should return the set value")
		}

		// Test empty key
		results, err := m.Get(KeyItem{IP: "192.168.1.1", Browser: "Edge"})
		assert.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("Set and Get with Clear", func(t *testing.T) {
		m := NewMultiLevelMap().(*MultiLevelMapImp)
		testCases := []struct {
			key   KeyItem
			value uuid.UUID
		}{
			{KeyItem{IP: "127.0.0.1", Browser: "Chrome"}, list[0]},
			{KeyItem{IP: "127.0.0.1", Browser: "Chrome"}, list[1]},
		}

		err := m.Set(testCases[0].key, testCases[0].value)
		assert.NoError(t, err, "Set() should not return an error")
		results, err := m.Get(testCases[0].key)
		assert.NoError(t, err, "Get() should not return an error")
		assert.Equal(t, 1, len(results), "Get() should return the set value")

		err = m.Set(testCases[1].key, testCases[1].value)
		assert.NoError(t, err, "Set() should not return an error")
		results, err = m.Get(testCases[1].key)
		assert.NoError(t, err, "Get() should not return an error")
		assert.Equal(t, 2, len(results), "Get() should return the set value")

		closeDuration = time.Second
		time.Sleep(2 * time.Second)
		results, err = m.Get(testCases[0].key)
		assert.Nil(t, err)
		assert.Equal(t, list[1], results[0])
	})

	t.Run("MaxSize", func(t *testing.T) {
		m := NewMultiLevelMap().(*MultiLevelMapImp)
		for i := 0; i <= maxSize+1; i++ {
			key := KeyItem{IP: fmt.Sprintf("192.168.1.%d", i), Browser: fmt.Sprintf("Browser%d", i)}
			assert.NoError(t, m.Set(key, uuid.New()))
		}
		assert.Equal(t, maxSize+1-clearSize+1, m.size, "Size should be reduced to clearSize after exceeding maxSize")
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		m := NewMultiLevelMap()
		var wg sync.WaitGroup
		numGoroutines := 100

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				defer wg.Done()
				key := KeyItem{IP: fmt.Sprintf("192.168.1.%d", i), Browser: fmt.Sprintf("Browser%d", i)}
				value := uuid.New()

				err := m.Set(key, value)
				assert.NoError(t, err, "Concurrent Set() should not return an error")

				results, err := m.Get(key)
				assert.NoError(t, err, "Concurrent Get() should not return an error")
				assert.Contains(t, results, value, "Concurrent Get() should return the set value")
			}(i)
		}

		wg.Wait()
	})

	t.Run("EdgeCases", func(t *testing.T) {
		m := NewMultiLevelMap()
		testCases := []struct {
			name  string
			key   KeyItem
			value uuid.UUID
		}{
			{"Empty IP", KeyItem{IP: "", Browser: "Chrome"}, uuid.New()},
			{"Empty Browser", KeyItem{IP: "192.168.1.1", Browser: ""}, uuid.New()},
			{"Both Empty", KeyItem{IP: "", Browser: ""}, uuid.New()},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := m.Set(tc.key, tc.value)
				assert.NoError(t, err, "Set() should not return an error for edge cases")

				results, err := m.Get(tc.key)
				assert.NoError(t, err, "Get() should not return an error for edge cases")
				assert.Contains(t, results, tc.value, "Get() should return the set value for edge cases")
			})
		}
	})

	t.Run("EdgeCasesError", func(t *testing.T) {
		m := NewMultiLevelMap()
		testCases := []struct {
			name  string
			key   KeyItem
			value uuid.UUID
		}{
			{"Nil Value", KeyItem{IP: "127.0.0.1", Browser: "Chrome"}, uuid.Nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := m.Set(tc.key, tc.value)
				assert.Error(t, err, "Set() should return an error for edge cases")

				results, err := m.Get(tc.key)
				assert.Error(t, err, "Get() should not return an error for edge cases")
				assert.NotContains(t, results, tc.value, "Get() should return the set value for edge cases")
			})
		}
	})
}

func BenchmarkMultiLevelMap(b *testing.B) {
	m := NewMultiLevelMap()
	rand.Seed(time.Now().UnixNano())

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := KeyItem{
				IP:      fmt.Sprintf("192.168.1.%d", rand.Intn(256)),
				Browser: fmt.Sprintf("Browser%d", rand.Intn(10)),
			}
			assert.NoError(b, m.Set(key, uuid.New()))
		}
	})

	b.Run("Get", func(b *testing.B) {
		keys := make([]KeyItem, 1000)
		for i := 0; i < 1000; i++ {
			key := KeyItem{
				IP:      fmt.Sprintf("192.168.1.%d", rand.Intn(256)),
				Browser: fmt.Sprintf("Browser%d", rand.Intn(10)),
			}
			keys[i] = key
			assert.NoError(b, m.Set(key, uuid.New()))
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := m.Get(keys[i%1000])
			if err == nil {
				//println(get)
			}
		}
	})
}
