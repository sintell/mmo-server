package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func readBytesAsBool(data []byte) bool {
	val, bytesRead := binary.Varint(data)
	if bytesRead == 0 {
		panic(fmt.Errorf("no bytes were read"))
	}

	return val != 0
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
type PacketsList map[uint]Packet

// NewPacketsList returns pointer to a new packets list
func (ph *GamePacketHandler) NewPacketsList() *PacketsList {
	return &PacketsList{
		33:   new(fp),
		5008: new(sp),
	}
}

// ReadHead reads packet id from stream
func (ph GamePacketHandler) ReadHead(c io.Reader) (uint, error) {
	if ph.HeadLength == 0 {
		return 0, fmt.Errorf("can't read packet head as it size set to 0")
	}

	data := make([]byte, ph.HeadLength)
	_, err := c.Read(data)
	if err != nil {
		return 0, err
	}

	pid := uint(binary.LittleEndian.Uint16(data))
	return pid, nil
}

// ReadBody convert bytes to structs
func (ph GamePacketHandler) ReadBody(id uint, c io.Reader, p *PacketsList) (Packet, error) {
	size := id - ph.HeadLength
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
		ph.Logger.Debugf("pid: %d\tread: %d\tleft: %d\n", id, read, bytesLeft)
	}

	if packetItem, exists := (*p)[id]; exists {
		err := packetItem.UnmarshalBinary(buf)
		if err != nil {
			return nil, err
		}

		return packetItem, nil
	}

	return nil, fmt.Errorf("Unknown packet id: %d", id)
}
