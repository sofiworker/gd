package inject

import (
	"go.uber.org/dig"
	"sync"
)

func New() *Container {
	container := dig.New(dig.RecoverFromPanics())
	return &Container{container: container, injectTag: "gd", locker: &sync.RWMutex{}}
}

func (c *Container) SetInjectTag(tag string) *Container {
	c.injectTag = tag
	return c
}
