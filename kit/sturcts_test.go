package kit

import (
	"testing"

	"github.com/emirpasic/gods/lists/arraylist"
)

func TestArrayList(t *testing.T) {
	list := arraylist.New()
	list.Add("a")      // ["a"]
	list.Add("c", "b") // ["a","c","b"]
}
