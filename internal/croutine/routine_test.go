package croutine

import (
	"context"
	"fmt"
	"testing"

	"github.com/Cai-ki/caia/internal/ctypes"
)

func TestManagerNormal(t *testing.T) {
	mailbox := make(chan ctypes.Message, 5)
	msgs := [...]ctypes.Message{
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
	}
	r := NewManager("test", 5, context.Background(), func(ctx context.Context, msg ctypes.Message) {
		msg.ReplyTo <- *ctypes.NewMessage(msg.Payload, nil)
	})

	r.Start()

	for _, msg := range msgs {
		r.SendMessage(msg)
	}

	for _, msg := range msgs {
		err := r.SendMessageAsync(msg)
		if err == nil {
			t.Errorf("err = %v; want %v", err, fmt.Errorf("channel is full"))
		}
	}

	r.Stop()

	if len(mailbox) != len(msgs) {
		t.Errorf("len(mailbox) = %d; want %d", len(mailbox), len(msgs))
	}
}

func TestManagerWhenPanic(t *testing.T) {
	mailbox := make(chan ctypes.Message, 10)
	msgs := [...]ctypes.Message{
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
		*ctypes.NewMessage("msg", mailbox),
	}
	r := NewManager("test", 10, context.Background(), func(ctx context.Context, msg ctypes.Message) {
		msg.ReplyTo <- *ctypes.NewMessage(msg.Payload, nil)
		panic("panic")
	})

	r.Start()

	for _, msg := range msgs {
		r.SendMessage(msg)
	}

	r.Stop()
}
