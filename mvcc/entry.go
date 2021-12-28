package mvcc

import (
	"github.com/google/btree"
)

type Entry struct {
	Version uint64
	Key     string
}

func (m *Entry) Less(item btree.Item) bool {
	return m.Key < (item.(*Entry)).Key
}
