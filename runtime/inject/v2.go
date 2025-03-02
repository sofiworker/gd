package inject

import (
	"bytes"
	"context"
	"errors"
	"github.com/chuck1024/gd/v2/reflectx"
	"go.uber.org/dig"
	"io"
	"os"
	"reflect"
	"sync"
)

var (
	NotStructError               = errors.New("not a struct")
	BasicTypeShouldWithNameError = errors.New("basic type should with name")
	NameOrGroupOnlyOneError      = errors.New("name or group only one")
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

func New() *Container {
	container := dig.New(dig.RecoverFromPanics())
	return &Container{container: container, injectTag: "gd", locker: &sync.RWMutex{}}
}

func (c *Container) SetInjectTag(tag string) *Container {
	c.injectTag = tag
	return c
}

type NameOption string

func (o NameOption) ApplyProvideOption(opt *ProvideOptions) {
	opt.Name = string(o)
}

func WithName(name string) ProvideOption {
	return NameOption(name)
}

type GroupOption string

func (o GroupOption) ApplyProvideOption(opt *ProvideOptions) {
	opt.Group = string(o)
}

func WithGroup(group string) ProvideOption {
	return GroupOption(group)
}

func (c *Container) ProvideWithName(data interface{}, name string) error {
	return c.Provide(data, WithName(name))
}

func (c *Container) Provide(data interface{}, opts ...ProvideOption) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	hasName := false
	hasGroup := false
	for _, opt := range opts {
		if _, ok := opt.(NameOption); ok {
			hasName = true
		}
		if _, ok := opt.(GroupOption); ok {
			hasGroup = true
		}
	}

	if reflectx.IsBasicType(data) && (!hasName || !hasGroup) {
		return BasicTypeShouldWithNameError
	}

	if hasName && hasGroup {
		return NameOrGroupOnlyOneError
	}

	options := BuildDigProvideOption(opts...)

	constructor, err := c.GenerateConstructor(data)
	if err != nil {
		return err
	}
	return c.container.Provide(constructor, options...)
}

func BuildDigProvideOption(opts ...ProvideOption) []dig.ProvideOption {
	var provideOpts ProvideOptions
	for _, opt := range opts {
		opt.ApplyProvideOption(&provideOpts)
	}
	ret := make([]dig.ProvideOption, 0)
	if provideOpts.Name != "" {
		ret = append(ret, dig.Name(provideOpts.Name))
	}
	if provideOpts.Group != "" {
		ret = append(ret, dig.Group(provideOpts.Group))
	}
	if provideOpts.As != nil {
		for _, as := range provideOpts.As {
			ret = append(ret, dig.As(as))
		}
	}
	if provideOpts.Group != "" {
		ret = append(ret, dig.Group(provideOpts.Group))
	}

	return ret
}

// GenerateConstructor generates a constructor function for the given instance using reflection
func (c *Container) GenerateConstructor(instance interface{}) (interface{}, error) {
	if instance == nil {
		return nil, errors.New("instance cannot be nil")
	}

	// Get instance type
	t := reflect.TypeOf(instance)
	if t == nil {
		return nil, errors.New("cannot get type information")
	}

	// Handle pointer types
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem()
		if t == nil {
			return nil, errors.New("nil pointer dereference")
		}
	}

	valueOf := reflect.ValueOf(instance)
	if isPtr {
		valueOf = valueOf.Elem()
	}

	if reflectx.IsBasicType(instance) {
		ctorType := reflect.FuncOf(nil, []reflect.Type{t}, false)
		fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {
			if isPtr {
				return []reflect.Value{valueOf.Addr()}
			}
			return []reflect.Value{valueOf}
		})
		return fn.Interface(), nil
	}

	// Collect struct field types as constructor parameters
	var paramTypes []reflect.Type
	var fieldIndices []int
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if tagVal := field.Tag.Get(c.injectTag); tagVal == SkipInjectTag {
			continue
		}
		if !isBasicType(field.Type.Kind()) {
			paramTypes = append(paramTypes, field.Type)
			fieldIndices = append(fieldIndices, i)
		}
	}

	// Set return type
	returnType := t
	if isPtr {
		returnType = reflect.PointerTo(t)
	}

	ctorType := reflect.FuncOf(paramTypes, []reflect.Type{returnType}, false)

	fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {
		newInstance := reflect.New(t).Elem()

		// Copy all fields from initial value
		if valueOf.IsValid() {
			newInstance.Set(valueOf)
		}

		// Set only injectable fields
		for i, fieldIdx := range fieldIndices {
			if i >= len(args) {
				break
			}
			newInstance.Field(fieldIdx).Set(args[i])
		}

		if isPtr {
			return []reflect.Value{newInstance.Addr()}
		}
		return []reflect.Value{newInstance}
	})

	return fn.Interface(), nil
}

func isBasicType(k reflect.Kind) bool {
	switch k {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		return true
	default:
		return false
	}
}

func (c *Container) Invoke(f interface{}) error {
	return c.container.Invoke(f)
}

//func (c *Container) InvokeWithName(f interface{}) error {
//	return c.container.Invoke(f)
//}

func (c *Container) PrintGraph(writers ...io.Writer) {
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}
	for _, w := range writers {
		err := dig.Visualize(c.container, w)
		if err != nil {
			panic(err)
		}
		return
	}
}

func (c *Container) GetGraphString() string {
	var b bytes.Buffer
	err := dig.Visualize(c.container, &b)
	if err != nil {
		panic(err)
	}
	return b.String()
}
