package inject

import (
	"fmt"
	"go.uber.org/dig"
	"testing"
)

type D struct {
	Email string `gd:"-"`
}

// 定义结构体 C
type C struct {
	Name string
	D    *D `gd:"-"`
}

// 定义结构体 B 并包含 C
type B struct {
	C *C
}

// B 的构造函数，需要 C 作为参数
func NewB(c *C) *B {
	return &B{C: c}
}

// 定义结构体 A 并包含 B
type A struct {
	B *B
}

// A 的构造函数，需要 B 作为参数
func NewA(b *B) *A {
	return &A{B: b}
}

func TestDig(t *testing.T) {
	container := dig.New()

	// 提供 C 的实例
	container.Provide(func() *C {
		return &C{Name: "Hello from C"}
	})

	// 提供 B 的实例，依赖于 C
	container.Provide(NewB)

	// 提供 A 的实例，依赖于 B
	container.Provide(NewA)

	// 解析出 A 的实例，同时也会解析出其依赖的 B 和 C
	var a *A
	err := container.Invoke(func(aInstance *A) {
		a = aInstance
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("A contains B, and B contains C with name: %s\n", a.B.C.Name)
}

func TestContainer(t *testing.T) {
	container := New()
	a := &A{}
	b := &B{}
	c := &C{
		Name: "1111111111111",
	}
	d := &D{}
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

	err = container.Invoke(func(a *A) {
		fmt.Println("==============")
		fmt.Printf("A contains B, and B contains C with name: %s\n", a.B.C.Name)
		fmt.Println("==============")
	})
	if err != nil {
		t.Error(err)
	}
}
