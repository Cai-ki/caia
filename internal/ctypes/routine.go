package ctypes

type Routine interface {
	CreateChild(name string, buffer int, handle HandleFunc) (Routine, error)
	GetName() string
	GetMailbox() Mailbox
	SetHandle(handle HandleFunc) error
	Start()
	Stop()
	SendMessage(msg Message) error
	SendMessageAsync(msg Message) error
	SendMessageToChildren(msg Message)
	SendMessageAsyncToChildren(msg Message)
}
