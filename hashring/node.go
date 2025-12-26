package hashring

import "github.com/goldeneas/gokv/store"

type Node interface {
	Identifier() string
	Store() store.Store
}
