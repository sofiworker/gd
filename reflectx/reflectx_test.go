package reflectx

import (
	"fmt"
	"log"
	"testing"
)

type Test struct {
	Name string
	Age  int
	// StopCh  chan struct{}
	Hobbies []string
}

func (t *Test) String() string {
	return fmt.Sprintf("%+v", *t)
}

func (t Test) Full(name string) error {
	fmt.Println(name)
	return nil
}

func (t *Test) ReturnTest() (tt *Test, err error) {
	tt = t
	return
}

func (t *Test) print() {
	fmt.Print("hello!!\n")
}

func TestNew(t *testing.T) {
	// ints := 120
	// r := New(ints)
	// log.Println(r)

	// floats := 10.0
	// r = New(floats)
	// log.Println(r)

	// ch := make(chan struct{})
	// r = New(ch)
	// log.Println(r)

	// m := make(map[string]struct{})
	// r = New(m)
	// log.Println(r)

	// slices := make([]string, 10)
	// r = New(slices)
	// log.Println(r)

	// arr := [5]int{1, 2, 3, 4, 5}
	// r = New(arr)
	// log.Println(r)

	tt := &Test{
		Name: "test",
		Age:  10,
		// StopCh:  make(chan struct{}),
		Hobbies: []string{"play", "sleep"},
	}
	r := New(tt)
	fmt.Println(r)

	// myFunc := func() {
	// 	fmt.Println("hello world!!!")
	// }
	// r = New(myFunc)
	// log.Println(r)

	// myFunc2 := func() int {
	// 	fmt.Println("hello world!!!")
	// 	return 0
	// }
	// r = New(myFunc2)
	// log.Println(r)
}

func TestCheckPrtDepth(t *testing.T) {
	var i int
	var j *int
	var k **int
	var l ***int
	var m map[string]struct{}

	ret := CheckPtrDepth(i)
	log.Println(ret)
	ret = CheckPtrDepth(j)
	log.Println(ret)
	ret = CheckPtrDepth(k)
	log.Println(ret)
	ret = CheckPtrDepth(l)
	log.Println(ret)
	ret = CheckPtrDepth(m)
	log.Println(ret)

}
