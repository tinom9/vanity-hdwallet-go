package main

import (
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"
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

func checkVanity(address string, currency string, vanity string) bool {
	return strings.HasPrefix(address, currencies.PrefixMap[currency]+vanity)
}

func generateAddress(currency string, words int, passphrase string) (string, string, error) {
	mnemonic, err := generateMnemonic(words)
	if err != nil {
		return "", "", err
	}
	var address string
	if currency == currencies.Cosmos {
		address, _ = addresses.GetCosmosAddress(mnemonic, passphrase)
	} else if currency == currencies.Bitcoin {
		address, _ = addresses.GetBitcoinAddress(mnemonic, passphrase)
	} else {
		return "", "", ErrInvalidCurrency
	}
	return address, mnemonic, nil
}

func generateVanityAddress(currency string, vanity string, words int, passphrase string, numWorkers int) (string, error) {
	resultChan := make(chan string, numWorkers)
	errorChan := make(chan error, numWorkers)
	doneChan := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-doneChan:
					return
				default:
					address, mnemonic, err := generateAddress(currency, words, passphrase)
					if err != nil {
						errorChan <- err
						return
					}
					if checkVanity(address, currency, vanity) {
						fmt.Printf("Address found: %s\nMnemonic: %s\n", address, mnemonic)
						resultChan <- address
						close(doneChan)
						return
					}
				}
			}
		}()
	}
	select {
	case result := <-resultChan:
		wg.Wait()
		return result, nil
	case err := <-errorChan:
		close(doneChan)
		return "", err
	}
}

func getCLIArgs() (*string, *string, *int, *string, *int) {
	currency := flag.String("currency", currencies.Bitcoin, "currency to use")
	vanity := flag.String("vanity", "", "vanity string to use")
	words := flag.Int("words", 12, "number of words to use")
	passphrase := flag.String("passphrase", "", "passphrase to use")
	numWorkers := flag.Int("num-workers", runtime.NumCPU(), "number of workers to use")
	flag.Parse()
	return currency, vanity, words, passphrase, numWorkers
}

func main() {
	currency, vanity, words, passphrase, numWorkers := getCLIArgs()
	_, err := generateVanityAddress(*currency, *vanity, *words, *passphrase, *numWorkers)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
