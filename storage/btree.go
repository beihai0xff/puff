package storage

import (
	"errors"

	"github.com/google/btree"
)

type btreeStorage struct {
	tree *btree.BTree
}

func newBtreeStorage() *btreeStorage {
	return &btreeStorage{
		tree: btree.New(100),
	}
}

func (t *btreeStorage) set(entry *Entry) error {
	if entry == nil {
		return errors.New("nil entry can not be added to btreeStorage")
	}
	t.tree.ReplaceOrInsert(entry)
	return nil
}

func (t *btreeStorage) Get(entry *Entry) (*Entry, error) {
	item := t.tree.Get(entry)
	if item == nil {
		return nil, nil
	}
	return item.(*Entry), nil
}
