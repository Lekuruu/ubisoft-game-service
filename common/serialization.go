package common

import (
	"encoding/binary"
	"math/big"
)

func WriteU32(value int) []byte {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(value))
	return result
}

func WriteU16(v int) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(v))
	return buf
}

func ReadAsU32List(buf []byte) []uint32 {
	size := len(buf) / 4
	result := make([]uint32, size)
	for i := 0; i < size; i++ {
		result[i] = binary.LittleEndian.Uint32(buf[i*4:])
	}
	return result
}

func WriteU32List(ints []uint32) []byte {
	buf := make([]byte, 4*len(ints))
	for i, val := range ints {
		binary.LittleEndian.PutUint32(buf[i*4:], val)
	}
	return buf
}

func WriteBigIntBE(value *big.Int, length int) []byte {
	result := make([]byte, length)
	copy(result[length-len(value.Bytes()):], value.Bytes())
	return result
}

func ReadBigIntBE(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}
