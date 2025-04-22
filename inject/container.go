package inject

import (
	"fmt"
	"sync"
	"time"
)

var (
	ErrModuleNotFound      = fmt.Errorf("module not found")
	ErrModuleNameEmpty     = fmt.Errorf("module name is empty")
	ErrModuleAlreadyExists = fmt.Errorf("module already exists")
)

type ContainerOptionFunc func(*ContainerOption)

type ContainerOption struct {
	injectTag string
	name      string
}

type Container struct {
	locker            sync.RWMutex
	opts              *ContainerOption
	root              *Module
	recoverFromPanics bool
	startTimeout      time.Duration
	stopTimeout       time.Duration
}

func WithInjectTag(tag string) ContainerOptionFunc {
	return func(opts *ContainerOption) {
		opts.injectTag = tag
	}
}

func WithName(name string) ContainerOptionFunc {
	return func(opts *ContainerOption) {
		opts.name = name
	}
}

func New(opts ...ContainerOptionFunc) *Container {
	c := &Container{root: &Module{Name: "root", Parent: nil}}

	o := &ContainerOption{
		injectTag: "gd",
	}
	for _, opt := range opts {
		opt(o)
	}

	c.opts = o

	return c
}

func (c *Container) ProvideModule(module *Module) error {
	return c.ProvideModuleWithParentName("", module)
}

func (c *Container) ProvideModuleWithParentName(parentName string, module *Module) error {
	if module == nil || module.Name == "" {
		return ErrModuleNameEmpty
	}

	c.locker.Lock()
	defer c.locker.Unlock()

	checkModule := c.findModuleByNameLocked(module.Name, c.root)
	if checkModule != nil {
		return ErrModuleAlreadyExists
	}
	parent := c.findModuleByNameLocked(parentName, c.root)
	if parent == nil {
		parent = c.root
	}
	parent.modules = append(parent.modules, module)
	return nil
}

func (c *Container) FindModuleByName(name string, current *Module) *Module {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return c.findModuleByNameLocked(name, current)
}

func (c *Container) findModuleByNameLocked(name string, current *Module) *Module {
	if current.Name == name {
		return current
	}
	for _, subModule := range current.modules {
		if found := c.findModuleByNameLocked(name, subModule); found != nil {
			return found
		}
	}
	return nil
}
