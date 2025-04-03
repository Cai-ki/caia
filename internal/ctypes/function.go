package ctypes

type MethodFunc func(interface{}) interface{}

type HandleFunc func(Message)

type HandleFactoryFunc func(Registry, Storage) func(Message)
