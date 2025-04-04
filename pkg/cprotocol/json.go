package cprotocol

import "encoding/json"

type JSONHandler struct{}

func (h *JSONHandler) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (h *JSONHandler) Decode(data []byte) (interface{}, int, error) {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, 0, err
	}
	return &msg, len(data), nil
}
