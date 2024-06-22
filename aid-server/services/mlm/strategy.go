package mlm

func clearKeys(m *MultiLevelMapImp) {
	for i := 0; i < clearSize; i++ {
		key, ok := m.keyListQueue.Dequeue()
		if !ok {
			break
		}
		m.hm.Remove(key)
	}
	m.size = m.keyListQueue.Size()
}
