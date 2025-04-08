package actor

import "context"

type Context interface {
	Init()
	Self() Ref
	Sender() Ref
	Parent() Ref
	Manager() Manager
	Message() interface{}
	SetSelf(self Ref)
	SetSender(sender Ref)
	SetParent(parent Ref)
	SetManager(mgr Manager)
}

type BaseContext struct {
	context.Context
	self   Ref
	parent Ref
	sender Ref
	mgr    Manager
	msg    interface{}
}

var _ Context = (*BaseContext)(nil)

func NewBaseContext() *BaseContext {
	return &BaseContext{}
}

func (bc *BaseContext) Init() {
}

func (bc *BaseContext) Self() Ref {
	return bc.self
}

func (bc *BaseContext) Sender() Ref {
	return bc.sender
}

func (bc *BaseContext) Parent() Ref {
	return bc.parent
}

func (bc *BaseContext) Manager() Manager {
	return bc.mgr
}

func (bc *BaseContext) Message() interface{} {
	return bc.msg
}

func (bc *BaseContext) SetSelf(self Ref) {
	bc.self = self
}

func (bc *BaseContext) SetSender(sender Ref) {
	bc.sender = sender
}

func (bc *BaseContext) SetParent(parent Ref) {
	bc.parent = parent
}

func (bc *BaseContext) SetManager(mgr Manager) {
	bc.mgr = mgr
}

func (bc *BaseContext) SetMessage(msg interface{}) {
	bc.msg = msg
}
