package storage

import "github.com/beihai0xff/puff/mvcc"

type Storage interface {
	set(key string) error
	Get(key string) (*mvcc.Entry, error)
	Delete(key string) error
	Range(start, end string) []mvcc.Entry
	Backup()
}
