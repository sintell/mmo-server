package packet

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	"github.com/sintell/mmo-server/resourse"
)

// HeaderPacket represents header for packets
type HeaderPacket struct {
	Length  uint16
	IsCrypt bool
	Number  uint8
	ID      uint16

	Internal bool
}

// MarshalBinary TODO: write doc
func (hp *HeaderPacket) MarshalBinary() []byte {
	buf := make([]byte, 6)
	putUint16AsBytes(buf[0:2], hp.Length)
	putBoolAsBytes(buf[2:3], hp.IsCrypt)
	putUint8AsBytes(buf[3:4], hp.Number)
	putUint16AsBytes(buf[4:6], hp.ID)

	return buf
}

// UnmarshalBinary TODO: write doc
func (hp *HeaderPacket) UnmarshalBinary(data []byte) error {
	hp.Length = readBytesAsUint16(data[0:2])
	hp.IsCrypt = readBytesAsBool(data[2:3])
	if hp.IsCrypt {
		DecryptHead(data)
	}
	hp.Number = readBytesAsUint8(data[3:4])
	hp.ID = readBytesAsUint16(data[4:6])

	return nil
}

/////////////////////// INTERNAL PACKETS /////////////////////

// ContextSwitch TODO
type ContextSwitch struct {
	HeaderPacket

	Ctx context.Context
}

