package common

import (
	"encoding/binary"
	"math/big"

	"golang.org/x/exp/constraints"
)

func WriteU32[T constraints.Integer](value T) []byte {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(value))
	return result
}

func WriteU16[T constraints.Integer](value T) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(value))
	return buf
}

func WriteU32List(ints []uint32) []byte {
	buf := make([]byte, 4*len(ints))
	for i, val := range ints {
		binary.LittleEndian.PutUint32(buf[i*4:], val)
	}
	return buf
}

func ReadAsU32List(data []byte) []uint32 {
	size := len(data) / 4
	result := make([]uint32, size)
	for i := 0; i < size; i++ {
		result[i] = binary.LittleEndian.Uint32(data[i*4:])
	}
	return result
}

func WriteBigIntBE(value *big.Int, length int) []byte {
	result := make([]byte, length)
	copy(result[length-len(value.Bytes()):], value.Bytes())
	return result
}

func ReadBigIntBE(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}
