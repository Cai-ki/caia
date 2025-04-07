package actor

import (
	"sync"
)

type Actor interface {
	Receive(ctx Context, msg interface{})
	PreStart(ctx Context)
	PostStop(ctx Context)
}

type BaseActor struct {
	ctx       Context
	children  map[uint32]Ref
	childLock sync.Mutex
}

func (b *BaseActor) SetContext(ctx Context) {
	b.ctx = ctx
	b.children = make(map[uint32]Ref)
}

func (b *BaseActor) PreStart(ctx *BaseContext) {
}

func (b *BaseActor) PostStop(ctx *BaseContext) {
}

func (b *BaseActor) SpawnChild(ctx Context, id uint32, actor Actor) Ref {
	bctx := b.ctx.(*BaseContext)
	child := bctx.Manager().(*BaseManager).RegisterActor(ctx, id, actor, ctx.Self())
	b.childLock.Lock()
	b.children[id] = child
	b.childLock.Unlock()
	return child
}

func (b *BaseActor) Context() Context {
	return b.ctx
}
