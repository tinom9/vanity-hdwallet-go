package derivations

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/pbkdf2"
)

const (
	PBKDF2Rounds  = 2048
	HardenedIndex = 0x80000000
)

func mnemonicToSeed(mnemonic string, passphrase string) []byte {
	passphrase = "mnemonic" + passphrase
	mnemonicBytes := []byte(mnemonic)
	passphraseBytes := []byte(passphrase)
	stretched := pbkdf2.Key(mnemonicBytes, passphraseBytes, PBKDF2Rounds, 64, sha512.New)
	return stretched
}

func derivePrivateKeyFromSeed(seed []byte) ([]byte, []byte) {
	key := []byte("Bitcoin seed")
	hmac512 := hmac.New(sha512.New, key)
	hmac512.Write(seed)
	secret := hmac512.Sum(nil)
	return secret[:32], secret[32:]
}

func privKeyToPubKey(privKey []byte) []byte {
	pubKeyX, pubKeyY := secp256k1.S256().ScalarBaseMult(privKey)
	return secp256k1.CompressPubkey(pubKeyX, pubKeyY)
}

func getPathList(path string) []uint32 {
	indexes := strings.Split(path, "/")[1:]
	listPath := make([]uint32, len(indexes))
	for i, idx := range indexes {
		if strings.HasSuffix(idx, "'") || strings.HasSuffix(idx, "h") || strings.HasSuffix(idx, "H") {
			n, _ := strconv.Atoi(idx[:len(idx)-1])
			listPath[i] = uint32(n) + HardenedIndex
		} else {
			n, _ := strconv.Atoi(idx)
			listPath[i] = uint32(n)
		}
	}
	return listPath
}

func derivePrivateChild(pubKey, privKey, chainCode []byte, index uint32, hardened bool) ([]byte, []byte) {
	var key []byte
	if hardened {
		key = append([]byte{0x00}, privKey...)
	} else {
		key = pubKey
	}
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, index)
	hmac512 := hmac.New(sha512.New, chainCode)
	hmac512.Write(append(key, buf...))
	payload := hmac512.Sum(nil)
	curve := btcec.S256()
	childKey, _ := btcec.PrivKeyFromBytes(curve, privKey)
	childKeyNum := new(big.Int).SetBytes(childKey.Serialize())
	factor := new(big.Int).SetBytes(payload[:32])
	childKeyNum.Add(childKeyNum, factor)
	childKeyNum.Mod(childKeyNum, curve.Params().N)
	childKey, _ = btcec.PrivKeyFromBytes(curve, childKeyNum.Bytes())
	return childKey.Serialize(), payload[32:]
}

func DerivePublicKey(mnemonic, passphrase, path string) []byte {
	seed := mnemonicToSeed(mnemonic, passphrase)
	privKey, chainCode := derivePrivateKeyFromSeed(seed)
	pubKey := privKeyToPubKey(privKey)
	pathList := getPathList(path)
	for _, index := range pathList {
		hardened := index & HardenedIndex
		privKey, chainCode = derivePrivateChild(pubKey, privKey, chainCode, index, hardened != 0)
		pubKey = privKeyToPubKey(privKey)
	}
	return pubKey
}
