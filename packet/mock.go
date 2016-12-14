package packet

import (
	"encoding/binary"
	"os"

	"github.com/golang/glog"
)

// MockPacket TODO
type MockPacket struct {
	Data []byte
}

// MarshalBinary TODO: write doc
func (m *MockPacket) MarshalBinary() []byte {
	return m.Data
}

// UnmarshalBinary TODO: write doc
func (m *MockPacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: write doc
func (m *MockPacket) Header() *HeaderPacket {
	return nil
}

// UnmarshalBinary TODO: write doc
func (m *MockPacket) setHeader(h *HeaderPacket) {
}

// ReadMockPacket generates packet from binary file
func ReadMockPacket(path string) *MockPacket {
	glog.V(10).Infof("reading mock from: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDONLY, 0555)
	if err != nil {
		glog.Warningf("error loading mock: %s", err.Error())
		return nil
	}
	fInfo, err := f.Stat()
	if err != nil {
		glog.Warningf("error loading mock: %s", err.Error())
		return nil
	}
	p := &MockPacket{Data: make([]byte, fInfo.Size())}
	err = binary.Read(f, binary.LittleEndian, p.Data)
	if err != nil {
		glog.Warningf("error loading mock: %s", err.Error())
		return nil
	}

	return p
}
