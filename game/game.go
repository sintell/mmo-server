package game

import (
	"context"
	"sync"

	"github.com/golang/glog"
	"github.com/sintell/mmo-server/packet"
	"github.com/sintell/mmo-server/resource"
)

// PacketsList holds all awailable packets
type PacketsList map[uint]packet.Packet

type actorList struct {
	list map[uint32][3]*resource.ActorShort
	*sync.Mutex
}

type ingameActorList struct {
	actor map[uint32]*resource.Actor
	*sync.Mutex
}

// Manager handles incoming game packets and controlls game mechanics and game resources
type Manager struct {
	binds          map[uint]data
	count          uint
	awaitingActors *actorList
	ingameActors   *ingameActorList
	rds            dataSource
}

type dataSource interface {
	GetInventory(uint32) ([]byte, error)
	RemoveItem(uint32, uint32, int32) (*packet.RemoveResultPacket, error)
	AddItem(uint32, []byte, int32) (*packet.AddItemResultPacket, error)
}

type data struct {
	source <-chan packet.Packet
	sink   chan<- packet.Packet
}

// NewManager initialized new game manager
func NewManager(remoteSource dataSource) *Manager {
	return &Manager{
		make(map[uint]data),
		0,
		&actorList{make(map[uint32][3]*resource.ActorShort), new(sync.Mutex)},
		&ingameActorList{make(map[uint32]*resource.Actor), new(sync.Mutex)},
		remoteSource,
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
	var gameActor *resource.Actor
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

			if al, ok := ctx.Value("ActorsList").([3]*resource.ActorShort); ok {
				glog.V(10).Infof("got actors list: %v", al)
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
			gameActor = new(resource.Actor)

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

			inventoryBuf, err := m.rds.GetInventory(gameActor.ID)

			if err != nil {
				glog.Warningf("no inventory list for: %d", gameActor.ID)
				continue
			}

			d.sink <- &packet.LoginInWorldPacket{
				HeaderPacket:  packet.HeaderPacket{Length: 9066, IsCrypt: false, Number: 0, ID: 5117},
				ActorData:     gameActor,
				InventoryData: inventoryBuf,
			}
		case 5188:
			cm := p.(*packet.ClientMovePacket)
			if gameActor == nil {
				glog.Warningf("no gameActor for %d", uid)
				continue
			}
			addItemResult, err := AddItem(m.rds, gameActor, 409, 1, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x99, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, resource.AddItemNormal, 5188)
			if err != nil {
				glog.Warningf("can't add item")
				continue
			}
			d.sink <- &packet.ServerMovePacket{
				HeaderPacket:     packet.HeaderPacket{Length: 33, IsCrypt: false, Number: 0, ID: 5189},
				ClientMovePacket: *cm,
				UniqueID:         gameActor.UniqueID,
				SpeedMove:        350, //TODO SpeedMove (gameActor.Stats.SpeedMove)
			}
			d.sink <- addItemResult
		case 5528:
			ri := p.(*packet.RemoveItemPacket)
			if gameActor == nil {
				glog.Warningf("no gameActor for %d", uid)
				continue
			}
			glog.V(10).Infof("Remove item ID: %d. Amount: %d", ri.ID, ri.Amount)
			removeItemResult, err := RemoveItem(m.rds, gameActor.ID, gameActor.UniqueID, ri.UniqueID, ri.Amount, resource.InventoryItemLoss, p.Header().ID)

			if err != nil {
				glog.Warningf("can't remove item %d, from gameActor %d", ri.UniqueID, gameActor.ID)
				continue
			}

			d.sink <- removeItemResult
		default:
			d.sink <- p
		}
	}
}
