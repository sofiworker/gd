package pcapx

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

type Handle struct {
	file       string
	reader     io.Reader
	fileHeader *PcapHeader
}

func (h *Handle) SetFilter() {}

func (h *Handle) Range(f func(p *Packet) error) {
	bufReader := bufio.NewReader(h.reader)
	for {
		packetHeader := new(PacketHeader)
		err := binary.Read(bufReader, binary.LittleEndian, packetHeader)
		if err != nil {
			if err != io.EOF {
				break
			}
		}
		data := make([]byte, packetHeader.CapturedLen)
		n, err := bufReader.Read(data)
		if err != nil {
			if err != io.EOF {
				break
			}
		}
		realData := data[0:n]
		packet := &Packet{
			Header: packetHeader,
			Data:   realData,
		}
		err = f(packet)
		if err != nil {

		}
	}
}

func (h *Handle) Packets() chan *Packet {
	return nil
}

func OpenFile(file string) (handle *Handle, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var fileHeader PcapHeader
	err = binary.Read(f, binary.LittleEndian, &fileHeader)
	if err != nil {
		return nil, err
	}
	return &Handle{file: file, reader: f, fileHeader: &fileHeader}, nil
}

func MergePcapFile(files ...string) error {
	return nil
}
