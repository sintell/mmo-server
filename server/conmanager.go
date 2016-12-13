package server

import (
	"io"
	"time"

	"github.com/golang/glog"
	"github.com/sintell/mmo-server/packet"
)

// PacketsList holds all awailable packets
type PacketsList map[uint]packet.Packet

// PacketHandler represent a handler of any client-server datagramms
type PacketHandler interface {
	ReadHead(io.Reader) (*packet.HeaderPacket, error)
	ReadBody(*packet.HeaderPacket, io.Reader, *packet.PacketsList) (packet.Packet, error)
	NewPacketsList() *packet.PacketsList
}

// ConnectionManager operates with connections
type ConnectionManager struct {
	PacketHandler PacketHandler
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
		glog.Errorf("catched panic: %s", err.(error).Error())
		c.Close()
		return
	}

	delete(cm.Connections, c)
	c.Close()
}

// ReadFrom handles connection reads
func (cm *ConnectionManager) ReadFrom(c TCPConnection) <-chan packet.Packet {
	sink := make(chan packet.Packet)
	cm.Connections[c] = true

	go func() {
		defer func() {
			cm.recoverConnectionPanic(c)
			close(sink)
		}()

		for {
			select {
			case _, more := <-cm.stop:
				if !more {
					glog.Infof("abort read: [%s]", c.RemoteAddr().String())
					return
				}
			default:
				{
					header, err := cm.PacketHandler.ReadHead(c)
					t := time.Now()
					if handleConnectionError(err) {
						glog.Errorf("error reading packet header: %s\n", err.Error())
						return
					}
					data, err := cm.PacketHandler.ReadBody(header, c, cm.PacketHandler.NewPacketsList())
					if handleConnectionError(err) {
						glog.Errorf("error reading packet body: %s\n", err.Error())
						return
					}
					glog.V(10).Infof("packet read complete in %s\n", time.Since(t).String())
					sink <- data
				}
			}
		}
	}()

	return sink
}

func (cm *ConnectionManager) Write(c TCPConnection, source <-chan packet.Packet) {
	go func() {
		defer func() {
			c.Close()
			glog.Infof("closing connection for: [%s]", c.RemoteAddr().String())
		}()

		for p := range source {
			select {
			case _, more := <-cm.stop:
				if !more {
					glog.Infof("abort write: [%s]", c.RemoteAddr().String())
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
				glog.Errorf("error writing data: %s", err.Error())
			}
			glog.V(10).Infof("packet write complete in %s\n", time.Since(t).String())
		}
	}()
}

// RegisterFilters add filtering functions to connection
func (cm *ConnectionManager) RegisterFilters(source <-chan packet.Packet,
	filters ...func(packet.Packet) packet.Packet) <-chan packet.Packet {
	sink := make(chan packet.Packet)

	go func() {
		for p := range source {
			for _, filter := range filters {
				if pass := filter(p); pass != nil {
					sink <- pass
				}
			}
		}
	}()

	return sink
}

// CloseAll closes all connections before shutdown
func (cm *ConnectionManager) CloseAll() {
	glog.Infof("stop requested, closing all connections")
	close(cm.stop)
	for c := range cm.Connections {
		c.Close()
	}
}
