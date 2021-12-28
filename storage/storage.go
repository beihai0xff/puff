package storage

type Storage interface {
	Set(entry Entry) error
	Get(key string) (Entry, error)
}
