package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

// AfterLoginPackets TODO
var AfterLoginPackets []byte

// CharListPacket TODO
var CharListPacket []byte

// StrangePacket TODO
var StrangePacket []byte

func init() {
	f, err := os.OpenFile("./packets/login.pac", os.O_RDONLY, 0555)
	if err != nil {
		return
	}
	fInfo, err := f.Stat()
	if err != nil {
		return
	}
	AfterLoginPackets = make([]byte, fInfo.Size())
	err = binary.Read(f, binary.LittleEndian, AfterLoginPackets)
	if err != nil {
		return
	}

	f, err = os.OpenFile("./packets/char.pac", os.O_RDONLY, 0555)
	if err != nil {
		return
	}
	fInfo, err = f.Stat()
	if err != nil {
		return
	}
	CharListPacket = make([]byte, fInfo.Size())
	err = binary.Read(f, binary.LittleEndian, CharListPacket)
	if err != nil {
		return
	}

	f, err = os.OpenFile("./packets/idk.pac", os.O_RDONLY, 0555)
	if err != nil {
		return
	}
	fInfo, err = f.Stat()
	if err != nil {
		return
	}
	StrangePacket = make([]byte, fInfo.Size())
	err = binary.Read(f, binary.LittleEndian, StrangePacket)
	if err != nil {
		return
	}
}

// Logger is interface for logging service
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// Packet represent datagramms and provides a convinient set of operations
type Packet interface {
	setHeader(*HeaderPacket)
	Header() *HeaderPacket
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
		1111: new(ClientLoginRequestPacket),
		5100: new(ClientLoginInfoPacket),
	}
}

// ReadHead reads packet id from stream
func (ph GamePacketHandler) ReadHead(c io.Reader) (*HeaderPacket, error) {
	if ph.HeadLength == 0 {
		return nil, fmt.Errorf("can't read packet head as it size set to 0")
	}

	buf := make([]byte, ph.HeadLength)
	_, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	ph.Logger.Debugf("head bytes: %+v", buf)

	hp := new(HeaderPacket)
	hp.UnmarshalBinary(buf)

	return hp, nil
}

// ReadBody convert bytes to structs
func (ph GamePacketHandler) ReadBody(header *HeaderPacket, c io.Reader, p *PacketsList) (Packet, error) {
	size := uint(header.Length - uint16(ph.HeadLength))
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

	ph.Logger.Debugf("body bytes: %+v", buf)

	if packetItem, exists := (*p)[header.ID]; exists {
		if header.IsCrypt {
			decryptBody(buf)
		}
		err := packetItem.UnmarshalBinary(buf)
		if err != nil {
			return nil, err
		}
		packetItem.setHeader(header)

		return packetItem, nil
	}

	return nil, fmt.Errorf("unknown packet id: %d", header.ID)
}
