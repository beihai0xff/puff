package mvcc

import (
	"sync"

	"github.com/beihai0xff/puff/mvcc/storage"
)

type MVCC interface {
	set(key string) error
	Get(key string) *Entry
	Delete(key string) error
	Range(start, end string) ([]*Entry, error)
	Backup()
}

type store struct {
	sync.RWMutex
	kvStorage storage.Storage
}
