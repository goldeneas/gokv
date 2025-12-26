package cluster

import (
	"testing"
	"time"

	"github.com/goldeneas/gokv/hashring"
	"github.com/goldeneas/gokv/store"
)

type ConcreteNode struct {
	id    string
	store store.Store
}

func (n *ConcreteNode) Identifier() string {
	return n.id
}

func (n *ConcreteNode) Store() store.Store {
	return n.store
}

func TestClusterFullFlow(t *testing.T) {
	c := NewCluster(
		hashring.SetEnableLogs(true),
	)

	nodes := []*ConcreteNode{
		{id: "node-1", store: store.NewMapStore()},
		{id: "node-2", store: store.NewMapStore()},
		{id: "node-3", store: store.NewMapStore()},
	}

	for _, n := range nodes {
		if err := c.Add(n); err != nil {
			t.Fatalf("unable to add node %s: %v", n.id, err)
		}
	}

	t.Run("verify put and get", func(t *testing.T) {
		key := "session_123"
		val := "active"

		err := c.Put(key, val, 10)
		if err != nil {
			t.Errorf("put error: %v", err)
		}

		got, err := c.Get(key)
		if err != nil {
			t.Errorf("get error: %v", err)
		}
		if got != val {
			t.Errorf("wrong value: expected %s, got %s", val, got)
		}
	})

	t.Run("verify ttl (expiration)", func(t *testing.T) {
		key := "temp_data"
		val := "burn_after_reading"

		c.Put(key, val, 1)

		if v, _ := c.Get(key); v != val {
			t.Fatal("data should be present")
		}

		t.Logf("waiting for expiration...")
		time.Sleep(1200 * time.Millisecond)

		_, err := c.Get(key)
		if err != store.ErrItemNotFound {
			t.Errorf("expected ErrItemNotFound error, got: %v", err)
		}
	})

	t.Run("verify deletion", func(t *testing.T) {
		key := "delete_me"
		c.Put(key, "data", 60)

		if err := c.Delete(key); err != nil {
			t.Errorf("delete error: %v", err)
		}

		_, err := c.Get(key)
		if err == nil {
			t.Error("the key should have been deleted")
		}
	})

	t.Run("verify distribution (consistency)", func(t *testing.T) {
		key := "consistent_key"
		c.Put(key, "value", 60)

		nodeInterface, _ := c.ring.Get(key)
		targetNode := nodeInterface.(*ConcreteNode)

		val, err := targetNode.Store().Get(key)
		if err != nil || val != "value" {
			t.Errorf("data not found in the node predicted by the hashring")
		}
	})
}
