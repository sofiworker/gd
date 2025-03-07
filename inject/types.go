package inject

import (
	"context"
	"errors"
	"go.uber.org/dig"
	"sync"
)

var (
	TypeError                    = errors.New("not a struct or basic type")
	BasicTypeShouldWithNameError = errors.New("basic type should with name")
	NameOrGroupOnlyOneError      = errors.New("name or group only one")
	NilError                     = errors.New("instance cannot be nil")
	NilPtrError                  = errors.New("nil pointer dereference")
	SkipInjectTag                = "-"
)

type Container struct {
	container *dig.Container
	injectTag string
	locker    *sync.RWMutex
}

type AppLifecycle interface {
	OnStart(context.Context) error
	OnStop(context.Context) error
}

type InvokeLifecycle interface {
	BeforeInject(context.Context) error
	AfterInject(context.Context) error
}

type ProvideOption interface {
	ApplyProvideOption(*ProvideOptions)
}

type ProvideOptions struct {
	Name  string
	Group string
	As    []interface{}
}
