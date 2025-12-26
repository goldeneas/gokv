package store

import (
	"errors"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	store := NewMapStore()

	var putValues = []struct {
		key   string
		value string
		ttl   int
	}{
		{"TEST1", "VALUE1", 0},
		{"TEST2", "VALUE2", 60},
	}

	for _, entry := range putValues {
		store.Put(entry.key, entry.value, entry.ttl)
	}

	var tests = []struct {
		key        string
		want_left  string
		want_right error
	}{
		{"TEST1", "", ErrItemNotFound},
		{"TEST2", "VALUE2", nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s,%s", tt.key, tt.want_left, tt.want_right)
		t.Run(testname, func(t *testing.T) {
			left, right := store.Get(tt.key)
			if left != tt.want_left || !errors.Is(right, tt.want_right) {
				t.Errorf("got '%s' and '%s', want '%s' and '%s'",
					left, right, tt.want_left, tt.want_right)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	store := NewMapStore()

	var putValues = []struct {
		key   string
		value string
		ttl   int
	}{
		{"TEST1", "VALUE1", 0},
		{"TEST2", "VALUE2", 60},
	}

	for _, entry := range putValues {
		store.Put(entry.key, entry.value, entry.ttl)
	}

	var tests = []struct {
		key  string
		want error
	}{
		{"TEST1", nil},
		{"TEST2", nil},
		{"TEST3", ErrItemNotFound},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.key, tt.want)
		t.Run(testname, func(t *testing.T) {
			ans := store.Delete(tt.key)
			if !errors.Is(ans, tt.want) {
				t.Errorf("got '%s', want '%s'", ans, tt.want)
			}
		})
	}
}
