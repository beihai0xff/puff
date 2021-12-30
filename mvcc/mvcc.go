package model

import (
	"sync"

	"github.com/beihai0xff/puff/internal/model"
	"github.com/beihai0xff/puff/mvcc/storage"
)

type StateMachine interface {
	Set(key string) error
	Get(key string) *model.Entry
	Delete(key string) error
	Range(start, end string) ([]*model.Entry, error)
	Backup() StateMachineStatus
}

type State struct {
	sync.Mutex
	kvStorage storage.Storage
	version   uint64

	isDumping bool
}

func (s *State) Set(key string) error {
	if key == "" {

	}
	var version uint64
	s.Lock()
	s.version++
	version = s.version
	s.Unlock()
	return s.kvStorage.Set(key, version)
}

func (s *State) Backup() error {
	if s.isDumping {
		return ErrStateMachineIsDumping
	}
	s.isDumping = true

	s.kvStorage.Backup()

	s.isDumping = false

	return nil
}
