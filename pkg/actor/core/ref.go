package actor

type Ref interface {
	Init()
	Id() uint32
	Mailbox() chan interface{}
	Parent() Ref
	Send(msg interface{})
}

type BaseRef struct {
	id      uint32
	mailbox chan interface{}
	parent  Ref
}

var _ Ref = (*BaseRef)(nil)

func NewBaseRef(id uint32, mailbox chan interface{}, parent Ref) *BaseRef {
	return &BaseRef{id: id, mailbox: mailbox, parent: parent}
}

func (br *BaseRef) Init() {}

func (br *BaseRef) Id() uint32 {
	return br.id
}

func (br *BaseRef) Mailbox() chan interface{} {
	return br.mailbox
}

func (br *BaseRef) Parent() Ref {
	return br.parent
}

func (br *BaseRef) Send(msg interface{}) {
	br.mailbox <- msg
}
