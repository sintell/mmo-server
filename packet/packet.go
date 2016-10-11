package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func readBytesAsBool(data []byte) bool {
	return readBytesAsUint8(data) > 0
}

func readBytesAsInt8(data []byte) int8 {
	return int8(data[0])
}

func readBytesAsUint8(data []byte) uint8 {
	return uint8(data[0])
}

func readBytesAsInt16(data []byte) int16 {
	return int16(binary.LittleEndian.Uint16(data))
}

func readBytesAsUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func readBytesAsInt32(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}

func readBytesAsUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func readBytesAsFloat32(data []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data))
}

func readBytesAsFloat64(data []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(data))
}

// Logger is interface for logging service
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// Packet represent datagramms and provides a convinient set of operations
type Packet interface {
	MarshalBinary() []byte
	UnmarshalBinary([]byte) error
}

// GamePacketHandler represent a set of funcitons to handle game client packets
type GamePacketHandler struct {
	HeadLength uint
	Logger     Logger
	PacketList PacketsList
}

// PacketsList holds all awailable packets
type PacketsList map[uint16]Packet

// NewPacketsList returns pointer to a new packets list
func (ph *GamePacketHandler) NewPacketsList() *PacketsList {
	return &PacketsList{
		5189: new(fp),
		5190: new(sp),
	}
}

// ReadHead reads packet id from stream
func (ph GamePacketHandler) ReadHead(c io.Reader) (*HeaderPacket, error) {
	if ph.HeadLength == 0 {
		return nil, fmt.Errorf("can't read packet head as it size set to 0")
	}

	buf := make([]byte, 6)
	hp := new(HeaderPacket)
	_, err := c.Read(buf)

	hp.UnmarshalBinary(buf)
	if err != nil {
		return nil, err
	}

	return hp, nil
}

// ReadBody convert bytes to structs
func (ph GamePacketHandler) ReadBody(header *HeaderPacket, c io.Reader, p *PacketsList) (Packet, error) {
	size := uint(header.length - 6)
	buf := make([]byte, 0, size)
	t := make([]byte, size/2)
	bytesLeft := size
	for bytesLeft > 0 {
		read, err := c.Read(t[:int(math.Min(float64(bytesLeft), float64(size/2)))])
		if err != nil {
			return nil, err
		}

		buf = append(buf, t[:read]...)
		bytesLeft = bytesLeft - uint(read)
		ph.Logger.Debugf("pid: %d\tread: %d\tleft: %d\n", header.ID, read, bytesLeft)
	}

	if packetItem, exists := (*p)[header.ID]; exists {
		err := packetItem.UnmarshalBinary(buf)
		if err != nil {
			return nil, err
		}

		return packetItem, nil
	}

	return nil, fmt.Errorf("Unknown packet id: %d", header.ID)
}
