package server

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/sintell/mmo-server/db"
	"github.com/sintell/mmo-server/packet"
)

// Logger is interface for logging service
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// GameManager is game manager
type GameManager interface {
	RegisterDataSource(context.Context, <-chan packet.Packet) <-chan packet.Packet
}

// AuthManager is auth manager
type AuthManager interface {
	RegisterDataSource(context.Context, <-chan packet.Packet) <-chan packet.Packet
}

// TCPServer is wrap around for basic tcp server listener
type TCPServer struct {
	NetAddr           *net.TCPAddr
	ConnectionManager ConnectionManager
	GameManager       GameManager
	AuthManager       AuthManager
	DB                db.Provider
	startTime         time.Time
}

func shouldReject(c TCPConnection) bool {
	if strings.Contains(c.RemoteAddr().String(), "0.0.0.0") ||
		strings.Contains(c.RemoteAddr().String(), "91.246.87.82") ||
		strings.Contains(c.RemoteAddr().String(), "91.246.101.158") ||
		strings.Contains(c.RemoteAddr().String(), "95.27.231.42") ||
		strings.Contains(c.RemoteAddr().String(), "138.201.123.151") ||
		strings.Contains(c.RemoteAddr().String(), "192.168.1.34") {

		return false
	}
	return true
}

// Listen start server and wait for incoming connections
func (s *TCPServer) Listen() {
	s.ConnectionManager.stop = make(chan interface{})
	srv, err := net.ListenTCP("tcp4", s.NetAddr)
	if err != nil {
		glog.Errorf("error creating server: %s\n", err.Error())
	}

	s.startTime = time.Now()

	glog.Infof("server listening on: %s\n", s.NetAddr.String())

	for {
		conn, err := srv.AcceptTCP()
		if err != nil {
			glog.Errorf("error accepting connection: %s\n", err.Error())
		}
		glog.V(10).Infof("got connection from: %s", conn.RemoteAddr().String())
		if shouldReject(conn) {
			glog.V(10).Infof("kicked: %s", conn.RemoteAddr().String())
			conn.Close()
			continue
		}
		conn.SetKeepAlive(true)
		_, err = conn.Write(packet.StrangePacket.MarshalBinary())
		if err != nil {
			glog.Errorf("error writing SP: %s", err.Error())
			return
		}

		ctx := context.WithValue(context.Background(), "UserConn", conn)

		source := s.ConnectionManager.ReadFrom(conn)
		source = s.AuthManager.RegisterDataSource(ctx, source)
		source = s.GameManager.RegisterDataSource(ctx, source)
		s.ConnectionManager.Write(conn, source)
	}
}

// Stop begins server teardown process
func (s *TCPServer) Stop() {
	s.ConnectionManager.CloseAll()
}
