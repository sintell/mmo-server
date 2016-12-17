package db

import (
	"io"
	"math"
	"net"

	"github.com/golang/glog"
	"github.com/sintell/mmo-server/packet"
)

// Provider represents database connection
type Provider struct {
	addr *net.TCPAddr
	pool []*net.TCPConn
}

// NewProvider creates and initialises a new data provider
func NewProvider(addr *net.TCPAddr, poolSize int) (*Provider, error) {
	p := &Provider{addr, make([]*net.TCPConn, poolSize)}
	err := p.connect()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Provider) connect() error {
	glog.V(10).Infof("connecting to DB server: %s", p.addr.String())
	c, err := net.DialTCP("tcp4", nil, p.addr)

	if err != nil {
		return err
	}

	p.pool[0] = c
	return nil
}

// Stop closes all connections
func (p *Provider) Stop() {
	for _, c := range p.pool {
		c.Close()
	}
}

func (p *Provider) getConn() *net.TCPConn {
	return p.pool[0]
}

// Query makes db requests
func (p *Provider) query(req []byte) (*packet.HeaderPacket, []byte, error) {
	glog.V(10).Infof("making request to db: % x", req)
	_, err := p.getConn().Write(req)
	if err != nil {
		if err == io.EOF {
			p.connect()
		}
		return nil, nil, err
	}
	buf := make([]byte, 6)
	read, err := p.getConn().Read(buf)
	if err != nil {
		if err == io.EOF {
			p.connect()
		}
		return nil, nil, err
	}
	glog.V(10).Infof("read queryl head: read: %d\tresponce: % x", read, buf)
	h := new(packet.HeaderPacket)
	h.UnmarshalBinary(buf)

	glog.V(10).Infof("header: %x", h)

	size := uint(h.Length - 6)
	buf = make([]byte, 0, size)
	t := make([]byte, size/2)
	bytesLeft := size
	for bytesLeft > 0 {
		read, err := p.getConn().Read(t[:int(math.Min(float64(bytesLeft), float64(size/2)))])
		if err != nil {
			if err == io.EOF {
				p.connect()
			}
			return nil, nil, err
		}

		buf = append(buf, t[:read]...)
		bytesLeft = bytesLeft - uint(read)
		glog.V(10).Infof("db responce: %d\tread: %d\tleft: %d\n", h.ID, read, bytesLeft)
	}

	return h, buf, nil
}

//GetActorsList TODO
func (p *Provider) GetActorsList(uid uint32) (*packet.ActorListPacket, error) {
	h, charListBuf, err := p.query((&packet.ActorListQueryPacket{
		HeaderPacket: packet.HeaderPacket{Length: 10, IsCrypt: false, Number: 0, ID: 11000},
		UID:          uid,
	}).MarshalBinary())
	if err != nil {
		glog.Warningf("character list query failed: %s", err.Error())
		return nil, err
	}

	cl := &packet.ActorListPacket{HeaderPacket: *h}
	err = cl.UnmarshalBinary(charListBuf)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (p *Provider) GetInventory(actorId uint32) ([]byte, error) {
	_, inventoryBuf, err := p.query((&packet.InventoryQueryPacket{
		HeaderPacket: packet.HeaderPacket{Length: 10, IsCrypt: false, Number: 0, ID: 11001},
		ActorID:      actorId,
	}).MarshalBinary())
	if err != nil {
		glog.Warningf("inventory query failed: %s", err.Error())
		return nil, err
	}

	return inventoryBuf, nil
}

func (p *Provider) RemoveItem(actorId uint32, itemID uint32) ([]byte, error) {
	_, inventoryBuf, err := p.query((&packet.InventoryQueryPacket{
		HeaderPacket: packet.HeaderPacket{Length: 10, IsCrypt: false, Number: 0, ID: 11002},
		ActorID:      actorId,
		ItemID:       itemID,
	}).MarshalBinary())
	if err != nil {
		glog.Warningf("inventory query failed: %s", err.Error())
		return nil, err
	}

	return inventoryBuf, nil
}
