package packet

import (
	"encoding/binary"
	"os"
)

var xorKey []byte

func init() {
	f, err := os.OpenFile("./key/packets.key", os.O_RDONLY, 0555)
	if err != nil {
		return
	}
	fInfo, err := f.Stat()
	if err != nil {
		return
	}

	xorKey = make([]byte, fInfo.Size())
	err = binary.Read(f, binary.LittleEndian, xorKey)
	if err != nil {
		return
	}
}

func decryptHead(data []byte) []byte {
	Transform(data, -3)
	return data
}

func encryptHead(data []byte) []byte {
	Transform(data, -3)
	return data
}

func decryptBody(data []byte) []byte {
	Transform(data, 3)
	return data
}

// Transform converts byte array to a crypted one
func Transform(data []byte, offset int) {
	size := len(data)
	buf := make([]byte, len(xorKey))
	copy(buf, xorKey)
	c1 := uint8(0)
	c2 := uint16(0)
	kc := uint8(0)
	key := uint8(0)
	for i := 0 - offset; i < size; i++ {
		c1++
		key = buf[4*c1+8]
		c2 = (uint16(key) + c2) & 255
		kc = buf[4*c2+8]
		buf[4*c1+8] = kc
		buf[4*c2+8] = key
		if i < 0 {
			continue
		}
		data[i] = data[i] ^ buf[4*uint16(uint8(key+kc))+8]
	}

}
