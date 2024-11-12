package logic

import (
	"protchain/internal/dep"
	"sync"
)

var (
	logic = new(Logic)
	once  sync.Once
)

type Logic struct{}

func New(dep *dep.Dependencies) *Logic {
	once.Do(func() {
		logic = &Logic{}
	})

	return logic
}
