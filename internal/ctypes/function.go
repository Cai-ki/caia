package ctypes

type MethodFunc func(interface{}) interface{}

type HandleFunc func(Actor, Message)

type OptionFunc func(interface{})
