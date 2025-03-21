package gredis

import (
	"log"
	"testing"
)

func TestTx(t *testing.T) {
	client, err := NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}
	errors := client.Tx(func(r Redis) error {
		val, err := r.Get("")
		if err != nil {
			return err
		}
		log.Println(val)
		return nil
	})
	if len(errors) > 0 {
		log.Fatal(errors)
	}
}
