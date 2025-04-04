package ctlv

// import (
// 	"encoding/binary"
// 	"errors"

// 	"github.com/Cai-ki/caia/internal/ctypes"
// 	"github.com/Cai-ki/caia/pkg/cprotocol"
// )

// type Data struct {
// 	Type   uint32
// 	Length uint32
// 	Body   []byte
// }

// func NewData() *Data {
// 	return &Data{
// 		Body: make([]byte, 0),
// 	}
// }

// type TLVEncoder struct{}

// type TLVDecoder struct {
// 	buffer []byte
// }

// var _ cprotocol.Encoder = (*TLVEncoder)(nil)
// var _ cprotocol.Decoder = (*TLVDecoder)(nil)

// func NewTLVDecoder(buffer int) *TLVDecoder {
// 	return &TLVDecoder{
// 		buffer: make([]byte, 0, buffer),
// 	}
// }

// // Encode encodes a Data structure into a byte slice.
// func (t *TLVEncoder) Encode(v interface{}) ([]byte, error) {
// 	data, ok := v.(*Data)
// 	if !ok {
// 		return nil, ctypes.ErrInvalidType
// 	}

// 	// Preallocate space for Type (4 bytes), Length (4 bytes), and Body (variable length).
// 	bytes := make([]byte, 8+len(data.Body))

// 	// Write Type and Length to the correct positions.
// 	binary.BigEndian.PutUint32(bytes[0:4], data.Type)
// 	binary.BigEndian.PutUint32(bytes[4:8], data.Length)

// 	// Copy Body data.
// 	copy(bytes[8:], data.Body)

// 	return bytes, nil
// }

// // Decode decodes a byte slice into a Data structure or a slice of Data structures.
// func (t *TLVDecoder) Decode(bytes []byte, v interface{}) error {
// 	t.buffer = append(t.buffer, bytes...)

// 	for len(t.buffer) >= 8 {
// 		// Extract Type and Length from the buffer.
// 		tp := binary.BigEndian.Uint32(t.buffer[:4])
// 		length := binary.BigEndian.Uint32(t.buffer[4:8])

// 		// Check if the buffer contains enough data for the Body.
// 		if uint32(len(t.buffer)-8) < length {
// 			return errors.New("bytes not enough")
// 		}

// 		switch data := v.(type) {
// 		case *Data:
// 			// Allocate space for Body and copy data.
// 			data.Type = tp
// 			data.Length = length
// 			data.Body = make([]byte, length)
// 			copy(data.Body, t.buffer[8:8+length])

// 			// Update the buffer.
// 			t.buffer = t.buffer[8+length:]
// 			return nil

// 		case *[]*Data:
// 			// Allocate space for Body and create a new Data item.
// 			item := &Data{
// 				Type:   tp,
// 				Length: length,
// 				Body:   make([]byte, length),
// 			}
// 			copy(item.Body, t.buffer[8:8+length])

// 			// Append the item to the slice.
// 			*data = append(*data, item)

// 			// Update the buffer.
// 			t.buffer = t.buffer[8+length:]

// 		default:
// 			return ctypes.ErrInvalidType
// 		}
// 	}

// 	// If we exit the loop, there are not enough bytes to decode.
// 	return errors.New("bytes not enough")
// }
