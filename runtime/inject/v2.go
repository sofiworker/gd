package inject

import (
	"context"
	"go.uber.org/dig"
	"reflect"
)

type Container struct {
	container *dig.Container
}

type Lifecycle interface {
	OnStart(context.Context) error
	OnStop(context.Context) error
}

type Instance interface {
	New() (interface{}, error)
}

func New() *Container {
	container := dig.New()
	return &Container{container: container}
}

func (c *Container) Provide(data interface{}) error {
	typeOf := reflect.TypeOf(data)
	iface := reflect.TypeOf((*Instance)(nil)).Elem()
	if typeOf.Implements(iface) {
		return c.container.Provide(data.(Instance).New)
	}

	// 构建构造函数类型 func() T
	ctorType := reflect.FuncOf(
		nil,                    // 无输入参数
		[]reflect.Type{typeOf}, // 返回类型
		false,                  // 非可变参数
	)

	// 创建构造函数实现
	ctorValue := reflect.MakeFunc(ctorType, func(_ []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(data)}
	})

	return c.container.Provide(ctorValue.Interface())
}

func (c *Container) Register() error {
	err := c.container.Provide(func() {

	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) RegisterSingle() error {
	err := c.container.Provide(func() {

	}, dig.Name("ro"))
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) Invoke(f interface{}) error {
	return c.container.Invoke(f)
}
