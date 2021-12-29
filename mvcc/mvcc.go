package mvcc

import (
	"sync"

	"github.com/beihai0xff/puff/mvcc/storage"
)

type StateMachine interface {
	Set(key string) error
	Get(key string) *Entry
	Delete(key string) error
	Range(start, end string) ([]*Entry, error)
	Backup() StateMachineStatus
}

type State struct {
	sync.RWMutex
	kvStorage storage.Storage
	version   uint64

	isDumping bool
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
