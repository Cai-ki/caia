package actor

import "context"

type Context interface {
	Self() Ref
	Sender() Ref
	Parent() Ref
	Manager() Manager
}

type BaseContext struct {
	context.Context
	self   Ref
	sender Ref
	mgr    Manager
}

func (ctx *BaseContext) Self() Ref {
	return ctx.self
}

func (ctx *BaseContext) Sender() Ref {
	return ctx.sender
}

func (bc *BaseContext) Parent() Ref {
	return bc.self.Parent()
}

func (bc *BaseContext) Manager() Manager {
	return bc.mgr
}
