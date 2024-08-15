package idmap

type CacheStrategy interface {
	remove(m *IDMap)
}

type FIFO struct {
}

func (f *FIFO) remove(m *IDMap) {
	for i := 0; i < m.DeleteCount; i++ {
		delete(m.Cache, m.Keys[i])
	}
	m.Keys = m.Keys[m.DeleteCount:]
	m.Size -= m.DeleteCount
}
