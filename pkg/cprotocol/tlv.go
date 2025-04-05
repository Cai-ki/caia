package cprotocol

import (
	"encoding/binary"

	"github.com/Cai-ki/caia/internal/clog"
)

// TLV实现示例
type TLVPacket struct {
	Type   uint32
	Length uint32
	Value  []byte
}

func (p *TLVPacket) Reset() {
	p.Type = 0
	p.Length = 0
	p.Value = p.Value[:0]
}

type TLVHandler struct {
	codec       *Codec
	MaxBodySize int
}

func (h *TLVHandler) SetCodec(codec *Codec) {
	h.codec = codec
	h.MaxBodySize = 1024
}

func (h *TLVHandler) Encode(v interface{}) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			clog.Error(err)
		}
	}()
	pkt, ok := v.(*TLVPacket)
	if !ok {
		return nil, ErrInvalidData
	}

	// if len(pkt.Value) > h.MaxBodySize {
	// 	return nil, ErrInvalidData
	// }

	buf := h.codec.GetBuffer()
	if len(buf) < int(8+len(pkt.Value)) {
		buf = make([]byte, 8+len(pkt.Value))
	}
	binary.BigEndian.PutUint32(buf[0:4], pkt.Type)
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(pkt.Value)))
	copy(buf[8:], pkt.Value)
	return buf, nil
}

func (h *TLVHandler) Decode(data []byte) (interface{}, int, error) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	if len(data) < 8 {
		return nil, 0, ErrInvalidData
	}

	pkt := h.codec.Acquire().(*TLVPacket)
	pkt.Type = binary.BigEndian.Uint32(data[0:4])
	pkt.Length = binary.BigEndian.Uint32(data[4:8])

	// if int(pkt.Length) > h.MaxBodySize {
	// 	return nil, 0, ErrInvalidData
	// }

	totalSize := 8 + int(pkt.Length)
	if len(data) < totalSize {
		return nil, 0, ErrInvalidData
	}

	pkt.Value = append(pkt.Value[:0], data[8:totalSize]...)
	data = data[totalSize:]
	return pkt, totalSize, nil
}
