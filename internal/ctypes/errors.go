package ctypes

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrRequestTimeout  = errors.New("request timeout")
	ErrNotFound        = errors.New("not found")
	ErrKeyRepeat       = errors.New("key repeat")
	ErrSetRepeat       = errors.New("Set repeat")

	ErrRoutineClosed = errors.New("routine: operation on closed routine")
	ErrChannelFull   = errors.New("routine: channel buffer full")

	ErrKeyNotFound = errors.New("storage: key not found")
	ErrInvalidType = errors.New("storage: invalid value type")

	ErrHandlerTimeout = errors.New("handler: execution timed out")
)
