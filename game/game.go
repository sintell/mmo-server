package game

import "github.com/sintell/mmo-server/packet"

// PacketsList holds all awailable packets
type PacketsList map[uint]packet.Packet

// Manager handles incoming game packets and controlls game mechanics and game resources
type Manager struct {
	binds map[uint]data
	count uint
}

type data struct {
	source <-chan packet.Packet
	sink   chan<- packet.Packet
}

// NewManager initialized new game manager
func NewManager() *Manager {
	return &Manager{make(map[uint]data), 0}
}

// RegisterDataSource adds source to the list of sources to listen from
func (m *Manager) RegisterDataSource(source <-chan packet.Packet) <-chan packet.Packet {
	sink := make(chan packet.Packet)
	m.binds[m.count] = data{source, sink}
	go m.handle(m.binds[m.count])
	return sink
}

func (m *Manager) handle(d data) {
	for {
		select {
		case packet := <-d.source:
			{
				d.sink <- packet
			}
		}
	}
}
