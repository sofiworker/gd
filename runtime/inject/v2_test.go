package inject

import (
	"fmt"
	"go.uber.org/dig"
	"reflect"
	"testing"
)

// ProvideStruct 将结构体包装为构造函数，并通过 dig.Provide 注入
func ProvideStruct(container *dig.Container, instance interface{}) error {
	// 获取实例的类型
	t := reflect.TypeOf(instance)

	// 处理指针类型（如传入 &MyStruct{}，则构造函数返回 *MyStruct）
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem() // 获取指针指向的类型（即 MyStruct）
	}

	// 确保传入的是结构体
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct type, got %v", t.Kind())
	}

	// 生成构造函数类型：func(...deps) <T or *T>
	ctor := generateConstructor(t, isPtr, instance)

	// 注入生成的构造函数
	return container.Provide(ctor)
}

// generateConstructor 通过反射生成结构体的构造函数
func generateConstructor(t reflect.Type, isPtr bool, data interface{}) interface{} {
	// 收集结构体字段的类型作为构造函数参数
	var paramTypes []reflect.Type
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // 跳过未导出的字段
		}
		paramTypes = append(paramTypes, field.Type)
	}

	// 构造函数的返回类型
	returnType := t
	if isPtr {
		returnType = reflect.PtrTo(t)
	}

	// 定义构造函数类型：func(...paramTypes) returnType
	ctorType := reflect.FuncOf(paramTypes, []reflect.Type{returnType}, false)

	// 创建函数实现
	fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {
		//instance := reflect.New(t).Elem() // 创建结构体实例
		instance := reflect.ValueOf(data).Elem()

		// 将参数赋值给结构体的导出字段
		argIndex := 0
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue // 跳过未导出的字段
			}
			if argIndex >= len(args) {
				break
			}
			instance.Field(i).Set(args[argIndex])
			argIndex++
		}

		if isPtr {
			return []reflect.Value{instance.Addr()} // 返回指针
		}
		return []reflect.Value{instance} // 返回值
	})

	return fn.Interface()
}

type Logger struct {
	Name string
}

func NewLogger() *Logger {
	return &Logger{
		Name: "logger",
	}
}

type MyStruct struct {
	Logger  Logger // 假设 Logger 是已注册的依赖
	Timeout int    `gd:"-"` // 依赖参数
}

func TestNew(t *testing.T) {
	container := dig.New()

	err := container.Provide(NewLogger)
	if err != nil {
		t.Error(err)
	}

	// 直接传入结构体实例（指针）
	err = ProvideStruct(container, &MyStruct{})
	if err != nil {
		t.Error(err)
	}

	// 从容器中获取实例
	err = container.Invoke(func(s *MyStruct) {
		fmt.Println("MyStruct injected:", s.Logger.Name)
	})
	if err != nil {
		t.Error(err)
	}
}

func TestBasic(t *testing.T) {
	container := New()
	err := container.Provide("aaa", dig.Name("a"))
	if err != nil {
		t.Error(err)
	}
	err = container.Provide("bbbbb", dig.Name("b"))
	if err != nil {
		t.Error(err)
	}
}
