package storage

type Storage interface {
	Set(key string) error
	Get(key string) (string, error)
}
