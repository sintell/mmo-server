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

func DecryptHead(data []byte) []byte {
	transform(data, -3)
	return data
}

func EncryptHead(data []byte) []byte {
	transform(data, -3)
	return data
}

func DecryptBody(data []byte) []byte {
	transform(data, 3)
	return data
}

func EncryptBody(data []byte) []byte {
	transform(data, 3)
	return data
}

func Encrypt(data []byte) []byte {
	transform(data, -3)
	return data
}

// transform converts byte array to a crypted one
func transform(data []byte, offset int) {
	size := len(data)
	keyBuf := make([]byte, len(xorKey))
	copy(keyBuf, xorKey)
	c1 := uint16(0)
	c2 := uint16(0)
	kc := uint16(0)
	key := uint16(0)
	for i := 0 - offset; i < size; i++ {
		c1 = (c1 + 1) & 255
		key = uint16(keyBuf[4*c1+8])
		c2 = (key + c2) & 255
		kc = uint16(keyBuf[4*c2+8])
		keyBuf[4*c1+8] = uint8(kc)
		keyBuf[4*c2+8] = uint8(key)
		if i < 0 {
			continue
		}
		data[i] = data[i] ^ keyBuf[4*((key+kc)&255)+8]
	}
}
