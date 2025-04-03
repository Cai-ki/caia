package ctypes

type Storage interface {
	Get(routineName string, key string) (interface{}, error)
	Put(routineName string, key string, value interface{}) error
	Delete(routineName string, key string) error
	Exists(routineName string, key string) bool
}
