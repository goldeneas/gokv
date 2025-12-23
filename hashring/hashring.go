package hashring

import (
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"slices"
	"sync"
)

var (
	ErrNoConnectedNodes = errors.New("no connected nodes available")
	ErrNodeExists       = errors.New("node already exists")
	ErrNodeNotFound     = errors.New("node not found")
	ErrInHashingKey     = errors.New("error in hashing the key")
)

type Node interface {
	Identifier() string
}

type config struct {
	HashFunction func() hash.Hash64
	EnableLogs   bool
}

type ConfigFn func(*config)

func SetHashFunction(f func() hash.Hash64) ConfigFn {
	return func(c *config) {
		c.HashFunction = f
	}
}

func SetEnableLogs(enabled bool) ConfigFn {
	return func(c *config) {
		c.EnableLogs = enabled
	}
}

type HashRing struct {
	mtx               sync.RWMutex
	config            config
	nodes             sync.Map
	sortedKeysOfNodes []uint64
}

func NewHashRing(opts ...ConfigFn) *HashRing {
	config := config{
		HashFunction: fnv.New64a,
		EnableLogs:   false,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return &HashRing{
		config:            config,
		sortedKeysOfNodes: make([]uint64, 0),
	}
}

func (h *HashRing) AddNode(node Node) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	hash, err := h.generateHash(node.Identifier())

	if err != nil {
		return fmt.Errorf("%w: %s", ErrInHashingKey, node.Identifier())
	}

	if _, exists := h.nodes.Load(hash); exists {
		return fmt.Errorf("%w: node %s", ErrNodeExists, node.Identifier())
	}

	h.nodes.Store(hash, node)
	h.sortedKeysOfNodes = append(h.sortedKeysOfNodes, hash)

	slices.Sort(h.sortedKeysOfNodes)

	if h.config.EnableLogs {
		log.Printf("[HashRing] Added Node with Identifier %s (hash: %d)",
			node.Identifier(),
			hash,
		)
	}

	return nil
}

func (h *HashRing) generateHash(key string) (uint64, error) {
	hash := h.config.HashFunction()
	if _, err := hash.Write([]byte(key)); err != nil {
		return 0, err
	}

	return hash.Sum64(), nil
}
