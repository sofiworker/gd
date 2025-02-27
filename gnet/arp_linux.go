package gnet

import (
	"encoding/binary"
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

// 以太网帧头结构
type EthernetHeader struct {
	DstMAC    [6]byte
	SrcMAC    [6]byte
	EtherType uint16
}

// ARP报文结构
type ArpHeader struct {
	HardwareType    uint16
	ProtocolType    uint16
	HardwareAddrLen uint8
	ProtocolAddrLen uint8
	Operation       uint16
	SenderMAC       [6]byte
	SenderIP        [4]byte
	TargetMAC       [6]byte
	TargetIP        [4]byte
}

func ListenArp() {
	// 创建原始套接字（需root权限）
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		log.Fatal("Socket error:", err)
	}
	defer syscall.Close(fd)

	// 缓冲区接收数据
	buffer := make([]byte, 1500)
	for {
		n, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			continue
		}
		if n < 14+28 { // 以太网头+ARP报文最小长度
			continue
		}

		// 解析以太网帧头
		eth := (*EthernetHeader)(unsafe.Pointer(&buffer[0]))
		ethType := ntohs(eth.EtherType)
		if ethType != 0x0806 { // 过滤非ARP包
			continue
		}

		// 解析ARP报文
		arpData := buffer[14:]
		arp := (*ArpHeader)(unsafe.Pointer(&arpData[0]))

		// 打印信息
		fmt.Printf("ARP Operation: %d (1=Request, 2=Reply)\n", ntohs(arp.Operation))
		fmt.Printf("Sender MAC: %02x:%02x:%02x:%02x:%02x:%02x, IP: %d.%d.%d.%d\n",
			arp.SenderMAC[0], arp.SenderMAC[1], arp.SenderMAC[2],
			arp.SenderMAC[3], arp.SenderMAC[4], arp.SenderMAC[5],
			arp.SenderIP[0], arp.SenderIP[1], arp.SenderIP[2], arp.SenderIP[3])
		fmt.Printf("Target MAC: %02x:%02x:%02x:%02x:%02x:%02x, IP: %d.%d.%d.%d\n",
			arp.TargetMAC[0], arp.TargetMAC[1], arp.TargetMAC[2],
			arp.TargetMAC[3], arp.TargetMAC[4], arp.TargetMAC[5],
			arp.TargetIP[0], arp.TargetIP[1], arp.TargetIP[2], arp.TargetIP[3])
		fmt.Println("----------------------------------------------------")
	}
}

func htons(x uint16) uint16 {
	return (x << 8) | (x >> 8)
}

func ntohs(x uint16) uint16 {
	return binary.BigEndian.Uint16((*[2]byte)(unsafe.Pointer(&x))[:])
}
