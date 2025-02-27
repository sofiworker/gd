package ghttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()
	server.Post("/1", func(ctx *gin.Context) {
		fmt.Println(1)
	})
	server.Post("/2", GinJsonWrap(func() {
		fmt.Println(2)
	}))
	server.Get("/3", GinJsonWrap(func() error {
		fmt.Println(3)
		return nil
	}))
	err := server.ListenAndServe()
	if err != nil {
		t.Fatal(err)
	}
}
