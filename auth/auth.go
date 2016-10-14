package auth

import "github.com/sintell/mmo-server/packet"

// Manager handles incoming auth packets and controlls game server authorization flow
type Manager struct {
	binds map[uint]data
	count uint
}

type data struct {
	source <-chan packet.Packet
	sink   chan<- packet.Packet
}

// NewManager creates new auth manager
func NewManager() *Manager {
	return &Manager{make(map[uint]data), 0}
}

// RegisterDataSource make 1-way channels to read from and to write to
func (m *Manager) RegisterDataSource(source <-chan packet.Packet) <-chan packet.Packet {
	sink := make(chan packet.Packet)
	m.binds[m.count] = data{source, sink}
	go m.handle(m.binds[m.count])
	return sink
}

func (m *Manager) handle(d data) {
	for p := range d.source {
		switch p.Header().ID {
		case 1111:
			t := p.(*packet.ClientLoginRequestPacket)
			resp := &packet.ClientLoginAcceptPacket{
				HeaderPacket: packet.HeaderPacket{Length: 39, IsCrypt: false, Number: 0, ID: 1112},
				Token:        t.Token,
				Accepted:     true,
			}
			d.sink <- resp
		case 5100:
			_ = p.(*packet.ClientLoginInfoPacket)
			resp := &packet.ServerTimePacket{
				HeaderPacket: packet.HeaderPacket{Length: 26, IsCrypt: false, Number: 0, ID: 5651},
			}
			d.sink <- resp
			d.sink <- &packet.MockPacket{Data: packet.AfterLoginPackets}
			d.sink <- &packet.MockPacket{Data: packet.CharListPacket}
		}
	}
}
