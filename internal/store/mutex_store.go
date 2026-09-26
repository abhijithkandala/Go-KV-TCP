package store

import (
	"sync"
)

type MutexStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMutexStore() *MutexStore {
	return &MutexStore{data: make(map[string]string)}
}

func (s *MutexStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value

}

func (s *MutexStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, isPresent := s.data[key]
	if !isPresent {
		return "", ErrEntryNotFound
	}
	return val, nil
}
func (s *MutexStore) Del(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, isPresent := s.data[key]
	if !isPresent {
		return ErrEntryNotFound
	}
	delete(s.data, key)
	return nil
}

var _ Store = (*MutexStore)(nil)
