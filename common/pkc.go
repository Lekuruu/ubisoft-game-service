package common

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

const MAX_RSA_MODULUS_LEN = 128
const PUBLIC_KEY_LEN = 512

func RsaKeygen() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, PUBLIC_KEY_LEN)
}

func RsaPublicKeyFromBuffer(data []byte) *rsa.PublicKey {
	key := rsa.PublicKey{}
	key.N = ReadBigIntBE(data[4 : MAX_RSA_MODULUS_LEN+4])
	key.E = int(ReadBigIntBE(data[MAX_RSA_MODULUS_LEN+4 : 2*MAX_RSA_MODULUS_LEN+4]).Int64())
	return &key
}

func RsaPublicKeyToBuffer(key *rsa.PublicKey) []byte {
	result := make([]byte, 0)
	result = append(result, WriteU32(PUBLIC_KEY_LEN)...)
	result = append(result, WriteBigIntBE(key.N, MAX_RSA_MODULUS_LEN)...)
	result = append(result, WriteBigIntBE(big.NewInt(int64(key.E)), MAX_RSA_MODULUS_LEN)...)
	return result
}
