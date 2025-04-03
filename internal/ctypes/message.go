package ctypes

import "time"

// Message 负责传输数据
type Message struct {
	Payload interface{}
	ReplyTo chan<- Message
}

func NewMessage(payload interface{}, replyTo chan<- Message) *Message {
	return &Message{
		Payload: payload,
		ReplyTo: replyTo,
	}
}

// Mailbox 管理消息接收和结果返回
type Mailbox struct {
	ch chan Message
}

func NewMailbox(buffer int) *Mailbox {
	return &Mailbox{
		ch: make(chan Message, buffer),
	}
}

// Result 表示处理结果
type Result struct {
	Data interface{}
	Err  error
}

//TODO 目前没想出一个合理的机制阻止向不属于自己的信箱中读取信息，目前只能靠自觉，可以抽像出邮箱接口，分别实现私人邮箱（只能接收），根据私人邮箱创建公共邮箱（只能发送）。

// 返回公共chan
func (m *Mailbox) PubChan() chan<- Message {
	return m.ch
}

// 返回私有chan
func (m *Mailbox) PriChan() <-chan Message {
	return m.ch
}

// 返回chan
func (m *Mailbox) Chan() chan Message {
	return m.ch
}

// 阻塞发送结果
func (m *Mailbox) SendResult(data interface{}, err error) {
	m.ch <- Message{Payload: Result{Data: data, Err: err}, ReplyTo: nil}
}

// 带超时的结果发送
func (m *Mailbox) SendResultTimeout(data interface{}, err error, timeout time.Duration) bool {
	select {
	case m.ch <- Message{Payload: Result{Data: data, Err: err}, ReplyTo: nil}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// 接收结果（阻塞）
func (m *Mailbox) Receive() Message {
	return <-m.ch
}

// 接收结果（带超时）
func (m *Mailbox) ReceiveTimeout(timeout time.Duration) (Message, bool) {
	select {
	case res := <-m.ch:
		return res, true
	case <-time.After(timeout):
		return Message{}, false
	}
}

// 关闭
func (m *Mailbox) Close() {
	close(m.ch)
}
