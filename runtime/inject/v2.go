package inject

import (
	"context"
	"errors"
	"github.com/chuck1024/gd/v2/reflectx"
	"go.uber.org/dig"
	"reflect"
)

var (
	NotStructError               = errors.New("not a struct")
	BasicTypeShouldWithNameError = errors.New("basic type should with name")
)

type Container struct {
	container *dig.Container
	injectTag string
}

type Lifecycle interface {
	OnStart(context.Context) error
	OnStop(context.Context) error
}

type Instance interface {
	New(params ...interface{}) (interface{}, error)
}

func New() *Container {
	container := dig.New()
	return &Container{container: container, injectTag: "gd"}
}

func (c *Container) SetInjectTag(tag string) {
	c.injectTag = tag
}

func (c *Container) Provide(data interface{}, opts ...dig.ProvideOption) error {
	if reflectx.IsBasicType(data) && len(opts) == 0 {
		return BasicTypeShouldWithNameError
	}

	constructor, err := c.GenerateConstructor(data)
	if err != nil {
		return err
	}
	return c.container.Provide(constructor, opts...)
}

// GenerateConstructor 通过反射生成结构体的构造函数
func (c *Container) GenerateConstructor(instance interface{}) (interface{}, error) {
	// 获取实例的类型
	t := reflect.TypeOf(instance)

	// 处理指针类型（如传入 &MyStruct{}，则构造函数返回 *MyStruct）
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem() // 获取指针指向的类型（即 MyStruct）
	}
	//// 确保传入的是结构体
	//if t.Kind() != reflect.Struct {
	//	return nil, NotStructError
	//}

	valueOf := reflect.ValueOf(instance)
	if reflectx.IsBasicType(instance) {
		ctorType := reflect.FuncOf(nil, []reflect.Type{t}, false)
		fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {
			if isPtr {
				return []reflect.Value{valueOf.Addr()} // 返回指针
			}
			return []reflect.Value{valueOf} // 返回值
		})
		return fn.Interface(), nil
	}

	// 收集结构体字段的类型作为构造函数参数
	var paramTypes []reflect.Type
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tagVal := field.Tag.Get(c.injectTag)
		if tagVal == "-" {
			continue
		}
		if field.Type.Kind() != reflect.Struct || field.Type.Elem().Kind() != reflect.Struct {
			continue
		}
		paramTypes = append(paramTypes, field.Type)
	}

	// 构造函数的返回类型
	returnType := t
	if isPtr {
		returnType = reflect.PointerTo(t)
	}

	// 定义构造函数类型：func(...paramTypes) returnType
	ctorType := reflect.FuncOf(paramTypes, []reflect.Type{returnType}, false)

	// 创建函数实现
	fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {

		// 将参数赋值给结构体的导出字段
		argIndex := 0
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue // 跳过未导出的字段
			}
			tagVal := field.Tag.Get(c.injectTag)
			if tagVal == "-" {
				continue
			}
			if field.Type.Kind() != reflect.Struct || field.Type.Elem().Kind() != reflect.Struct {
				continue
			}
			if argIndex >= len(args) {
				break
			}
			valueOf.Field(i).Set(args[argIndex])
			argIndex++
		}

		if isPtr {
			return []reflect.Value{valueOf.Addr()} // 返回指针
		}
		return []reflect.Value{valueOf} // 返回值
	})

	return fn.Interface(), nil
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
