package ctypes

type Factory interface {
	RegisterFactory(name string, factory HandleFactoryFunc)
	GetFactory(name string) (HandleFactoryFunc, bool)
}