func (cs *ContextSwitch) MarshalBinary() []byte             { return nil }
func (cs *ContextSwitch) UnmarshalBinary(data []byte) error { return nil }
func (cs *ContextSwitch) Header() *HeaderPacket {
	return &cs.HeaderPacket
}
func (cs *ContextSwitch) setHeader(h *HeaderPacket) {
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

	return buf
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

// ActorLoginPacket
type ActorLoginPacket struct {
	HeaderPacket
	ActorID uint32
}

// MarshalBinary TODO: write doc
func (cl *ActorLoginPacket) MarshalBinary() []byte {
	return nil
}

// UnmarshalBinary TODO: write doc
func (cl *ActorLoginPacket) UnmarshalBinary(data []byte) error {
	cl.ActorID = readBytesAsUint32(data[0:4])
	return nil
}

// Header TODO: write doc
func (cl *ActorLoginPacket) Header() *HeaderPacket {
	return &cl.HeaderPacket
}

func (cl *ActorLoginPacket) setHeader(h *HeaderPacket) {
	cl.HeaderPacket = *h
}

/////////////////////// RDS PACKETS ////////////////////////

// ActorListQueryPacket TODO
type ActorListQueryPacket struct {
	HeaderPacket
	UID uint32
}

// MarshalBinary TODO: write doc
func (clq *ActorListQueryPacket) MarshalBinary() []byte {
	buf := make([]byte, clq.HeaderPacket.Length)
	copy(buf[:6], clq.HeaderPacket.MarshalBinary())
	putUint32AsBytes(buf[6:10], clq.UID)
	return buf
}

// UnmarshalBinary TODO: write doc
func (clq *ActorListQueryPacket) UnmarshalBinary(data []byte) error {
	return nil
}

// Header TODO: write doc
func (clq *ActorListQueryPacket) Header() *HeaderPacket {
	return nil
}

// UnmarshalBinary TODO: write doc
func (clq *ActorListQueryPacket) setHeader(h *HeaderPacket) {
}

// ActorListPacket TODO
type ActorListPacket struct {
	HeaderPacket
	SlotsTaken uint64
	List       [3]*resourse.ActorShort
}

// MarshalBinary TODO: write doc
func (cl *ActorListPacket) MarshalBinary() []byte {
	// 1528 is character list packet length
	buf := make([]byte, 1528)
	copy(buf[0:6], cl.HeaderPacket.MarshalBinary())
	offset := 0
	for cell, character := range cl.List {
		if character == nil {
			continue
		}

		offset = cell * 144
		putUint32AsBytes(buf[11+offset:11+offset+4], character.ID)
		putUint16AsBytes(buf[15+offset:15+offset+2], character.Class)
		putUint8AsBytes(buf[19+offset:19+offset+1], character.Sex)
		putUint8AsBytes(buf[20+offset:20+offset+1], character.Hair)
		putUint8AsBytes(buf[21+offset:21+offset+1], character.Face)
		putUint16AsBytes(buf[87+offset:87+offset+2], character.Level)
		copy(buf[89+offset:], []byte(character.Name))
		putUint32AsBytes(buf[115+offset:115+offset+4], character.ID)
		putUint16AsBytes(buf[123+offset:123+offset+2], character.Level)

		copy(buf[cell*320+439:], character.Equipment.Data)

		putUint16AsBytes(buf[1435+cell*4:1435+cell*4+2], character.Str)
		putUint16AsBytes(buf[1447+cell*4:1447+cell*4+2], character.Int)
		putUint16AsBytes(buf[1459+cell*4:1459+cell*4+2], character.Dex)
		putInt32AsBytes(buf[1471+cell*4:1471+cell*4+4], character.Rating)

		putFloat32AsBytes(buf[1483+cell*12:1483+cell*12+4], character.X)
		putFloat32AsBytes(buf[1487+cell*12:1487+cell*12+4], character.Y)
		putFloat32AsBytes(buf[1491+cell*12:1491+cell*12+4], character.Z)
	}
	return buf
}

// UnmarshalBinary TODO: write doc
func (cl *ActorListPacket) UnmarshalBinary(data []byte) error {
	var bOffset int
	var equipOffset int
	var cell int
	var name string

	cl.SlotsTaken = readBytesAsUint64(data[8 : 8+8])
	for slot := 0; slot < int(cl.SlotsTaken); slot++ {
		name = ""
		bOffset = 24 + (slot * 112)
		cell = int(readBytesAsUint8(data[8+bOffset : 8+bOffset+2]))
		equipOffset = 24 + int(cl.SlotsTaken)*112

		// This is beacause name can contain some trash
		// We should split them till first \x00 symbol
		// And took the first part
		name = string(bytes.Runes(bytes.Split(data[bOffset+10:bOffset+22], []byte{0x00})[0]))

		cl.List[cell] = &resourse.ActorShort{
			ID:    readBytesAsUint32(data[bOffset : bOffset+4]),
			Name:  name,
			Class: readBytesAsUint16(data[bOffset+22 : bOffset+24]),
			Level: readBytesAsUint16(data[bOffset+24 : bOffset+26]),
			Appearance: resourse.Appearance{
				Sex:  readBytesAsUint8(data[bOffset+40 : bOffset+41]),
				Hair: readBytesAsUint8(data[bOffset+41 : bOffset+42]),
				Face: readBytesAsUint8(data[bOffset+42 : bOffset+43]),
			},
			Stats: resourse.Stats{
				Str:    readBytesAsUint16(data[bOffset+44 : bOffset+46]),
				Dex:    readBytesAsUint16(data[bOffset+46 : bOffset+48]),
				Int:    readBytesAsUint16(data[bOffset+48 : bOffset+50]),
				Rating: readBytesAsInt32(data[bOffset+60 : bOffset+64]),
			},
			Position: resourse.Position{
				X: readBytesAsFloat32(data[bOffset+64 : bOffset+68]),
				Y: readBytesAsFloat32(data[bOffset+68 : bOffset+72]),
				Z: readBytesAsFloat32(data[bOffset+72 : bOffset+76]),
			},
			Equipment: resourse.Equipment{
				Data: data[(slot+1)*equipOffset : (slot+1)*equipOffset+(16*20)],
			},
		}
	}

	return nil
}

// Header TODO: write doc
func (cl *ActorListPacket) Header() *HeaderPacket {
	return &cl.HeaderPacket
}

// UnmarshalBinary TODO: write doc
func (cl *ActorListPacket) setHeader(h *HeaderPacket) {
	cl.HeaderPacket = *h
}
