package capture

//
//import (
//	"errors"
//	"fmt"
//	"golang.org/x/sys/unix"
//	"log"
//	"net"
//	"syscall"
//	"unsafe"
//)
//
//type LinuxHandler struct {
//}
//
//const (
//	ETH_P_ALL     = 0x0003    // 捕获所有以太网协议
//	ETH_ALEN      = 6         // MAC地址长度
//	VLAN_HLEN     = 4         // VLAN标签长度（4字节）
//	ETH_FRAME_LEN = 14        // 以太网帧头长度
//	ETH_TYPE_VLAN = 0x8100    // VLAN以太网类型标识
//	BUFFER_SIZE   = 65535     // 抓包缓冲区大小
//	BATCH_SIZE    = 64        // 批量读取数量
//	BLOCK_SIZE    = 4096 * 64 // 内存映射块大小
//	BLOCK_NR      = 64        // 内存映射块数量
//)
//
//type ethhdr struct {
//	DestMAC   [ETH_ALEN]byte
//	SrcMAC    [ETH_ALEN]byte
//	EtherType uint16
//}
//
//type vlanhdr struct {
//	TCI   uint16
//	Proto uint16
//}
//
//type ifreq struct {
//	Name    [16]byte
//	Ifindex int32
//	Flags   uint16
//}
//
//func main() {
//	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, int(htons(ETH_P_ALL)))
//	if err != nil {
//		log.Fatalf("创建套接字失败: %v", err)
//	}
//	defer unix.Close(fd)
//
//	sourceIfaceName := "eth0"
//	sourceIface, err := net.InterfaceByName(sourceIfaceName)
//	if err != nil {
//		log.Fatalf("获取源网卡信息失败: %v", err)
//	}
//
//	addr := unix.SockaddrLinklayer{
//		Protocol: htons(ETH_P_ALL),
//		Ifindex:  sourceIface.Index,
//	}
//	if err := unix.Bind(fd, &addr); err != nil {
//		log.Fatalf("绑定源网卡失败: %v", err)
//	}
//
//	var req unix.TpacketReq = unix.TpacketReq{
//		Block_size: BLOCK_SIZE,
//		Block_nr:   BLOCK_NR,
//		Frame_size: BLOCK_SIZE / BLOCK_NR,
//		Frame_nr:   BLOCK_SIZE / (BLOCK_SIZE / BLOCK_NR),
//	}
//
//	if err := unix.Setsockopt(fd, unix.SOL_PACKET, unix.PACKET_RX_RING, unsafe.Pointer(&req), unsafe.Sizeof(req)); err != nil {
//		log.Fatalf("配置内存映射环形缓冲区失败: %v", err)
//	}
//
//	buffer, err := syscall.Mmap(int(fd), 0, int(req.Block_size*req.Block_nr), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
//	if err != nil {
//		log.Fatalf("内存映射失败: %v", err)
//	}
//	defer syscall.Munmap(buffer)
//
//	if err := setPromiscuous(fd, sourceIfaceName); err != nil {
//		log.Fatalf("设置混杂模式失败: %v", err)
//	}
//
//	targetIfaceName := "eth1"
//	targetIface, err := net.InterfaceByName(targetIfaceName)
//	if err != nil {
//		log.Fatalf("获取目标网卡信息失败: %v", err)
//	}
//
//	targetAddr := unix.SockaddrLinklayer{
//		Protocol: htons(ETH_P_ALL),
//		Ifindex:  targetIface.Index,
//	}
//
//	for {
//		var mmsgs [BATCH_SIZE]unix.Msghdr
//		var iovecs [BATCH_SIZE]unix.Iovec
//		for i := range mmsgs {
//			iovecs[i] = unix.Iovec{
//				Base: &buffer[int(i)*int(req.Frame_size)],
//				Len:  uint64(req.Frame_size),
//			}
//			mmsgs[i] = unix.Msghdr{
//				Iov:    &iovecs[i],
//				Iovlen: 1,
//			}
//		}
//
//		n, err := unix.Recvmsg(fd, buffer, nil, 0)
//		if err != nil {
//			if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
//				continue
//			}
//			log.Printf("批量读取失败: %v", err)
//			continue
//		}
//
//		for i := 0; i < n; i++ {
//			msg := &mmsgs[i]
//			if msg.MsgLen == 0 {
//				continue
//			}
//
//			packet := buffer[int(i)*int(req.Frame_size) : int(i+1)*int(req.Frame_size)]
//			packetLen := int(msg.MsgLen)
//
//			eth := (*ethhdr)(unsafe.Pointer(&packet[0]))
//			fmt.Printf("源MAC: %s -> 目标MAC: %s\n",
//				net.HardwareAddr(eth.SrcMAC[:]), net.HardwareAddr(eth.DestMAC[:]))
//
//			offset := ETH_FRAME_LEN
//			if eth.EtherType == htons(ETH_TYPE_VLAN) {
//				vlan := (*vlanhdr)(unsafe.Pointer(&packet[offset]))
//				vlanID := vlan.TCI & 0x0FFF
//				fmt.Printf("检测到VLAN ID: %d\n", vlanID)
//				offset += VLAN_HLEN
//			}
//
//			if err := unix.Sendto(fd, packet[:packetLen], 0, &targetAddr); err != nil {
//				log.Printf("转发报文失败: %v", err)
//			}
//		}
//	}
//}
//
//func setPromiscuous(fd int, ifaceName string) error {
//	ifr, err := unix.NewIfreq(ifaceName)
//	if err != nil {
//		return err
//	}
//	if err := unix.IoctlIfreq(fd, unix.SIOCGIFFLAGS, ifr); err != nil {
//		return err
//	}
//	flags := ifr.Flags()
//	flags |= unix.IFF_PROMISC
//	ifr.SetFlags(flags)
//	return unix.IoctlIfreq(fd, unix.SIOCSIFFLAGS, ifr)
//}
//
//func htons(n uint16) uint16 {
//	return (n << 8) | (n >> 8)
//}
