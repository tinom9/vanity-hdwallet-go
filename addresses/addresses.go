package addresses

import (
	"vanityhdwallet/crypto"
	"vanityhdwallet/currency"
	"vanityhdwallet/derivations"

	"github.com/btcsuite/btcutil/bech32"
)

func getBech32Address(pk []byte, hrp string, version *int) (string, error) {
	s := crypto.Sha3(pk)
	r := crypto.Ripemd160Hash(s)
	fiveBitR, _ := bech32.ConvertBits(r, 8, 5, true)
	var data []byte
	if version != nil {
		data = append(data, byte(*version))
	}
	data = append(data, fiveBitR...)
	return bech32.Encode(hrp, data)
}

func GetBitcoinAddress(mnemonic, passphrase string) (string, error) {
	pk := derivations.DerivePublicKey(mnemonic, passphrase, currency.PathMap[currency.Bitcoin])
	hrp := "bc"
	version := 0
	return getBech32Address(pk, hrp, &version)
}

func GetCosmosAddress(mnemonic, passphrase string) (string, error) {
	pk := derivations.DerivePublicKey(mnemonic, passphrase, currency.PathMap[currency.Cosmos])
	hrp := "cosmos"
	return getBech32Address(pk, hrp, nil)
}
