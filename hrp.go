package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
)

const addressSep = "-"

var (
	ErrNoSeparator = errors.New("no separator found in address")
	errBits5To8    = errors.New("unable to convert address from 5-bit to 8-bit formatting")
	errBits8To5    = errors.New("unable to convert address from 8-bit to 5-bit formatting")
)

// Format takes in a chain prefix, HRP, and byte slice to produce a string for an address.
func Format(chainIDAlias, hrp string, addr []byte) (string, error) {
	addrStr, err := FormatBech32(hrp, addr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", chainIDAlias, addressSep, addrStr), nil
}

// FormatBech32 takes an address's bytes as input and returns a bech32 address.
func FormatBech32(hrp string, payload []byte) (string, error) {
	fiveBits, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return "", errBits8To5
	}
	return bech32.Encode(hrp, fiveBits)
}

// Parse takes in an address string and returns the corresponding parts.
func Parse(addrStr string) (string, string, []byte, error) {
	addressParts := strings.SplitN(addrStr, addressSep, 2)
	if len(addressParts) < 2 {
		return "", "", nil, ErrNoSeparator
	}
	chainID := addressParts[0]
	rawAddr := addressParts[1]

	hrp, addr, err := ParseBech32(rawAddr)
	return chainID, hrp, addr, err
}

// ParseBech32 takes a bech32 address as input and returns the HRP and data section of a bech32 address.
func ParseBech32(addrStr string) (string, []byte, error) {
	hrp, decoded, err := bech32.Decode(addrStr)
	if err != nil {
		return "", nil, err
	}
	addrBytes, err := bech32.ConvertBits(decoded, 5, 8, true)
	if err != nil {
		return "", nil, errBits5To8
	}
	return hrp, addrBytes, nil
}

func main() {
	// Example Ethereum address (without the 0x prefix)
	ethereumAddressHex := "EcaBC9480D5a8CdbBb509ac697adeEcB3356dc74"

	// Decode the hex string to bytes
	ethereumAddressBytes, err := hex.DecodeString(ethereumAddressHex)
	if err != nil {
		log.Fatalf("Failed to decode Ethereum address: %v", err)
	}

	// Encode the Ethereum address to Bech32 with a specified HRP
	hrp := "cryft"
	encodedAddress, err := FormatBech32(hrp, ethereumAddressBytes)
	if err != nil {
		log.Fatalf("Failed to encode address: %v", err)
	}

	// Print the resulting Bech32 address
	fmt.Printf("Bech32 address: %s\n", encodedAddress)

	// Attempt to parse the address to check its validity
	chainID, hrpRecovered, _, parseErr := Parse(fmt.Sprintf("chain%s%s", addressSep, encodedAddress))
	if parseErr != nil {
		log.Fatalf("Failed to parse address: %v", parseErr)
	}

	// Confirm that the parsing went correctly
	fmt.Printf("Parsed successfully with HRP: %s\n", hrpRecovered)
}
