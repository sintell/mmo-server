package game

import (
	"context"
	"sync"

	"github.com/golang/glog"
	"github.com/sintell/mmo-server/packet"
	"github.com/sintell/mmo-server/resourse"
)

// PacketsList holds all awailable packets
type PacketsList map[uint]packet.Packet

type actorList struct {
	list map[uint32][3]*resourse.ActorShort
	*sync.Mutex
}

type ingameActorList struct {
	actor map[uint32]*resourse.Actor
	*sync.Mutex
}

// Manager handles incoming game packets and controlls game mechanics and game resources
type Manager struct {
	binds          map[uint]data
	count          uint
	awaitingActors *actorList
	ingameActors   *ingameActorList
}

type data struct {
	source <-chan packet.Packet
	sink   chan<- packet.Packet
}

// NewManager initialized new game manager
func NewManager() *Manager {
	return &Manager{
		make(map[uint]data),
		0,
		&actorList{make(map[uint32][3]*resourse.ActorShort), new(sync.Mutex)},
		&ingameActorList{make(map[uint32]*resourse.Actor), new(sync.Mutex)},
	}
}

// RegisterDataSource adds source to the list of sources to listen from
func (m *Manager) RegisterDataSource(ctx context.Context, source <-chan packet.Packet) <-chan packet.Packet {
	sink := make(chan packet.Packet)
	m.binds[m.count] = data{source, sink}
	go m.handle(ctx, m.binds[m.count])
	return sink
}

func (m *Manager) handle(ctx context.Context, d data) {
	var uid uint32
	for p := range d.source {
		if p == nil || p.Header() == nil {
			d.sink <- p
			continue
		}
		glog.V(10).Infof("hadnling packet with ID: %d", p.Header().ID)
		switch p.Header().ID {
		case 1:
			ctx = p.(*packet.ContextSwitch).Ctx
			if id, ok := ctx.Value("UserID").(uint32); ok {
				glog.V(10).Infof("got uid: %d", uid)
				uid = id
			}

			if al, ok := ctx.Value("ActorsList").([3]*resourse.ActorShort); ok {
				glog.Infof("got actors list: %v", al)
				m.awaitingActors.Lock()
				m.awaitingActors.list[uid] = al
				m.awaitingActors.Unlock()
			}

		case 5116:
			al := p.(*packet.ActorLoginPacket)
			// Выбрать для пользователя его список
			m.awaitingActors.Lock()
			list, exists := m.awaitingActors.list[uid]
			m.awaitingActors.Unlock()

			if !exists {
				glog.Warningf("no character list for: %d", uid)
				continue
			}

			// Создать полного перса
			gameActor := new(resourse.Actor)

			for slotID, ac := range list {
				// Выбрать из списка перса с нужной айдишкой
				if ac.ID == al.ActorID {
					gameActor.ActorShort = *list[slotID]
				}
			}

			m.ingameActors.Lock()
			m.ingameActors.actor[uid] = gameActor
			m.ingameActors.Unlock()

			glog.V(10).Infof("selected actor: %v", gameActor)

		default:
			d.sink <- p
		}
	}
}
