package ctypes

import "sync"

// Registry 定义服务注册发现接口
type Registry interface {
	RegisterMethod(name string, method MethodFunc) error
	RegisterHandle(name string, handle HandleFunc) error
	RegisterPool(name string, pool *sync.Pool) error
	GetMethod(name string) (MethodFunc, error)
	GetHandle(name string) (HandleFunc, error)
	GetPool(name string) (*sync.Pool, error)
}
