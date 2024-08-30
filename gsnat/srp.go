package gsnat

import (
	"encoding/binary"
	"fmt"

	"github.com/lekuruu/ubisoft-game-service/common"
)

const PACKET_BUFFER_SIZE = 1024
const SRP_HEADER_SIZE = 12
const SRP_WINDOW_SIZE = 8

const (
	FLAGS_FIN    = 1
	FLAGS_SYN    = 2
	FLAGS_ACK    = 4
	FLAGS_URG    = 8
	FLAGS_SRP_ID = 0x3040
)

type SRPWindow struct {
	Tail              uint16
	SenderSignature   uint16
	ChecksumInitValue uint16
	WndBufferSize     uint16
}

func (window *SRPWindow) String() string {
	return fmt.Sprintf(
		"SRPWindow{Tail: %d, SenderSignature: %d, ChecksumInitValue: %d, WndBufferSize: %d}",
		window.Tail, window.SenderSignature, window.ChecksumInitValue, window.WndBufferSize,
	)
}

func (window *SRPWindow) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, common.WriteU16(window.Tail)...)
	data = append(data, common.WriteU16(window.SenderSignature)...)
	data = append(data, common.WriteU16(window.ChecksumInitValue)...)
	data = append(data, common.WriteU16(window.WndBufferSize)...)
	return data
}

type SRPPacket struct {
	Checksum  uint16
	Signature uint16
	DataSize  uint16
	Flags     uint16
	Seg       uint16
	Ack       uint16
	Window    SRPWindow
}

func (srp *SRPPacket) String() string {
	return fmt.Sprintf(
		"SRPPacket{Checksum: %d, Signature: %d, DataSize: %d, Flags: %d, Seg: %d, Ack: %d, Window: %s}",
		srp.Checksum, srp.Signature, srp.DataSize, srp.Flags, srp.Seg, srp.Ack, srp.Window.String(),
	)
}

func (srp *SRPPacket) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, common.WriteU16(srp.Checksum)...)
	data = append(data, common.WriteU16(srp.Signature)...)
	data = append(data, common.WriteU16(srp.DataSize)...)
	data = append(data, common.WriteU16(srp.Flags)...)
	data = append(data, common.WriteU16(srp.Seg)...)
	data = append(data, common.WriteU16(srp.Ack)...)
	data = append(data, srp.Window.Serialize()...)
	return data
}

func ReadSRPPacket(client *Client) (*SRPPacket, error) {
	header := make([]byte, SRP_HEADER_SIZE)
	_, err := client.Reader.Read(header)

	if err != nil {
		return nil, err
	}

	packet := &SRPPacket{
		Checksum:  common.ReadU16(header[0:2]),
		Signature: common.ReadU16(header[2:4]),
		DataSize:  common.ReadU16(header[4:6]),
		Flags:     common.ReadU16(header[6:8]),
		Seg:       common.ReadU16(header[8:10]),
		Ack:       common.ReadU16(header[10:12]),
	}

	if packet.Checksum <= 0 || packet.DataSize <= 0 {
		// Empty data
		return nil, nil
	}

	windowHeader := make([]byte, SRP_WINDOW_SIZE)
	_, err = client.Reader.Read(windowHeader)

	if err != nil {
		return nil, err
	}

	packet.Window = SRPWindow{
		Tail:              common.ReadU16(windowHeader[0:2]),
		SenderSignature:   common.ReadU16(windowHeader[2:4]),
		ChecksumInitValue: common.ReadU16(windowHeader[4:6]),
		WndBufferSize:     common.ReadU16(windowHeader[6:8]),
	}

	return packet, nil
}

func NewSRPPacketFromRequest(request *SRPPacket) SRPPacket {
	packet := SRPPacket{
		Checksum:  request.Window.ChecksumInitValue,
		Signature: request.Window.SenderSignature,
		DataSize:  SRP_WINDOW_SIZE,
		Flags:     FLAGS_SRP_ID | FLAGS_SYN | FLAGS_ACK,
		Seg:       request.Seg + 1,
		Ack:       request.Seg,
	}

	packet.Window = SRPWindow{
		Tail:              10,
		SenderSignature:   2,
		ChecksumInitValue: 0,
		WndBufferSize:     536,
	}

	packet.Checksum = makeChecksum(packet.Serialize())
	return packet
}

func makeChecksum(data []byte) uint16 {
	var truncPos int
	var checkBase uint32
	halfLen := len(data) >> 1
	oddLen := len(data)%2 == 1

	if oddLen {
		// Add the first byte as extra
		checkBase += uint32(data[0])
		truncPos++
	}
	if halfLen > 0 {
		for i := 0; i < halfLen; i++ {
			checkBase += uint32(binary.LittleEndian.Uint16(data[truncPos:]))
			truncPos += 2
		}
	}

	checksum := checkBase & 0xFFFF
	checksum += checkBase >> 16
	checksum += checksum >> 16
	checksum = ^checksum & 0xFFFF
	return uint16(checksum)
}
