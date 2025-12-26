# gokv

gokv is a distributed key-value store implementation in Go, utilizing consistent hashing for data partitioning and node distribution.

## Features

- **Consistent Hashing**: efficient distribution of data across nodes with minimal reorganization when nodes are added or removed
- **Thread-Safe**: concurrent access handling using RWMutex
- **TTL Support**: built-in expiration mechanism for key-value pairs
- **Pluggable Storage**: interface-based design allowing custom backend storage implementations
- **Configurable**: customizable hash functions and logging options
