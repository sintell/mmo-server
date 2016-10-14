package packet

import (
	"encoding/binary"
	"math"
)

func readBytesAsBool(data []byte) bool {
	return readBytesAsUint8(data) > 0
}

func readBytesAsInt8(data []byte) int8 {
	return int8(data[0])
}

func readBytesAsUint8(data []byte) uint8 {
	return uint8(data[0])
}

func readBytesAsInt16(data []byte) int16 {
	return int16(binary.LittleEndian.Uint16(data))
}

func readBytesAsUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func readBytesAsInt32(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}

func readBytesAsUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func readBytesAsFloat32(data []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data))
}

func readBytesAsFloat64(data []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(data))
}

//////////////////////// PUT BYTES ///////////////////////////////

func putBoolAsBytes(data []byte, v bool) {
	if v {
		putUint8AsBytes(data, 1)
	} else {
		putUint8AsBytes(data, 0)
	}
}

func putInt8AsBytes(data []byte, v int8) {
	data[0] = uint8(v)
}

func putUint8AsBytes(data []byte, v uint8) {
	data[0] = uint8(v)
}

func putInt16AsBytes(data []byte, v int16) {
	binary.LittleEndian.PutUint16(data, uint16(v))
}

func putUint16AsBytes(data []byte, v uint16) {
	binary.LittleEndian.PutUint16(data, v)
}

func putInt32AsBytes(data []byte, v int32) {
	binary.LittleEndian.PutUint32(data, uint32(v))
}

func putUint32AsBytes(data []byte, v uint32) {
	binary.LittleEndian.PutUint32(data, v)
}

func putFloat32AsBytes(data []byte, v float32) {
	binary.LittleEndian.PutUint32(data, math.Float32bits(v))
}

func putFloat64AsBytes(data []byte, v float64) {
	binary.LittleEndian.PutUint64(data, math.Float64bits(v))
}
