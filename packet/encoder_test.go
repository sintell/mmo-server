package packet

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestReadBytesAsBool(t *testing.T) {
	t.Parallel()

	defer func() {
		if err := recover(); err != nil {
			t.Error("panic during bool conversion")
		}
	}()

	trueb := []byte{1}
	falseb := []byte{0}

	if !readBytesAsBool(trueb) {
		t.Error("expecting true, got false")
	}

	if readBytesAsBool(falseb) {
		t.Error("expectiong false, got true")
	}
}

func TestReadBytesAsInt8(t *testing.T) {
	t.Parallel()

	posb := []byte{math.MaxInt8}     // 0111 1111
	negb := []byte{math.MaxInt8 + 1} // 1111 1111

	res := readBytesAsInt8(posb)
	if res != math.MaxInt8 {
		t.Errorf("expected %d, got %d", math.MaxInt8, res)
	}

	res = readBytesAsInt8(negb)
	if res != math.MinInt8 {
		t.Errorf("expected %d, got %d", math.MinInt8, res)
	}
}

func TestReadBytesAsUint8(t *testing.T) {
	t.Parallel()

	b := []byte{math.MaxUint8}

	res := readBytesAsUint8(b)
	if res != math.MaxUint8 {
		t.Errorf("expected %d, got %d", math.MaxUint8, res)
	}
}

func TestReadBytesAsInt16(t *testing.T) {
	t.Parallel()

	posb := make([]byte, 2)
	negb := make([]byte, 2)

	binary.LittleEndian.PutUint16(posb, math.MaxInt16)
	binary.LittleEndian.PutUint16(negb, math.MaxInt16+1)

	res := readBytesAsInt16(posb)
	if res != math.MaxInt16 {
		t.Errorf("expected %d, got %d", math.MaxInt16, res)
	}

	res = readBytesAsInt16(negb)
	if res != math.MinInt16 {
		t.Errorf("expected %d, got %d", math.MinInt16, res)
	}
}

func TestReadBytesAsUint16(t *testing.T) {
	t.Parallel()

	b := make([]byte, 2)

	binary.LittleEndian.PutUint16(b, math.MaxUint16)

	res := readBytesAsUint16(b)
	if res != math.MaxUint16 {
		t.Errorf("expected %d, got %d", math.MaxUint16, res)
	}
}

func TestReadBytesAsInt32(t *testing.T) {
	t.Parallel()

	posb := make([]byte, 4)
	negb := make([]byte, 4)

	binary.LittleEndian.PutUint32(posb, math.MaxInt32)
	binary.LittleEndian.PutUint32(negb, math.MaxInt32+1)

	res := readBytesAsInt32(posb)
	if res != math.MaxInt32 {
		t.Errorf("expected %d, got %d", math.MaxInt32, res)
	}

	res = readBytesAsInt32(negb)
	if res != math.MinInt32 {
		t.Errorf("expected %d, got %d", math.MinInt32, res)
	}
}

func TestReadBytesAsUint32(t *testing.T) {
	t.Parallel()

	b := make([]byte, 4)

	binary.LittleEndian.PutUint32(b, math.MaxUint32)

	res := readBytesAsUint32(b)
	if res != math.MaxUint32 {
		t.Errorf("expected %d, got %d", math.MaxUint32, res)
	}
}

func TestReadBytesAsFloat32(t *testing.T) {
	t.Parallel()

	posb := make([]byte, 4)
	negb := make([]byte, 4)

	binary.LittleEndian.PutUint32(posb, math.Float32bits(math.MaxFloat32))
	binary.LittleEndian.PutUint32(negb, math.Float32bits(math.SmallestNonzeroFloat32))

	res := readBytesAsFloat32(posb)
	if res != math.MaxFloat32 {
		t.Errorf("expected %f, got %f", math.MaxFloat32, res)
	}

	res = readBytesAsFloat32(negb)
	if res != math.SmallestNonzeroFloat32 {
		t.Errorf("expected %f, got %f", math.SmallestNonzeroFloat32, res)
	}
}

func TestReadBytesAsFloat64(t *testing.T) {
	t.Parallel()

	posb := make([]byte, 8)
	negb := make([]byte, 8)

	binary.LittleEndian.PutUint64(posb, math.Float64bits(math.MaxFloat64))
	binary.LittleEndian.PutUint64(negb, math.Float64bits(math.SmallestNonzeroFloat64))

	res := readBytesAsFloat64(posb)
	if res != math.MaxFloat64 {
		t.Errorf("expected %f, got %f", math.MaxFloat64, res)
	}

	res = readBytesAsFloat64(negb)
	if res != math.SmallestNonzeroFloat64 {
		t.Errorf("expected %f, got %f", math.SmallestNonzeroFloat64, res)
	}
}
