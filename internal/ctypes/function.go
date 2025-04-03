package ctypes

import "context"

type MethodFunc func(interface{}) interface{}

type HandleFunc func(context.Context, Message)

type OptionFunc func(interface{})
