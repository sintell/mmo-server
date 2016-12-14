package packet

import (
	"encoding/binary"
	"time"
)

// HeaderPacket represents header for packets
type HeaderPacket struct {
	Length  uint16
	IsCrypt bool
	Number  uint8
	ID      uint16
}

// MarshalBinary TODO: write doc
func (hp *HeaderPacket) MarshalBinary() []byte {
	buf := make([]byte, 6)
	putUint16AsBytes(buf[0:2], hp.Length)
	putBoolAsBytes(buf[2:3], hp.IsCrypt)
	putUint8AsBytes(buf[3:4], hp.Number)
	putUint16AsBytes(buf[4:6], hp.ID)
	if hp.IsCrypt {
		return encryptHead(buf)
	}
	return buf
}

// UnmarshalBinary TODO: write doc
func (hp *HeaderPacket) UnmarshalBinary(data []byte) error {
	hp.Length = readBytesAsUint16(data[0:2])
	hp.IsCrypt = readBytesAsBool(data[2:3])
	if hp.IsCrypt {
		decryptHead(data)
	}
	hp.Number = readBytesAsUint8(data[3:4])
	hp.ID = readBytesAsUint16(data[4:6])

	return nil
}

/////////////////////// RECV PACKETS ////////////////////////

// ClientLoginRequestPacket - packet with ID=1111
type ClientLoginRequestPacket struct {
	HeaderPacket
	AccountID uint32
	Username  string
	Password  string
	Token     string
}

// MarshalBinary TODO: write doc
func (clr *ClientLoginRequestPacket) MarshalBinary() []byte {
	return nil
}

// UnmarshalBinary TODO: write doc
func (clr *ClientLoginRequestPacket) UnmarshalBinary(data []byte) error {
	clr.AccountID = readBytesAsUint32(data[0:4])
	clr.Username = string(data[4:24])
	clr.Password = string(data[24:44])
	clr.Token = string(data[44:76])

	return nil
}

func (clr *ClientLoginRequestPacket) setHeader(h *HeaderPacket) {
	clr.HeaderPacket = *h
}

// Header TODO: write doc
func (clr *ClientLoginRequestPacket) Header() *HeaderPacket {
	return &clr.HeaderPacket
}

// ClientLoginInfoPacket - packet with ID=5100
type ClientLoginInfoPacket struct {
	HeaderPacket
	AccountID uint32
	Password  string
}

// MarshalBinary TODO: write doc
func (cli *ClientLoginInfoPacket) MarshalBinary() []byte {
	return nil
}

// UnmarshalBinary TODO: write doc
func (cli *ClientLoginInfoPacket) UnmarshalBinary(data []byte) error {
	cli.AccountID = readBytesAsUint32(data[0:4])
	cli.Password = string(data[12:32])

	return nil
}

func (cli *ClientLoginInfoPacket) setHeader(h *HeaderPacket) {
	cli.HeaderPacket = *h
}

// Header TODO: write doc
func (cli *ClientLoginInfoPacket) Header() *HeaderPacket {
	return &cli.HeaderPacket
}

/////////////////////// SEND PACKETS ////////////////////////

// ClientLoginAcceptPacket - packet with ID=1112
type ClientLoginAcceptPacket struct {
	HeaderPacket
	Token    string
	Accepted bool
}

// MarshalBinary TODO: write doc
func (cla *ClientLoginAcceptPacket) MarshalBinary() []byte {
	buf := make([]byte, cla.Length)
	copy(buf[:6], cla.HeaderPacket.MarshalBinary())
	copy(buf[6:38], []byte(cla.Token))
	putBoolAsBytes(buf[38:39], cla.Accepted)
	return buf
}

// UnmarshalBinary TODO: write doc
func (cla *ClientLoginAcceptPacket) UnmarshalBinary([]byte) error {
	return nil
}

func (cla *ClientLoginAcceptPacket) setHeader(h *HeaderPacket) {
	cla.HeaderPacket = *h
}

// Header TODO: write doc
func (cla *ClientLoginAcceptPacket) Header() *HeaderPacket {
	return &cla.HeaderPacket
}

// ServerTimePacket - packet with ID=5651
type ServerTimePacket struct {
	HeaderPacket
	MsSinceStart uint32
	Year         uint16
	Month        uint16
	DayOfWeek    uint16
	DayNumber    uint16
	Hour         uint16
	Minute       uint16
	Second       uint16
	Millisecond  uint16
}

// MarshalBinary TODO: write doc
func (st *ServerTimePacket) MarshalBinary() []byte {
	buf := make([]byte, st.HeaderPacket.Length)
	copy(buf[:6], st.HeaderPacket.MarshalBinary())
	t := time.Now()
	y, m, d := t.Date()
	hr, min, s := t.Clock()
	ms := t.Nanosecond()
	dow := t.Weekday()

	binary.LittleEndian.PutUint32(buf[6:10], uint32(st.MsSinceStart))
	binary.LittleEndian.PutUint16(buf[10:12], uint16(y))
	binary.LittleEndian.PutUint16(buf[12:14], uint16(m))
	binary.LittleEndian.PutUint16(buf[14:16], uint16(dow))
	binary.LittleEndian.PutUint16(buf[16:18], uint16(d))
	binary.LittleEndian.PutUint16(buf[18:20], uint16(hr))
	binary.LittleEndian.PutUint16(buf[20:22], uint16(min))
	binary.LittleEndian.PutUint16(buf[22:24], uint16(s))
	binary.LittleEndian.PutUint16(buf[24:26], uint16(ms))

	return []byte{}
}

// UnmarshalBinary TODO: write doc
func (st *ServerTimePacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: write doc
func (st *ServerTimePacket) Header() *HeaderPacket {
	return &st.HeaderPacket
}

func (st *ServerTimePacket) setHeader(h *HeaderPacket) {
	st.HeaderPacket = *h
}

/////////////////////// RDS PACKETS ////////////////////////

// CharacterListQueryPacket TODO
type CharacterListQueryPacket struct {
	HeaderPacket
	UID uint32
}

// MarshalBinary TODO: write doc
func (clq *CharacterListQueryPacket) MarshalBinary() []byte {
	buf := make([]byte, clq.HeaderPacket.Length)
	copy(buf[:6], clq.HeaderPacket.MarshalBinary())
	putUint32AsBytes(buf[6:10], clq.UID)
	return buf
}

// UnmarshalBinary TODO: write doc
func (clq *CharacterListQueryPacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: write doc
func (clq *CharacterListQueryPacket) Header() *HeaderPacket {
	return nil
}

// UnmarshalBinary TODO: write doc
func (clq *CharacterListQueryPacket) setHeader(h *HeaderPacket) {
}
