package cstorage

import (
	"sync"

	"github.com/Cai-ki/caia/internal/ctypes"
)

type Option func(*Manager)

type Manager struct {
	groups sync.Map // key: string (group name), value: *Resource
}

type Resource struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

var _ ctypes.Storage = (*Manager)(nil)

func NewManager(opts ...Option) *Manager {
	m := &Manager{}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (ms *Manager) Get(groupName, key string) (interface{}, error) {
	res, _ := ms.groups.LoadOrStore(groupName, &Resource{
		data: make(map[string]interface{}),
	})

	rr := res.(*Resource)
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	val, ok := rr.data[key]
	if !ok {
		return nil, ctypes.ErrKeyNotFound
	}
	return val, nil
}

func (ms *Manager) Put(groupName, key string, value interface{}) error {
	if value == nil {
		return ctypes.ErrInvalidArgument
	}

	res, _ := ms.groups.LoadOrStore(groupName, &Resource{
		data: make(map[string]interface{}),
	})

	rr := res.(*Resource)
	rr.mu.Lock()
	defer rr.mu.Unlock()

	rr.data[key] = value
	return nil
}

func (ms *Manager) Delete(groupName, key string) error {
	res, ok := ms.groups.Load(groupName)
	if !ok {
		return ctypes.ErrKeyNotFound
	}

	rr := res.(*Resource)
	rr.mu.Lock()
	defer rr.mu.Unlock()

	delete(rr.data, key)
	return nil
}

func (ms *Manager) Exists(groupName, key string) bool {
	res, ok := ms.groups.Load(groupName)
	if !ok {
		return false
	}

	rr := res.(*Resource)
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	_, exists := rr.data[key]
	return exists
}
