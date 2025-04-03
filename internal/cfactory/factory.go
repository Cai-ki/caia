package cfactory

import (
	"github.com/Cai-ki/caia/internal/ctypes"
)

type Option func(*Manager)

type Manager struct {
	name     string
	factorys map[string]ctypes.HandleFactoryFunc
}

var _ ctypes.Factory = (*Manager)(nil)

func NewManager(name string, opts ...Option) *Manager {
	m := &Manager{
		name: name,
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Manager) RegisterFactory(name string, factory ctypes.HandleFactoryFunc) {
	m.factorys[name] = factory
}

func (m *Manager) GetFactory(name string) (ctypes.HandleFactoryFunc, bool) {
	factory, ok := m.factorys[name]
	return factory, ok
}
