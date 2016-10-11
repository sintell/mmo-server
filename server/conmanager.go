package server

import (
	"io"
	"time"

	"github.com/sintell/mmo-server/packet"
)

// Packet represent datagramms and provides a convinient set of operations
type Packet interface {
	UnmarshalBinary([]byte) error
}

// PacketsList holds all awailable packets
type PacketsList map[uint]Packet

// PacketHandler represent a handler of any client-server datagramms
type PacketHandler interface {
	ReadHead(io.Reader) (uint, error)
	ReadBody(uint, io.Reader, *packet.PacketsList) (packet.Packet, error)
	NewPacketsList() *packet.PacketsList
}

// ConnectionManager operates with connections
type ConnectionManager struct {
	PacketHandler PacketHandler
	Logger        Logger
	Connections   map[TCPConnection]bool
	stop          chan interface{}
}

var handleConnectionError = func(err error) bool {
	if err != nil {
		if err == io.EOF {
			return true
		}
		return true
	}
	return false
}

func (cm *ConnectionManager) recoverConnectionPanic(c TCPConnection) {
	if err := recover(); err != nil {
		cm.Logger.Errorf("catched panic: %s", err.(error).Error())
		c.Close()
	}

	delete(cm.Connections, c)
	c.Close()
}

// ReadFrom handles connection reads
func (cm *ConnectionManager) ReadFrom(c TCPConnection) <-chan packet.Packet {
	sink := make(chan packet.Packet)
	cm.Connections[c] = true

	go func() {
		defer c.Close()
		defer close(sink)
		defer cm.recoverConnectionPanic(c)

		for {
			select {
			case _, more := <-cm.stop:
				if !more {
					cm.Logger.Infof("abort read: [%s]", c.RemoteAddr().String())
					return
				}
			default:
				{
					t := time.Now()
					pid, err := cm.PacketHandler.ReadHead(c)
					if handleConnectionError(err) {
						return
					}
					cm.Logger.Debugf("read pid of %d", pid)
					data, err := cm.PacketHandler.ReadBody(pid, c, cm.PacketHandler.NewPacketsList())
					if handleConnectionError(err) {
						cm.Logger.Errorf("error reading packet body: %s\n", err.Error())
						return
					}
					cm.Logger.Debugf("packet read complete in %s\n", time.Since(t).String())
					sink <- data
				}
			}
		}
	}()

	return sink
}

func (cm *ConnectionManager) Write(c TCPConnection, source <-chan packet.Packet) {
	go func() {
		defer c.Close()
		defer cm.Logger.Infof("closing connection for: [%s]", c.RemoteAddr().String())

		for p := range source {
			select {
			case _, more := <-cm.stop:
				if !more {
					cm.Logger.Infof("abort write: [%s]", c.RemoteAddr().String())
					return
				}

			default:
			}

			if p == nil {
				return
			}
			t := time.Now()
			_, err := c.Write(p.MarshalBinary())
			if handleConnectionError(err) {
				cm.Logger.Errorf("error writing data: %s", err.Error())
			}
			cm.Logger.Debugf("packet write complete in %s\n", time.Since(t).String())

		}
	}()
}

// CloseAll closes all connections before shutdown
func (cm *ConnectionManager) CloseAll() {
	for c := range cm.Connections {
		cm.Logger.Infof("stop requested, closing all connections")
		close(cm.stop)
		c.Close()
	}
}
