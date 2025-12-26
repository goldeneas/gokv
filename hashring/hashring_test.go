package hashring

import (
	"errors"
	"fmt"
	"testing"

	"github.com/goldeneas/gokv/store"
)

type ConcreteNode struct {
	id    string
	store *store.MapStore
}

func NewConcreteNode(id string) *ConcreteNode {
	return &ConcreteNode{
		id:    id,
		store: store.NewMapStore(),
	}
}

func (c *ConcreteNode) Identifier() string {
	return c.id
}

func (c *ConcreteNode) Store() store.Store {
	return c.store
}

func TestAdd(t *testing.T) {
	ring := NewHashRing()

	var tests = []struct {
		id   string
		want error
	}{
		{"TEST1", nil},
		{"TEST2", nil},
		{"TEST1", ErrNodeExists},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.id, tt.want)
		t.Run(testname, func(t *testing.T) {
			node := NewConcreteNode(tt.id)
			ans := ring.Add(node)

			if !errors.Is(ans, tt.want) {
				t.Errorf("got '%s', want '%s'", ans, tt.want)
			}
		})
	}

}

func TestGet(t *testing.T) {
	ring := NewHashRing()

	node := NewConcreteNode("TEST2")
	ring.Add(node)

	var tests = []struct {
		id   string
		want error
	}{
		{"TEST1", nil},
		{"TEST2", nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.id, tt.want)
		t.Run(testname, func(t *testing.T) {
			_, ans := ring.Get(tt.id)

			if !errors.Is(ans, tt.want) {
				t.Errorf("got '%s', want '%s'", ans, tt.want)
			}
		})
	}
}

func TestGet_EmptyRing(t *testing.T) {
	ring := NewHashRing()
	want := ErrNoConnectedNodes

	_, ans := ring.Get("ANYKEY")
	if ans != want {
		t.Errorf("got '%s', want '%s'", ans, want)
	}
}

func TestRemove(t *testing.T) {
	ring := NewHashRing()

	node := NewConcreteNode("TEST2")
	ring.Add(node)

	var tests = []struct {
		id   string
		want error
	}{
		{"TEST1", ErrNodeNotFound},
		{"TEST2", nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.id, tt.want)
		t.Run(testname, func(t *testing.T) {
			ans := ring.Remove(tt.id)

			if !errors.Is(ans, tt.want) {
				t.Errorf("got '%s', want '%s'", ans, tt.want)
			}
		})
	}
}
