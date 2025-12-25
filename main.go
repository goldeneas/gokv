package gokv

import (
	"github.com/goldeneas/gokv/hashring"
)

type ConcreteNode struct {
	id string
}

func NewConcreteNode(id string) *ConcreteNode {
	return &ConcreteNode{
		id: id,
	}
}

func (c *ConcreteNode) Identifier() string {
	return c.id
}

func main() {
	hr := hashring.NewHashRing()
	n1 := NewConcreteNode("TESTID1")

	hr.Add(n1)
}
