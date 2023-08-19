package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"vanityhdwallet/addresses"
	currencies "vanityhdwallet/currency"

	"github.com/tyler-smith/go-bip39"
)

var (
	ErrInvalidWordCount = errors.New("word count should be a multiple of 3 and between 12 and 24")
	ErrInvalidCurrency  = errors.New("currency should be one of: bitcoin, cosmos")
)

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
		m, err := generateMnemonic(words)
		if err != nil {
			return "", err
		}
		var address string
		if currency == currencies.Cosmos {
			address, _ = addresses.GetCosmosAddress(m, passphrase)
		} else if currency == currencies.Bitcoin {
			address, _ = addresses.GetBitcoinAddress(m, passphrase)
		} else {
			return "", ErrInvalidCurrency
		}
		fmt.Printf("Try: %d\n", count)
		if checkVanity(address, currencies.PrefixMap[currency], vanity) {
			fmt.Printf("Address found: %s\nMnemonic: %s\n", address, m)
			return address, nil
		}
	}
}

func main() {
	currency := flag.String("currency", currencies.Bitcoin, "currency to use")
	vanity := flag.String("vanity", "", "vanity string to use")
	words := flag.Int("words", 12, "number of words to use")
	passphrase := flag.String("passphrase", "", "passphrase to use")
	flag.Parse()
	_, err := generateVanityAddress(*currency, *vanity, *words, *passphrase)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
