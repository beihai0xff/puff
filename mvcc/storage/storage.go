package storage

import "github.com/beihai0xff/puff/internal/model"

type Storage interface {
	Set(key string, version uint64) error
	Get(key string) (*model.Entry, error)
	Delete(key string) error
	Range(start, end string) []model.Entry
	Backup()
}
