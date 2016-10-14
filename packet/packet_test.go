package packet

import (
	"bytes"
	"testing"
)

func noop() {}

type TestPacket struct {
	id    uint
	field uint8
}

func (tp *TestPacket) MarshalBinary() []byte {
	return []byte{byte(tp.id), byte(tp.field)}
}
func (tp *TestPacket) UnmarshalBinary(data []byte) error {
	return nil
}

func (tp *TestPacket) setHeader(h *HeaderPacket) {
}
func (tp *TestPacket) Header() *HeaderPacket {
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

func TestReadHead(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("panic while reading HeaderPacket: %s", err.(error).Error())
		}
	}()
	t.Parallel()

	testPacket := []byte{0x10, 0x0, 0x0, 0x0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	emptyPacket := []byte{}

	ph := GamePacketHandler{}
	header, err := ph.ReadHead(bytes.NewBuffer(testPacket))
	if err == nil {
		t.Error("calling ReadHead with zero HeadLength yeilds no error")
	}

	if header != nil {
		t.Error("calling ReadHead with zero HeadLength results in header not being nil")
	}

	ph = GamePacketHandler{HeadLength: 6}
	header, err = ph.ReadHead(bytes.NewBuffer(emptyPacket))
	if err == nil {
		t.Error("reading empty packet yeilds no error")
	}

	if header != nil {
		t.Error("reading empty packet results in in header not being nil")
	}

	header, err = ph.ReadHead(bytes.NewBuffer(testPacket))
	if err != nil {
		t.Errorf("error reading head: %s", err.Error())
	}
	if header.ID != 0x10 {
		t.Errorf("wrong header.ID: expected %d, got %d", 16, header.ID)
	}
}
