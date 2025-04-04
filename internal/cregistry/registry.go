package cregistry

import (
	"sync"

	"github.com/Cai-ki/caia/internal/ctypes"
)

type Option func(*Manager)

type Manager struct {
	name    string
	methods map[string]ctypes.MethodFunc
	handles map[string]ctypes.HandleFunc
	pools   map[string]*sync.Pool
}

var _ ctypes.Registry = (*Manager)(nil)

func NewManager(name string, opts ...Option) *Manager {
	m := &Manager{
		name:    name,
		methods: make(map[string]ctypes.MethodFunc),
		handles: make(map[string]ctypes.HandleFunc),
		pools:   make(map[string]*sync.Pool),
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Manager) GetName() string {
	return m.name
}

func (m *Manager) GetMethod(name string) (ctypes.MethodFunc, error) {
	method, ok := m.methods[name]
	if !ok {
		return nil, ctypes.ErrNotFound
	}
	return method, nil
}

func (m *Manager) GetHandle(name string) (ctypes.HandleFunc, error) {
	handle, ok := m.handles[name]
	if !ok {
		return nil, ctypes.ErrNotFound
	}
	return handle, nil
}

func (m *Manager) GetPool(name string) (*sync.Pool, error) {
	pool, ok := m.pools[name]
	if !ok {
		return nil, ctypes.ErrNotFound
	}
	return pool, nil
}

func (m *Manager) RegisterMethod(name string, method ctypes.MethodFunc) error {
	if method == nil {
		return ctypes.ErrInvalidArgument
	}
	m.methods[name] = method
	return nil
}

func (m *Manager) RegisterHandle(name string, handle ctypes.HandleFunc) error {
	if handle == nil {
		return ctypes.ErrInvalidArgument
	}
	m.handles[name] = handle
	return nil
}

func (m *Manager) RegisterPool(name string, pool *sync.Pool) error {
	if pool == nil {
		return ctypes.ErrInvalidArgument
	}
	m.pools[name] = pool
	return nil
}
