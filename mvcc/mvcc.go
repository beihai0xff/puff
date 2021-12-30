package model

import (
	"sync"

	"github.com/beihai0xff/puff/internal/model"
	"github.com/beihai0xff/puff/mvcc/storage"
)

type StateMachine interface {
	Set(key string) error
	Get(key string) (*model.Entry, error)
	Delete(key string) error
	Range(start, end string) ([]model.Entry, error)
	Backup() error
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

func (s *State) Get(key string) (*model.Entry, error) {
	if key == "" {

	}
	return s.kvStorage.Get(key)
}

func (s *State) Delete(key string) error {
	if key == "" {

	}
	return s.kvStorage.Delete(key)
}

func (s *State) Range(start, end string) ([]model.Entry, error) {
	if start == "" || end == "" {

	}
	return s.kvStorage.Range(start, end), nil
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
