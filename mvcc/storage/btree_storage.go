package storage

import (
	"errors"
	"sync"

	"github.com/google/btree"

	"github.com/beihai0xff/puff/internal/model"
)

type btreeStorage struct {
	sync.RWMutex
	tree *btree.BTree
}

func NewBtreeStorage() Storage {
	return &btreeStorage{
		tree: btree.New(100),
	}
}

func (t *btreeStorage) Set(key string, version uint64) error {

	if key == "" {
		return errors.New("empty string can not be added to btreeStorage")
	}

	entry := &model.Entry{Key: key, Version: version}

	t.Lock()
	defer t.Unlock()
	if t.tree.Has(entry) {
		return nil
	}
	t.tree.ReplaceOrInsert(entry)
	return nil
}

func (t *btreeStorage) Get(key string) (*model.Entry, error) {
	entry := &model.Entry{Key: key}
	t.RLock()
	item := t.tree.Get(entry)
	t.RUnlock()
	if item == nil {
		return nil, nil
	}
	return item.(*model.Entry), nil
}

func (t *btreeStorage) Delete(key string) error {
	entry := &model.Entry{Key: key}
	t.Lock()
	item := t.tree.Delete(entry)
	t.Unlock()
	if item == nil {
		// TODO: 删除不存在的 key 时向上层返回错误信息
	}

	return nil
}

func (t *btreeStorage) Range(start, end string) []model.Entry {
	greater := &model.Entry{Key: start}
	lessThan := &model.Entry{Key: end}
	var res []model.Entry
	t.RLock()
	defer t.RUnlock()
	t.tree.AscendRange(greater, lessThan, func(item btree.Item) bool {
		f := item.(*model.Entry)
		res = append(res, *f)
		return true
	})

	return res
}

// Backup 将某个时刻的存储引擎中的数据备份到磁盘中
func (t *btreeStorage) Backup() {
	_ = t.tree.Clone()

}
