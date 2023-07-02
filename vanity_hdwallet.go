package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/tyler-smith/go-bip39"
)

const (
	BitcoinCurrency = "bitcoin"
	CosmosCurrency  = "cosmos"
)

var (
	currencyPrefixMap = map[string]string{
		BitcoinCurrency: "bc1q",
		CosmosCurrency:  "cosmos1",
	}
	currencyPathMap = map[string]string{
		BitcoinCurrency: "m/84'/0'/0'/0/0",
		CosmosCurrency:  "m/44'/118'/0'/0/0",
	}
)

var (
	ErrInvalidWordCount = errors.New("word count should be a multiple of 3 and between 12 and 24")
	ErrInvalidCurrency  = errors.New("currency should be one of: bitcoin, cosmos")
)

func getBech32Address(pk []byte, hrp string, version *int) (string, error) {
	s := sha3(pk)
	r := ripemd160Hash(s)
	fiveBitR, _ := bech32.ConvertBits(r, 8, 5, true)
	var data []byte
	if version != nil {
		data = append(data, byte(*version))
	}
	data = append(data, fiveBitR...)
	return bech32.Encode(hrp, data)
}

func getBitcoinAddress(mnemonic, passphrase string) (string, error) {
	pk := derivePublicKey(mnemonic, passphrase, currencyPathMap[BitcoinCurrency])
	hrp := "bc"
	version := 0
	return getBech32Address(pk, hrp, &version)
}

func getCosmosAddress(mnemonic, passphrase string) (string, error) {
	pk := derivePublicKey(mnemonic, passphrase, currencyPathMap[CosmosCurrency])
	hrp := "cosmos"
	return getBech32Address(pk, hrp, nil)
}

func generateMnemonic(words int) (string, error) {
	if words%3 != 0 || words < 12 || words > 24 {
		return "", ErrInvalidWordCount
	}
	entropy, _ := bip39.NewEntropy(words * 32 / 3)
	return bip39.NewMnemonic(entropy)
}

func checkVanity(address, prefix, vanity string) bool {
	return strings.HasPrefix(address, prefix+vanity)
}

func generateVanityAddress(currency, vanity string, words int, passphrase string) (string, error) {
	for count := 1; ; count++ {
		fmt.Printf("Try: %d\n", count)
		m, _ := generateMnemonic(12)
		var address string
		if currency == CosmosCurrency {
			address, _ = getCosmosAddress(m, passphrase)
		} else if currency == BitcoinCurrency {
			address, _ = getBitcoinAddress(m, passphrase)
		} else {
			return "", ErrInvalidCurrency
		}
		if checkVanity(address, currencyPrefixMap[currency], vanity) {
			fmt.Printf("Address found: %s\nMnemonic: %s\n", address, m)
			return address, nil
		}
	}
}

func main() {
	_, _ = generateVanityAddress(BitcoinCurrency, "", 24, "")
}
