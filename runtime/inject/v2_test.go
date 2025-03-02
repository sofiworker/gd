package inject

import (
	"fmt"
	"go.uber.org/dig"
	"reflect"
	"testing"
)

type Logger struct {
	Name string
}

type MyStruct struct {
	Logger  Logger
	Timeout int `gd:"-"`
}

type D struct {
	Email string
}

type C struct {
	Name string
	//D    *D `gd:"-"`
	D *D
}

type B struct {
	C *C
}

type A struct {
	B *B
}

func TestContainer(t *testing.T) {
	container := New()
	a := &A{}
	b := &B{}
	c := &C{
		Name: "1111111111111",
	}
	d := &D{
		Email: "2222222222222",
	}
	err := container.Provide(a)
	if err != nil {
		t.Error(err)
	}
	err = container.Provide(b)
	if err != nil {
		t.Error(err)
	}

	err = container.Provide(c)
	if err != nil {
		t.Error(err)
	}

	err = container.Provide(d)
	if err != nil {
		t.Error(err)
	}

	var a1, a2 *A

	err = container.Invoke(func(a *A) {
		fmt.Println("==============")
		fmt.Printf("A contains B, and B contains C with name: %s\n", a.B.C.Name)
		fmt.Printf("D email: %s\n", a.B.C.D.Email)
		fmt.Println("==============")
		a1 = a
	})
	if err != nil {
		t.Error(err)
	}

	err = container.Invoke(func(a *A) {
		a2 = a
	})
	if err != nil {
		t.Error(err)
	}

	if a1 == a2 {
		fmt.Println("xxxxx")
	}
	if reflect.DeepEqual(a1, a2) {
		fmt.Println("yyyyy")
	}

	container.PrintGraph()
}

func TestBasicInject(t *testing.T) {
	container := New()
	name := "test11111111"
	err := container.ProvideWithName(name, "name")
	if err != nil {
		t.Error(err)
	}

	type MyStruct struct {
		dig.In
		Name string `name:"name"`
	}
	err = container.Invoke(func(b MyStruct) {
		fmt.Println(b.Name)
	})
	if err != nil {
		t.Error(err)
	}
}

func TestMix(t *testing.T) {

}

//// 包装函数：注册基本类型到容器
//func ProvideNamed(c *dig.Container, name string, value interface{}) error {
//	t := reflect.TypeOf(value)
//	of := reflect.ValueOf(value)
//	ctorType := reflect.FuncOf(nil, []reflect.Type{t}, false)
//	fn := reflect.MakeFunc(ctorType, func(args []reflect.Value) []reflect.Value {
//		//if isPtr {
//		//	return []reflect.Value{valueOf.Addr()}
//		//}
//		return []reflect.Value{of}
//	})
//	return c.Provide(fn.Interface(), dig.Name(name))
//}
//
//func TestDig(t *testing.T) {
//	c := dig.New()
//
//	// 注册基本类型
//	err := ProvideNamed(c, "port", 8080)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = ProvideNamed(c, "timeout", 30*time.Second)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = ProvideNamed(c, "env", "production")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// 注入示例
//	err = c.Invoke(func(params struct {
//		dig.In
//		Port    int           `name:"port"`
//		Timeout time.Duration `name:"timeout"`
//		Env     string        `name:"env"`
//	}) {
//		fmt.Printf("Port: %d, Timeout: %s, Env: %s\n", params.Port, params.Timeout, params.Env)
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//}
