package inject

import (
	"fmt"
)

var (
	ErrModuleNotFound      = fmt.Errorf("module not found")
	ErrModuleNameEmpty     = fmt.Errorf("module name is empty")
	ErrModuleAlreadyExists = fmt.Errorf("module already exists")
)

type ContainerOptionFunc func(*ContainerOption)

type ContainerOption struct {
	//injectTag string
	name string
}

type Container struct {
	//locker            sync.RWMutex
	opts *ContainerOption
	root *Module
	//recoverFromPanics bool
	//startTimeout      time.Duration
	//stopTimeout       time.Duration
}

//func WithInjectTag(tag string) ContainerOptionFunc {
//	return func(opts *ContainerOption) {
//		opts.injectTag = tag
//	}
//}

func WithName(name string) ContainerOptionFunc {
	return func(opts *ContainerOption) {
		opts.name = name
	}
}

func New(opts ...ContainerOptionFunc) *Container {
	c := &Container{root: &Module{Name: "root", Parent: nil}}

	o := &ContainerOption{
		//injectTag: "gd",
	}
	for _, opt := range opts {
		opt(o)
	}

	c.opts = o

	return c
}

func (c *Container) NewModule(name string, opts ...ModuleOptionFunc) *Module {
	return &Module{Name: name, Parent: c.root}
}

func (c *Container) Root() *Module {
	return c.root
}
