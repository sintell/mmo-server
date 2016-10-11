package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

type ErrorTCPConnection struct {
}

func (ec *ErrorTCPConnection) Read([]byte) (int, error) {
	return 0, fmt.Errorf("read error")
}
func (ec *ErrorTCPConnection) Write([]byte) (int, error) {
	return 0, fmt.Errorf("write error")

}
func (ec *ErrorTCPConnection) Close() error {
	return nil
}
func (ec *ErrorTCPConnection) RemoteAddr() net.Addr {
	return &DummyAddr{}
}

func TestHandleConnectionError(t *testing.T) {
	if handleConnectionError(nil) {
		t.Error("want false, got true")
	}

	if !handleConnectionError(fmt.Errorf("error")) {
		t.Error("want true, got false")
	}

	if !handleConnectionError(io.EOF) {
		t.Error("want true, got false")
	}
}

func TestRecoverConnectionPanic(t *testing.T) {
	cm := ConnectionManager{
		PacketHandler: ErrorPacketHandler{},
		Logger:        DummyLogger{},
		Connections:   make(map[TCPConnection]bool),
		stop:          make(chan interface{}),
	}
	c := &DummyNetConnection{}
	func() {
		defer cm.recoverConnectionPanic(c)
		panic(fmt.Errorf("should panic"))
	}()

	if c.closeCalled == 0 {
		t.Error("connection close must be called on panic")
	}
}

func TestReadFrom(t *testing.T) {
	errorCalls := 0
	calls := 0
	ec := &ErrorTCPConnection{}
	testData := []byte{0x2, 0x1}
	cm := ConnectionManager{
		PacketHandler: ErrorPacketHandler{ReadBodyValue: testData},
		Logger:        DummyLogger{},
		Connections:   make(map[TCPConnection]bool),
		stop:          make(chan interface{}),
	}

	handleConnectionError = func(err error) bool {
		if err != nil {
			errorCalls = errorCalls + 1
			return true
		}
		calls = calls + 1
		return false
	}

	cm.ReadFrom(ec)

	if calls < 0 || errorCalls < 0 {
		t.Errorf("not handled errors: calls %d\terror calls: %d", calls, errorCalls)
	}

	c := &DummyNetConnection{}
	ch := cm.ReadFrom(c)
	data := <-ch
	if bytes.Compare(data.MarshalBinary(), testData) != 0 {

	}
	close(cm.stop)
	ch = cm.ReadFrom(c)
	<-time.After(time.Millisecond)
	if c.closeCalled == 0 {
		t.Error("should close connection on handler shutdown")
	}

	if _, open := <-ch; open {
		t.Error("should close sink on handler shutdown")
	}
}
