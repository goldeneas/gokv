package cluster

import (
	"github.com/goldeneas/gokv/hashring"
)

type Cluster struct {
	ring *hashring.HashRing
}

func NewCluster(opts ...hashring.ConfigFn) *Cluster {
	return &Cluster{
		ring: hashring.NewHashRing(opts...),
	}
}

func (c *Cluster) Add(node hashring.Node) error {
	return c.ring.Add(node)
}

func (c *Cluster) Remove(key string) error {
	return c.ring.Remove(key)
}

func (c *Cluster) Put(key string, value string, ttl int) error {
	node, err := c.ring.Get(key)
	if err != nil {
		return err
	}

	node.Store().Put(key, value, ttl)
	return nil
}

func (c *Cluster) Delete(key string) error {
	node, err := c.ring.Get(key)
	if err != nil {
		return err
	}

	return node.Store().Delete(key)
}

func (c *Cluster) Get(key string) (string, error) {
	node, err := c.ring.Get(key)
	if err != nil {
		return "", err
	}

	return node.Store().Get(key)
}
