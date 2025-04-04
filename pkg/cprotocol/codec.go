package cprotocol

import (
	"errors"
	"sync"
)

var (
	ErrInvalidData    = errors.New("invalid protocol data")
	ErrHandlerMissing = errors.New("handler not registered")
)

// 核心接口定义
type CodecHandler interface {
	SetCodec(*Codec)
	Encode(interface{}) ([]byte, error)
	Decode([]byte) (interface{}, int, error)
}

type PooledObject interface {
	Reset()
}

// 通用编解码器
type Codec struct {
	handler    CodecHandler
	objectPool *sync.Pool
	bufferPool *sync.Pool
	mu         sync.Mutex
	buf        []byte
}

func NewCodec(handler CodecHandler, poolFunc func() interface{}) *Codec {
	c := &Codec{
		handler: handler,
		objectPool: &sync.Pool{
			New: poolFunc,
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024)
			},
		},
	}
	c.handler.SetCodec(c)
	return c
}

func (c *Codec) Acquire() interface{} {
	obj := c.objectPool.Get().(PooledObject)
	obj.Reset()
	return obj
}

func (c *Codec) Release(obj interface{}) {
	c.objectPool.Put(obj)
}

func (c *Codec) Encode(v interface{}) ([]byte, error) {
	if c.handler == nil {
		return nil, ErrHandlerMissing
	}
	return c.handler.Encode(v)
}

func (c *Codec) Decode(data []byte) (interface{}, int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 使用缓冲区处理粘包
	c.buf = append(c.buf, data...)

	return c.handler.Decode(c.buf)
}

func (c *Codec) GetBuffer() []byte {
	return c.bufferPool.Get().([]byte)
}

func (c *Codec) PutBuffer(b []byte) {
	if len(b) > 1024 {
		b = b[:1024]
	}
	c.bufferPool.Put(b)
}
