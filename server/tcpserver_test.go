package server

import (
	"io"
	"net"

	"github.com/sintell/mmo-server/packet"
)

func noop() {}

type DummyPacketHandler struct {
	HeadLength uint
}

func (dph DummyPacketHandler) ReadHead(c io.Reader) (uint, error) {
	return 0, nil
}

func (dph DummyPacketHandler) ReadBody(id uint, c io.Reader, pl *packet.PacketsList) (packet.Packet, error) {
	return nil, nil
}

func (dph DummyPacketHandler) NewPacketsList() *packet.PacketsList {
	return new(packet.PacketsList)
}

type ErrorPacketHandler struct {
	HeadLength    uint
	ReadBodyValue []byte
}

func (eph ErrorPacketHandler) ReadHead(c io.Reader) (*packet.HeaderPacket, error) {
	_, err := c.Read([]byte{})
	return new(packet.HeaderPacket), err
}

func (eph ErrorPacketHandler) ReadBody(header *packet.HeaderPacket, c io.Reader, pl *packet.PacketsList) (packet.Packet, error) {
	_, err := c.Read([]byte{})
	return &testPacket{2, 1}, err
}

func (eph ErrorPacketHandler) NewPacketsList() *packet.PacketsList {
	return new(packet.PacketsList)
}

type testPacket struct {
	id    uint
	field uint8
}

func (tp *testPacket) MarshalBinary() []byte {
	return []byte{byte(tp.id), byte(tp.field)}
}
func (tp *testPacket) UnmarshalBinary(data []byte) error {
	return nil
}

type DummyLogger struct {
}

func (l DummyLogger) Infof(format string, args ...interface{}) {
	noop()
}

func (l DummyLogger) Errorf(format string, args ...interface{}) {
	noop()
}

func (l DummyLogger) Debugf(format string, args ...interface{}) {
	noop()
}

type DummyAddr struct {
}

func (da *DummyAddr) String() string {
	return "0.0.0.0"
}
func (da *DummyAddr) Network() string {
	return "dummy network"
}

type DummyNetConnection struct {
	readCalled  uint
	writeCalled uint
	closeCalled uint
}

func (dnc *DummyNetConnection) Read(buf []byte) (int, error) {
	dnc.readCalled = dnc.readCalled + 1
	return 2, nil
}
func (dnc *DummyNetConnection) Write(buf []byte) (int, error) {
	dnc.writeCalled = dnc.writeCalled + 1
	return 2, nil
}
func (dnc *DummyNetConnection) Close() error {
	dnc.closeCalled = dnc.closeCalled + 1

	return nil
}
func (dnc *DummyNetConnection) RemoteAddr() net.Addr {
	return &DummyAddr{}
}

var (
	netAddr = &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 9999}
)

// func TestHandleConnection(t *testing.T) {
// 	server := TCPServer{
// 		NetAddr:       netAddr,
// 		PacketHandler: DummyPacketHandler{HeadLength: 2},
// 		Logger:        DummyLogger{},
// 	}
// 	dnc := &DummyNetConnection{}
//
// 	stop := server.handleConnection(dnc)
//
// 	stop <- true
// 	if dnc.closeCalled == 0 {
// 		t.Error("close was not called after stop channel close")
// 	}
// }
