package packet

type fp struct {
	ID   uint16
	A    bool    // 1
	B    uint8   // 1
	C    uint16  // 2
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
	f.ID = 33
	f.A = readBytesAsBool(data[0:1])
	f.B = readBytesAsUint8(data[1:2])
	f.C = readBytesAsUint16(data[2:4])
	f.D = readBytesAsUint32(data[4:8])
	f.E = readBytesAsFloat32(data[8:12])
	f.F = readBytesAsFloat32(data[12:16])
	f.G = readBytesAsFloat32(data[16:20])
	f.H = readBytesAsUint16(data[20:22])
	f.Test = readBytesAsInt8(data[22:23])
	f.I = readBytesAsFloat32(data[23:27])
	f.J = readBytesAsInt32(data[27:31])

	return nil
}

type sp struct {
	//bool (1, true)
	A bool
	//uint8 (250 максимум, потом обнуляет) (нечетное)
	B uint8
	//uint16 (5190)
	C uint16
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
	s.A = readBytesAsBool(data[0:1])
	s.B = readBytesAsUint8(data[1:2])
	s.C = readBytesAsUint16(data[2:4])
	s.D = readBytesAsUint16(data[4:6])
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
			readBytesAsInt8(data[offset+6 : offset+7]),
			readBytesAsFloat32(data[offset+7 : offset+11]),
			readBytesAsInt8(data[offset+11 : offset+12]),
			readBytesAsInt16(data[offset+12 : offset+14]),
			readBytesAsInt16(data[offset+14 : offset+16]),
		}
	}
	return nil
}
