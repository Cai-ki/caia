package ctypes

import (
	"context"
)

type Actor interface {
	CreateChild(name string, buffer int, handle HandleFunc, opts ...OptionFunc) (Actor, error)
	DeleteChild(name string) error
	GetName() string
	GetMailbox() Mailbox
	GetContext() context.Context
	GetParent() Actor
	GetChildren() map[string]Actor
	Start()
	Stop()
	StopWithErase()
	SendMessage(msg Message) error
	SendMessageAsync(msg Message) error
	SendMessageToChildren(msg Message)
	SendMessageAsyncToChildren(msg Message)
}
