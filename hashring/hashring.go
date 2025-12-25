package hashring

import (
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"log"
	"slices"
	"sort"
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

func (h *HashRing) Add(node Node) error {
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
		log.Printf("added node with key %s (hash: %d)",
			node.Identifier(),
			hash,
		)
	}

	return nil
}

func (h *HashRing) Get(key string) (Node, error) {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	hash, err := h.generateHash(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInHashingKey, key)
	}

	index, err := h.search(hash)
	if err != nil {
		return nil, err
	}

	nodeHash := h.sortedKeysOfNodes[index]
	if node, ok := h.nodes.Load(nodeHash); ok {
		if h.config.EnableLogs {
			log.Printf("Key %s (hash: %d) mapped to node (hash: %d)", key, hash, nodeHash)
		}

		return node.(Node), nil
	}

	return nil, fmt.Errorf("%w: no node found for key %s (hash: %d)", ErrNodeNotFound, key, hash)
}

func (h *HashRing) Remove(key string) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	hash, err := h.generateHash(key)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInHashingKey, key)
	}

	if _, found := h.nodes.LoadAndDelete(hash); !found {
		return fmt.Errorf("%w: %s (hash: %d)", ErrNodeNotFound, key, hash)
	}

	index, err := h.search(hash)
	if err != nil {
		return err
	}

	h.sortedKeysOfNodes = append(h.sortedKeysOfNodes[:index], h.sortedKeysOfNodes[:index+1]...)

	if h.config.EnableLogs {
		log.Printf("removed node with key %s (hash: %d)", key, hash)
	}

	return nil
}

func (h *HashRing) search(hash uint64) (int, error) {
	if len(h.sortedKeysOfNodes) == 0 {
		return -1, ErrNoConnectedNodes
	}

	index := sort.Search(len(h.sortedKeysOfNodes), func(i int) bool {
		return h.sortedKeysOfNodes[i] >= hash
	})

	// wrap around. if no key is found that respects the filter,
	// index where the hash would be inserted
	if index == len(h.sortedKeysOfNodes) {
		index = 0
	}

	return index, nil
}

func (h *HashRing) generateHash(key string) (uint64, error) {
	hash := h.config.HashFunction()
	if _, err := hash.Write([]byte(key)); err != nil {
		return 0, err
	}

	return hash.Sum64(), nil
}
