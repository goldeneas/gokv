package store

import (
	"sync"
	"time"
)

type MapStore struct {
	data map[string]item
	mtx  sync.RWMutex
}

type item struct {
	value     string
	expiresAt time.Time
}

func NewMapStore() *MapStore {
	return &MapStore{
		data: make(map[string]item),
	}
}

func (s *MapStore) Put(key string, value string, ttl int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	item := item{
		value:     value,
		expiresAt: time.Now().Add(time.Duration(ttl) * time.Second),
	}

	s.data[key] = item
}

func (s *MapStore) Get(key string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	it, ok := s.data[key]
	if !ok || time.Now().After(it.expiresAt) {
		s.mtx.RUnlock()
		s.mtx.Lock()
		delete(s.data, key)
		s.mtx.Unlock()
		s.mtx.RLock()

		return "", ErrItemNotFound
	}

	return it.value, nil
}

func (s *MapStore) Delete(key string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, exists := s.data[key]
	if !exists {
		return ErrItemNotFound
	}

	delete(s.data, key)
	return nil
}
