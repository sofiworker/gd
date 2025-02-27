package gnet

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
)

func ListenICMP(network, address string) (net.PacketConn, error) {
	if network == "" {
		network = "ip:icmp"
	}
	if address == "" {
		address = "0.0.0.0"
	}
	conn, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ParseICMPMessage(bs []byte) (*icmp.Message, error) {
	return icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), bs)
}
