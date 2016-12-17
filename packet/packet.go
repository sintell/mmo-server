package packet

import (
	"fmt"
	"io"
	"math"

	"github.com/golang/glog"
)

// AfterLoginPackets TODO
var AfterLoginPackets Packet

// CharListPacket TODO
var CharListPacket Packet

// StrangePacket TODO
var StrangePacket Packet

func init() {
	AfterLoginPackets = ReadMockPacket("./packets/login.pac")
	CharListPacket = ReadMockPacket("./packets/char.pac")
	StrangePacket = ReadMockPacket("./packets/idk.pac")
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
	PacketList PacketsList
}

// PacketsList holds all awailable packets
type PacketsList map[uint16]Packet

// NewPacketsList returns pointer to a new packets list
func (ph *GamePacketHandler) NewPacketsList() *PacketsList {
	return &PacketsList{
		1111: new(ClientLoginRequestPacket),
		5100: new(ClientLoginInfoPacket),
		5116: new(ActorLoginPacket),
		5188: new(ClientMovePacket),
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
	glog.V(10).Infof("head bytes: %+v", buf)

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
		glog.V(10).Infof("pid: %d\tread: %d\tleft: %d\n", header.ID, read, bytesLeft)
	}

	glog.V(10).Infof("body bytes: %+v", buf)

	if packetItem, exists := (*p)[header.ID]; exists {
		if header.IsCrypt {
			DecryptBody(buf)
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
