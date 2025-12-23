package store

import (
	"errors"
	"sync"
	"time"
)

type Store struct {
	data map[string]item
	mtx  sync.RWMutex
}

type item struct {
	value     string
	expiresAt time.Time
}

func (s *Store) Put(key string, value string, ttl int) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	item := item{
		value:     value,
		expiresAt: time.Now().Add(time.Duration(ttl) * time.Second),
	}

	s.data[key] = item
}

func (s *Store) Get(key string) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	it, ok := s.data[key]
	if !ok || time.Now().After(it.expiresAt) {
		s.mtx.RUnlock()
		s.mtx.Lock()
		delete(s.data, key)
		s.mtx.Unlock()
		s.mtx.RLock()

		return "", errors.New("item not found or expired")
	}

	return it.value, nil
}

func (s *Store) Delete(key string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, exists := s.data[key]
	if !exists {
		return errors.New("no entry found with key")
	}

	delete(s.data, key)
	return nil
}
