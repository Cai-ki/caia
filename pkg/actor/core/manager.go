package actor

import (
	"sync"
)

type Manager interface {
	RegisterActor(ctx Context, id uint32, actor Actor, parent Ref) Ref
	runActor(ctx Context, actor Actor, ref Ref)
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

func (m *BaseManager) RegisterActor(ctx Context, id uint32, actor Actor, parent Ref) Ref {
	actorRef := ctx.Self()

	m.mu.Lock()
	m.actors[id] = actorRef
	m.mu.Unlock()

	go m.runActor(ctx, actor, actorRef)
	return actorRef
}

func (m *BaseManager) runActor(ctx Context, actor Actor, ref Ref) {

	actor.PreStart(ctx)
	for msg := range ref.Mailbox() {
		if env, ok := msg.(Envelope); ok {
			ctx.(*BaseContext).sender = env.Sender
			actor.Receive(ctx, env.Msg)
		} else {
			ctx.(*BaseContext).sender = nil
			actor.Receive(ctx, msg)
		}
	}
	actor.PostStop(ctx)
}
