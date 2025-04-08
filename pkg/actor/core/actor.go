package actor

import (
	"sync"
)

type Actor interface {
	Init()
	GetChild(id uint32) (Ref, bool)
	Receive(ctx Context, msg interface{})
	PreStart(ctx Context)
	PostStop(ctx Context)
	RegisterChild(child Ref)
	Mailbox() chan interface{}
}

type BaseActor struct {
	id         uint32
	self       Ref
	parent     Ref
	children   map[uint32]Ref
	childLock  sync.RWMutex
	mgr        Manager
	mailbox    chan interface{}
	stopped    bool
	restartCnt int
}

var _ Actor = (*BaseActor)(nil)

func NewBaseActor() *BaseActor {
	return &BaseActor{}
}

func (b *BaseActor) Init() {
	b.mailbox = make(chan interface{})
	b.children = make(map[uint32]Ref)
}

func (b *BaseActor) Self() Ref {
	return b.self
}

func (b *BaseActor) NewRef() Ref {
	return NewBaseRef(b.id, b.mailbox, b.parent)
}

func (b *BaseActor) Parent() Ref {
	return b.parent
}

func (b *BaseActor) Manager() Manager {
	return b.mgr
}

func (b *BaseActor) Mailbox() chan interface{} {
	return b.mailbox
}

func (b *BaseActor) NewMailbox(buffer int) chan interface{} {
	return make(chan interface{}, buffer)
}

func (b *BaseActor) SetParent(parent Ref) {
	b.parent = parent
}

func (b *BaseActor) SetManager(mgr Manager) {
	b.mgr = mgr
}

func (b *BaseActor) Receive(ctx Context, msg interface{}) {
}

func (b *BaseActor) PreStart(ctx Context) {
}

func (b *BaseActor) PostStop(ctx Context) {
}

func (b *BaseActor) RegisterChild(child Ref) {
	b.childLock.Lock()
	defer b.childLock.Unlock()
	b.children[child.Id()] = child
}

func (b *BaseActor) GetChild(id uint32) (Ref, bool) {
	b.childLock.RLock()
	defer b.childLock.RUnlock()
	child, ok := b.children[id]
	if ok {
		return child, true
	}
	return nil, false
}
