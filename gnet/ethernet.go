package gnet

type EthernetHeader struct {
	DstMac    [6]byte
	SrcMac    [6]byte
	EtherType uint16
	//Data      []byte
	//Fcs       [4]byte
}
