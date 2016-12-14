package auth

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/sintell/mmo-server/packet"
)

type usersList struct {
	list map[uint32]*User
	*sync.Mutex
}

// Manager handles incoming auth packets and controlls game server authorization flow
type Manager struct {
	binds map[uint]data
	count uint
	rds   remoteDataSource

	users usersList
}

type remoteDataSource interface {
	Query([]byte) ([]byte, error)
}

type data struct {
	source <-chan packet.Packet
	sink   chan<- packet.Packet
}

// NewManager creates new auth manager
func NewManager(remoteSource remoteDataSource) *Manager {
	return &Manager{make(map[uint]data), 0, remoteSource, usersList{make(map[uint32]*User), new(sync.Mutex)}}
}

// RegisterDataSource make 1-way channels to read from and to write to
func (m *Manager) RegisterDataSource(source <-chan packet.Packet) <-chan packet.Packet {
	sink := make(chan packet.Packet)
	m.binds[m.count] = data{source, sink}
	go m.handle(m.binds[m.count])
	return sink
}

// GetUsers return users which are currently online
func (m *Manager) GetUsers() map[uint32]*User {
	m.users.Lock()
	defer m.users.Unlock()

	return m.users.list
}

func (m *Manager) getUser(uid uint32) (*User, error) {
	if u, exists := m.GetUsers()[uid]; exists {
		return u, nil
	}
	return nil, fmt.Errorf("user with id: %d not logged in", uid)
}

func (m *Manager) registerUser(u *User) error {
	m.users.Lock()
	defer m.users.Unlock()

	if oldUser, exists := m.users.list[u.uid]; exists {
		return fmt.Errorf("user already logged in: %s", oldUser)
	}

	m.users.list[u.uid] = u
	return nil
}

func (m *Manager) handle(d data) {
	for p := range d.source {
		glog.V(10).Infof("hadnling packet with ID: %d", p.Header().ID)
		switch p.Header().ID {
		case 1111:
			userData := p.(*packet.ClientLoginRequestPacket)
			resp := &packet.ClientLoginAcceptPacket{
				HeaderPacket: packet.HeaderPacket{Length: 39, IsCrypt: false, Number: 0, ID: 1112},
				Token:        userData.Token,
				Accepted:     true,
			}

			err := m.registerUser(&User{
				uid:       userData.AccountID,
				pwdHash:   userData.Password,
				account:   userData.Username,
				loginTime: time.Now(),
			})
			if err != nil {
				glog.Warningf("can't login user %d: %s", err.Error())
				resp.Accepted = false
			}

			d.sink <- resp
		case 5100:
			clientData := p.(*packet.ClientLoginInfoPacket)
			u, err := m.getUser(clientData.AccountID)
			glog.V(10).Infof("login from user: %s", u)
			err = u.checkCredentials(clientData.AccountID, clientData.Password)
			if err != nil {
				glog.Warningf("can't authenticate client: %s", err.Error())
				return
			}

			charListData, err := m.rds.Query((&packet.CharacterListQueryPacket{
				HeaderPacket: packet.HeaderPacket{Length: 10, IsCrypt: false, Number: 0, ID: 11000},
				UID:          u.uid,
			}).MarshalBinary())
			if err != nil {
				glog.Warningf("character list query failed: %s", err.Error())
			}

			resp := &packet.ServerTimePacket{
				HeaderPacket: packet.HeaderPacket{Length: 26, IsCrypt: false, Number: 0, ID: 5651},
			}
			d.sink <- resp
			d.sink <- packet.AfterLoginPackets
			d.sink <- &packet.MockPacket{Data: charListData}
		default:
			d.sink <- p
		}
	}
}
