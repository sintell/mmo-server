package server

import (
	"net"

	"github.com/sintell/mmo-server/game"
)

// Logger is interface for logging service
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// TCPServer is wrap around for basic tcp server listener
type TCPServer struct {
	NetAddr           *net.TCPAddr
	Logger            Logger
	ConnectionManager ConnectionManager
	GameManager       *game.Manager
}

// Listen start server and wait for incoming connections
func (s *TCPServer) Listen() {
	s.ConnectionManager.stop = make(chan interface{})
	srv, err := net.ListenTCP("tcp4", s.NetAddr)
	if err != nil {
		s.Logger.Errorf("error creating server: %s\n", err.Error())
	}

	s.Logger.Infof("server listening on: %s\n", s.NetAddr.String())

	for {
		conn, err := srv.AcceptTCP()
		if err != nil {
			s.Logger.Errorf("error accepting connection: %s\n", err.Error())
		}
		conn.SetKeepAlive(true)
		source := s.ConnectionManager.ReadFrom(conn)
		source = s.GameManager.RegisterDataSource(source)
		s.ConnectionManager.Write(conn, source)
	}
}

// Stop begins server teardown process
func (s *TCPServer) Stop() {
	s.ConnectionManager.CloseAll()
}
