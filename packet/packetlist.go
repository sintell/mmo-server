package packet

type fp struct {
	D    uint32  // 4
	E    float32 // 4
	F    float32 // 4
	G    float32 // 4
	H    uint16  // 2
	Test int8    // 1
	I    float32 // 4
	J    int32   // 4
}

func (f *fp) MarshalBinary() []byte {
	return []byte{33, 10, 01, 10}
}

func (f *fp) UnmarshalBinary(data []byte) error {
	f.D = readBytesAsUint32(data[0:4])
	f.E = readBytesAsFloat32(data[4:8])
	f.F = readBytesAsFloat32(data[8:12])
	f.G = readBytesAsFloat32(data[12:16])
	f.H = readBytesAsUint16(data[16:18])
	f.Test = readBytesAsInt8(data[18:19])
	f.I = readBytesAsFloat32(data[19:23])
	f.J = readBytesAsInt32(data[23:27])

	return nil
}

type sp struct {
	//uint16 (500) (является количеством следующих элементов, длинной по 10 байт)
	D uint16
	//[
	E [500]struct {
		////int8 (9)
		F int8
		////float32 (500-600)
		G float32
		////int8 (0-5)
		H int8
		////int16 (1000-2000)
		I int16
		////int16 (15)
		J int16
	}
	//]
}

func (s *sp) MarshalBinary() []byte {
	return []byte{33, 10, 01, 10}
}

func (s *sp) UnmarshalBinary(data []byte) error {
	s.D = readBytesAsUint16(data[0:2])
	for i := 0; i < int(s.D); i++ {
		offset := i * 10
		s.E[i] = struct {
			////int8 (9)
			F int8
			////float32 (500-600)
			G float32
			////int8 (0-5)
			H int8
			////int16 (1000-2000)
			I int16
			////int16 (15)
			J int16
		}{
			readBytesAsInt8(data[offset+2 : offset+3]),
			readBytesAsFloat32(data[offset+3 : offset+7]),
			readBytesAsInt8(data[offset+7 : offset+8]),
			readBytesAsInt16(data[offset+8 : offset+10]),
			readBytesAsInt16(data[offset+10 : offset+12]),
		}
	}
	return nil
}

// HeaderPacket represents header for packets
type HeaderPacket struct {
	length  uint16
	isCrypt bool
	number  uint8
	ID      uint16
}

// MarshalBinary TODO: write doc
func (hp *HeaderPacket) MarshalBinary() []byte {
	return []byte{}
}

// UnmarshalBinary TODO: write doc
func (hp *HeaderPacket) UnmarshalBinary(data []byte) error {
	hp.length = readBytesAsUint16(data[0:2])
	hp.isCrypt = readBytesAsBool(data[2:3])
	hp.number = readBytesAsUint8(data[3:4])
	hp.ID = readBytesAsUint16(data[4:6])

	return nil
}

type clientLoginRequestPacket struct {
	isCrypt   bool
	number    uint8
	ID        uint16
	accountID uint32
	username  [20]byte
	password  [20]byte
	token     [32]byte
}

type clientLoginAccept struct {
	length  uint16
	isCrypt bool
	number  uint8
	ID      uint16
	token   [32]byte
}
