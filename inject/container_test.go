package inject

import (
	"fmt"
	"testing"
)

func TestOptions(t *testing.T) {
	container := New(WithInjectTag("test"))
	fmt.Printf("%+v", *container.opts)
	err := container.Provide(1)
	if err != nil {
		t.Fatal(err)
	}
}
