package ctypes

type Actor interface {
	CreateChild(name string, buffer int, handle HandleFunc) (Actor, error)
	GetName() string
	GetMailbox() Mailbox
	Start()
	Stop()
	SendMessage(msg Message) error
	SendMessageAsync(msg Message) error
	SendMessageToChildren(msg Message)
	SendMessageAsyncToChildren(msg Message)
}
