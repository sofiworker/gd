package pcapx

import (
	"testing"
)

func TestHandle(t *testing.T) {
	handle, err := OpenFile("test.pcap")
	if err != nil {
		t.Fatal(err)
	}
	handle.Range(func(p *Packet) error {
		//fmt.Printf("%+v \n", *p)
		return nil
	})
}
