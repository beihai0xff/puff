package storage

import "github.com/beihai0xff/puff/mvcc"

type Storage interface {
	Set(key string, version uint64) error
	Get(key string) (*mvcc.Entry, error)
	Delete(key string) error
	Range(start, end string) []mvcc.Entry
	Backup()
}
