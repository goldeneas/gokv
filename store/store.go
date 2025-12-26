package store

import "errors"

var (
	ErrItemNotFound = errors.New("item not found or expired")
)

type Store interface {
	Put(key string, value string, ttl int)
	Get(key string) (string, error)
	Delete(key string) error
}
