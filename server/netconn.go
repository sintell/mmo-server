package server

import "net"

// TCPConnection wraps connection with custom interface for portability
type TCPConnection interface {
	Close() error
	RemoteAddr() net.Addr
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}
