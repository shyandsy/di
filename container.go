package di

import (
	"sync"
)

type container struct {
	mu             *sync.Mutex
	singletonStore sync.Map
	typeCache      sync.Map
}

func NewContainer() Container {
	return &container{
		mu:             &sync.Mutex{},
		singletonStore: sync.Map{},
		typeCache:      sync.Map{},
	}
}
