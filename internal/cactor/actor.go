package cactor

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/Cai-ki/caia/internal/clog"
	"github.com/Cai-ki/caia/internal/ctypes"
)

func WithValue(key, val interface{}) ctypes.OptionFunc {
	return func(i interface{}) {
		m := i.(*Manager)
		m.ctx = context.WithValue(m.ctx, key, val)
	}
}

func Registry(registry ctypes.Registry) ctypes.OptionFunc {
	return func(i interface{}) {
		m := i.(*Manager)
		m.registry = registry
	}
}

type Manager struct {
	name     string
	mailbox  ctypes.Mailbox
	ctx      context.Context
	registry ctypes.Registry
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	parent   ctypes.Actor
	children map[string]ctypes.Actor
	mu       sync.Mutex
	closed   atomic.Bool
	handle   ctypes.HandleFunc
}

var _ ctypes.Actor = (*Manager)(nil)

func NewManager(name string, buffer int, parent ctypes.Actor, handle ctypes.HandleFunc, opts ...ctypes.OptionFunc) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	if parent != nil {
		ctx, cancel = context.WithCancel(parent.GetContext())
	}
	m := &Manager{
		name:     name,
		mailbox:  *ctypes.NewMailbox(buffer),
		ctx:      ctx,
		cancel:   cancel,
		parent:   parent,
		children: map[string]ctypes.Actor{},
		handle:   handle,
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Manager) CreateChild(name string, buffer int, handle ctypes.HandleFunc, opts ...ctypes.OptionFunc) (ctypes.Actor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.children[name]
	if ok {
		return nil, ctypes.ErrKeyRepeat
	}
	child := NewManager(name, buffer, m, handle, opts...)
	m.children[name] = child

	return child, nil
}

func (m *Manager) DeleteChild(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	child, ok := m.children[name]
	if !ok {
		return ctypes.ErrChildNotFound
	}

	child.Stop()
	delete(m.children, name)
	return nil
}

func (m *Manager) GetName() string {
	return m.name
}

func (m *Manager) GetMailbox() ctypes.Mailbox {
	return m.mailbox
}

func (m *Manager) GetContext() context.Context {
	return m.ctx
}

func (m *Manager) GetParent() ctypes.Actor {
	return m.parent
}

func (m *Manager) GetChildren() map[string]ctypes.Actor {
	return m.children
}

func (m *Manager) GetRegistry() (ctypes.Registry, error) {
	if m.registry == nil {
		return nil, ctypes.ErrNotFound
	}
	return m.registry, nil
}

func (m *Manager) Start() {
	m.wg.Add(1)
	m.serve() //TODO 这里可能有协程启动时不符合预期父协程先，子协程后的顺序。建议之后压测一下。

	for _, child := range m.children {
		child.Start()
	}
}

func (m *Manager) Stop() {
	m.stop()
}

func (m *Manager) StopWithErase() {
	m.stop()
	if m.parent != nil {
		err := m.parent.DeleteChild(m.name)
		if err != nil {
			clog.Error("actor: fail to delete actor:", m.name)
		}
	}
}

func (m *Manager) SendMessage(msg ctypes.Message) error {
	select {
	case m.mailbox.Chan() <- msg:
		return nil
	case <-m.ctx.Done():
		return ctypes.ErrActorClosed
	}
}

func (m *Manager) SendMessageAsync(msg ctypes.Message) error {
	select {
	case m.mailbox.Chan() <- msg:
		return nil
	case <-m.ctx.Done():
		return ctypes.ErrActorClosed
	default:
		return ctypes.ErrChannelFull
	}
}

func (m *Manager) SendMessageToChildren(msg ctypes.Message) {
	for _, child := range m.children {
		if child != nil {
			child.SendMessage(msg)
		}
	}
}

func (m *Manager) SendMessageAsyncToChildren(msg ctypes.Message) {
	for _, child := range m.children {
		if child != nil {
			child.SendMessageAsync(msg)
		}
	}
}

func (m *Manager) serve() {
	go func() {
		defer m.stop()
		defer m.wg.Done()

		for {
			select {
			case msg := <-m.mailbox.Chan():
				m.handleMessage(msg)
			case <-m.ctx.Done():
				if m.closed.Load() { // 正常退出流程，保证子协程正常退出后父协程退出。
					for {
						select {
						case msg := <-m.mailbox.Chan():
							m.handleMessage(msg)
						default:
							return
						}
					}
				} else { // 异常退出流程，协程树同时退出，不确保资源安全，不执行handleMessage
					return
				}
			}
		}
	}()
}

func (m *Manager) cleanup() {
	m.mailbox.Close()
}

func (m *Manager) stop() {
	if m.closed.Load() {
		return
	}
	m.closed.Store(true)

	for _, child := range m.children {
		child.Stop()
	}

	m.cancel()
	m.wg.Wait()
	m.cleanup()

	clog.Debugf("actor: %s actor stop", m.name)
}

func (m *Manager) handleMessage(msg ctypes.Message) {
	defer func() { // panic终止于此处，默认不信任所有执行的handle，Manager只负责接受message，并根据构造时传入的handle进行相应处理，进行一个逻辑转发的工作。
		if err := recover(); err != nil {
			clog.Errorf("actor: %s panic when handle message: %v with error: %v", m.name, msg, err)
		}
	}()

	if m.handle == nil {
		clog.Errorf("actor: %s handle not found", m.name)
		return
	}

	m.handle(m, msg)
}
