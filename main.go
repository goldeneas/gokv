package main

import (
	"log"

	"github.com/goldeneas/gokv/cluster"
	"github.com/goldeneas/gokv/store"
)

type ConcreteNode struct {
	id    string
	store store.Store
}

func NewConcreteNode(id string) *ConcreteNode {
	return &ConcreteNode{
		id:    id,
		store: store.NewMapStore(),
	}
}

func (c *ConcreteNode) Store() store.Store {
	return c.store
}

func (c *ConcreteNode) Identifier() string {
	return c.id
}

func main() {
	cluster := cluster.NewCluster()

	nodeMontoro := NewConcreteNode("montoro.com")
	nodeSalerno := NewConcreteNode("salerno.com")

	cluster.Add(nodeMontoro)
	cluster.Add(nodeSalerno)

	cluster.Put("test", "nicola")
	cluster.Put("picarella", "10")

	value, err := cluster.Get("test")
	if err != nil {
		log.Printf("got err %s", err)
	} else {
		log.Printf("got value %s", value)
	}
}
