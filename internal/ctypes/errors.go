package ctypes

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrRequestTimeout  = errors.New("request timeout")
	ErrNotFound        = errors.New("not found")
	ErrKeyRepeat       = errors.New("key repeat")
	ErrSetRepeat       = errors.New("set repeat")

	ErrActorClosed   = errors.New("actor: operation on closed actor")
	ErrChannelFull   = errors.New("actor: channel buffer full")
	ErrChildNotFound = errors.New("actor: child not found")

	ErrKeyNotFound = errors.New("storage: key not found")
	ErrInvalidType = errors.New("storage: invalid value type")

	ErrHandlerTimeout = errors.New("handler: execution timed out")
)
