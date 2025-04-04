package ctypes

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrRequestTimeout  = errors.New("request timeout")
	ErrNotFound        = errors.New("not found")
	ErrKeyRepeat       = errors.New("key repeat")
	ErrSetRepeat       = errors.New("set repeat")

	ErrActorClosed   = errors.New("operation on closed actor")
	ErrChannelFull   = errors.New("channel buffer full")
	ErrChildNotFound = errors.New("child not found")

	ErrKeyNotFound = errors.New("key not found")
	ErrInvalidType = errors.New("invalid value type")

	ErrHandlerTimeout = errors.New("execution timed out")
)
