package inject

import (
	"testing"
)

//type A struct {
//	Age int
//	b   B
//}
//
//func NewA() *A {
//	return &A{}
//}
//
//type B struct {
//	Email string
//	c     *C
//}
//
//func NewB() *B {
//	return &B{}
//}
//
//type C struct {
//	name string
//}
//
//func NewC() *C {
//	return &C{}
//}

func TestNew(t *testing.T) {
	//container := New()
	//a := &A{Age: 1}
	//err := container.Provide(a)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//b := &B{Email: "test"}
	//err = container.Provide(b)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//c := &C{name: "test"}
	//err = container.Provide(c)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//c := dig.New()
	//err := c.Provide(NewA)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = c.Provide(NewB)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = c.Provide(NewC)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//err = c.Invoke(func(a *A) {
	//	fmt.Println(a.Age)
	//	fmt.Println(a.b.Email)
	//	fmt.Println(a.b.c.name)
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
}
