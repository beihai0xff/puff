package storage

import (
	"errors"
	"sync"

	"github.com/google/btree"

	"github.com/beihai0xff/puff/mvcc"
)

type btreeStorage struct {
	sync.RWMutex
	tree *btree.BTree
}

func newBtreeStorage() Storage {
	return &btreeStorage{
		tree: btree.New(100),
	}
}

func (t *btreeStorage) set(key string) error {

	if key == "" {
		return errors.New("empty string can not be added to btreeStorage")
	}

	entry := &mvcc.Entry{Key: key}

	t.Lock()
	defer t.Unlock()
	if t.tree.Has(entry) {
		return nil
	}
	t.tree.ReplaceOrInsert(entry)
	return nil
}

func (t *btreeStorage) Get(key string) (*mvcc.Entry, error) {
	entry := &mvcc.Entry{Key: key}
	t.RLock()
	item := t.tree.Get(entry)
	t.RUnlock()
	if item == nil {
		return nil, nil
	}
	return item.(*mvcc.Entry), nil
}

func (t *btreeStorage) Delete(key string) error {
	entry := &mvcc.Entry{Key: key}
	t.Lock()
	item := t.tree.Delete(entry)
	t.Unlock()
	if item == nil {
		// TODO: 删除不存在的 key 时向上层返回错误信息
	}

	return nil
}

func (t *btreeStorage) Range(start, end string) []mvcc.Entry {
	greater := &mvcc.Entry{Key: start}
	lessThan := &mvcc.Entry{Key: end}
	var res []mvcc.Entry
	t.RLock()
	defer t.RUnlock()
	t.tree.AscendRange(greater, lessThan, func(item btree.Item) bool {
		f := item.(*mvcc.Entry)
		res = append(res, *f)
		return true
	})

	return res
}

func (t *btreeStorage) Backup() {

}
