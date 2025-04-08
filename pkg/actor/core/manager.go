package actor

import (
	"sync"
)

type Manager interface {
	RegisterActor(ctx Context, id uint32, actor Actor)
	runActor(ctx Context, actor Actor)
}

type BaseManager struct {
	actors map[uint32]Ref
	mu     sync.RWMutex
}

var _ Manager = (*BaseManager)(nil)

func NewBaseManager() *BaseManager {
	return &BaseManager{
		actors: make(map[uint32]Ref),
	}
}

func (m *BaseManager) RegisterActor(ctx Context, id uint32, actor Actor) {
	m.mu.Lock()
	m.actors[id] = ctx.Self()
	m.mu.Unlock()

	go m.runActor(ctx, actor)
}

func (m *BaseManager) runActor(ctx Context, actor Actor) {
	actor.PreStart(ctx)
	defer actor.PostStop(ctx)

	mailbox := actor.Mailbox()
	for {
		select {
		case msg := <-mailbox:
			if env, ok := msg.(Envelope); ok {
				ctx.SetSender(env.Sender)
				actor.Receive(ctx, env.Msg)
			} else {
				ctx.SetSender(nil)
				actor.Receive(ctx, msg)
			}
		case <-ctx.(*BaseContext).Done():
			return
		}
	}
}
