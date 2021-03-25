package memorycache

import (
	"sync"
)

type memoryCache struct {
	data map[string]interface{}
	m    sync.RWMutex
}

func New() *memoryCache {
	return &memoryCache{
		data: make(map[string]interface{}),
	}
}

func (m *memoryCache) Get(key string) interface{} {
	m.m.RLock()
	defer m.m.RUnlock()

	if v, ok := m.data[key]; ok {
		return v
	}
	return nil
}

func (m *memoryCache) Set(key string, value interface{}) {
	m.m.Lock()
	defer m.m.Unlock()

	m.data[key] = value
}
