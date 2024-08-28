package common

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"math/big"
)

const MAX_RSA_MODULUS_LEN = 128
const PUBLIC_KEY_LEN = 512

func RsaKeygen() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, PUBLIC_KEY_LEN)
}

func RsaPublicKeyFromBuffer(data []byte) *rsa.PublicKey {
	key := rsa.PublicKey{}
	key.N = readBigIntBE(data[4 : MAX_RSA_MODULUS_LEN+4])
	key.E = int(readBigIntBE(data[MAX_RSA_MODULUS_LEN+4 : 2*MAX_RSA_MODULUS_LEN+4]).Int64())
	return &key
}

func RsaPublicKeyToBuffer(key *rsa.PublicKey) []byte {
	result := make([]byte, 0)
	result = append(result, writeU32(PUBLIC_KEY_LEN)...)
	result = append(result, writeBigIntBE(key.N, MAX_RSA_MODULUS_LEN)...)
	result = append(result, writeBigIntBE(big.NewInt(int64(key.E)), MAX_RSA_MODULUS_LEN)...)
	return result
}

func writeU32(value int) []byte {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(value))
	return result
}

func writeBigIntBE(value *big.Int, length int) []byte {
	result := make([]byte, length)
	copy(result[length-len(value.Bytes()):], value.Bytes())
	return result
}

func readBigIntBE(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}
