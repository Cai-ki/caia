package cactor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Cai-ki/caia/internal/ctypes"
)

const (
	KeyManager = "Manager"
)

func WithValue(key, val interface{}) ctypes.OptionFunc {
	return func(i interface{}) {
		m := i.(*Manager)
		m.ctx = context.WithValue(m.ctx, key, val)
	}
}

type Manager struct {
	name    string
	mailbox ctypes.Mailbox
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	//children    []*Manager
	children map[string]*Manager
	mu       sync.Mutex
	closed   atomic.Bool
	handle   ctypes.HandleFunc
}

var _ ctypes.Actor = (*Manager)(nil)

func NewManager(name string, buffer int, parentCtx context.Context, handle ctypes.HandleFunc, opts ...ctypes.OptionFunc) *Manager {
	ctx, cancel := context.WithCancel(parentCtx)
	m := &Manager{
		name:     name,
		mailbox:  *ctypes.NewMailbox(buffer),
		ctx:      ctx,
		cancel:   cancel,
		children: map[string]*Manager{},
		handle:   handle,
	}

	WithValue(KeyManager, m)(m)

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (r *Manager) CreateChild(name string, buffer int, handle ctypes.HandleFunc, opts ...ctypes.OptionFunc) (ctypes.Actor, error) {
	if r.closed.Load() {
		panic("can't create child when actor running")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.children[name]
	if ok {
		return nil, ctypes.ErrKeyRepeat
	}
	child := NewManager(name, buffer, r.ctx, handle)
	r.children[name] = child

	WithValue(KeyManager, child)(child)

	for _, opt := range opts {
		opt(child)
	}
	return child, nil
}

func (r *Manager) GetName() string {
	return r.name
}

func (r *Manager) GetMailbox() ctypes.Mailbox {
	return r.mailbox
}

func (r *Manager) Start() {
	r.start()
}

func (r *Manager) Stop() {
	r.stop()
}

func (r *Manager) SendMessage(msg ctypes.Message) error {
	select {
	case r.mailbox.Chan() <- msg:
		return nil
	case <-r.ctx.Done():
		return ctypes.ErrActorClosed
	}
}

func (r *Manager) SendMessageAsync(msg ctypes.Message) error {
	select {
	case r.mailbox.Chan() <- msg:
		return nil
	case <-r.ctx.Done():
		return ctypes.ErrActorClosed
	default:
		return ctypes.ErrChannelFull
	}
}

func (r *Manager) SendMessageToChildren(msg ctypes.Message) {
	for _, child := range r.children {
		if child != nil {
			child.SendMessage(msg)
		}
	}
}

func (r *Manager) SendMessageAsyncToChildren(msg ctypes.Message) {
	for _, child := range r.children {
		if child != nil {
			child.SendMessageAsync(msg)
		}
	}
}

func (r *Manager) start() {
	r.wg.Add(1)
	r.serve() //TODO 这里可能有协程启动时不符合预期父协程先，子协程后的顺序。建议之后压测一下。

	for _, child := range r.children {
		child.start()
	}
}

func (r *Manager) serve() {
	go func() {
		defer r.stop()
		defer r.wg.Done()

		for {
			select {
			case msg := <-r.mailbox.Chan():
				r.handleMessage(msg)
			case <-r.ctx.Done():
				if r.closed.Load() { // 正常退出流程，保证子协程正常退出后父协程退出。
					for {
						select {
						case msg := <-r.mailbox.Chan():
							r.handleMessage(msg)
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

func (r *Manager) cleanup() {
	r.mailbox.Close()
	fmt.Printf("actor: %s close channel\n", r.name)
}

func (r *Manager) stop() {
	if r.closed.Load() {
		return
	}
	r.closed.Store(true)

	for _, child := range r.children {
		child.stop()
	}

	r.cancel()
	r.wg.Wait()

	r.cleanup()
}

func (r *Manager) handleMessage(msg ctypes.Message) {
	defer func() { // panic终止于此处，默认不信任所有执行的handle，Manager只负责接受message，并根据构造时传入的handle进行相应处理，进行一个逻辑转发的工作。
		if err := recover(); err != nil {
			fmt.Printf("actor: %s panic when handle message: %v with error: %v\n", r.name, msg, err)
		}
	}()

	if r.handle == nil {
		fmt.Printf("actor: %s handle not found\n", r.name)
		return
	}

	r.handle(r.ctx, msg)
}
