package pcapx

type PcapFile struct {
	Header  *PcapHeader
	Packets []*Packet
}

type PcapHeader struct {
	MagicNumber  uint32
	VersionMajor uint16
	VersionMinor uint16
	ThisZone     int32
	SigFigs      uint32
	SnapLen      uint32
	Network      uint32
}

type PacketHeader struct {
	TimestampSec  uint32
	TimestampUSec uint32
	CapturedLen   uint32
	OrigLen       uint32
}

// Packet is a complete packet including its header and data.
type Packet struct {
	Header *PacketHeader
	Data   []byte
}

type EthernetData struct {
}

type IpData struct {
}

type TcpData struct{}

type UdpData struct{}
