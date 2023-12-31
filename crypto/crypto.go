package crypto

import (
	"crypto/sha256"

	//lint:ignore SA1019 use of deprecated package.
	"golang.org/x/crypto/ripemd160"
)

func Sha3(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func Ripemd160Hash(data []byte) []byte {
	hash := ripemd160.New()
	hash.Write(data)
	return hash.Sum(nil)
}
