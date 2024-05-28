package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"
)

const addressSep = "-"

var (
	ErrNoSeparator = errors.New("no separator found in address")
	errBits5To8    = errors.New("unable to convert address from 5-bit to 8-bit formatting")
	errBits8To5    = errors.New("unable to convert address from 8-bit to 5-bit formatting")
)

// Parse takes in an address string and returns the corresponding parts.
// This returns the chain ID alias, bech32 HRP, address bytes, and an error if it occurs.
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

// Format takes in a chain prefix, HRP, and byte slice to produce a string for an address.
func Format(chainIDAlias string, hrp string, addr []byte) (string, error) {
	addrStr, err := FormatBech32(hrp, addr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", chainIDAlias, addressSep, addrStr), nil
}

// ParseBech32 takes a bech32 address as input and returns the HRP and data section of a bech32 address.
func ParseBech32(addrStr string) (string, []byte, error) {
	rawHRP, decoded, err := bech32.Decode(addrStr)
	if err != nil {
		return "", nil, err
	}
	addrBytes, err := bech32.ConvertBits(decoded, 5, 8, true)
	if err != nil {
		return "", nil, errBits5To8
	}
	return rawHRP, addrBytes, nil
}

// FormatBech32 takes an address's bytes as input and returns a bech32 address.
func FormatBech32(hrp string, payload []byte) (string, error) {
	fiveBits, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return "", errBits8To5
	}
	return bech32.Encode(hrp, fiveBits)
}

func main() {
	// Example private key from MetaMask (in hex format without 0x)
	privateKeyHex := "PRIVATE_KEY_HERE"

	// Decode the private key from hex
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// Derive the public key from the private key
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	// Convert the public key to the Ethereum address format
	pubBytes := crypto.FromECDSAPub(publicKey)
	hash := sha256.Sum256(pubBytes[1:])
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(hash[:])
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// Convert the address to 5-bit groups
	data, err := bech32.ConvertBits(publicKeyHash, 8, 5, true)
	if err != nil {
		log.Fatalf("Failed to convert bits: %v", err)
	}

	// Encode the address using Bech32 with the specified HRP
	hrp := "cryft"
	encodedAddress, err := bech32.Encode(hrp, data)
	if err != nil {
		log.Fatalf("Failed to encode address: %v", err)
	}

	// Format the address with chain ID alias
	chainIDAlias := "chain"
	formattedAddress, err := Format(chainIDAlias, hrp, publicKeyHash)
	if err != nil {
		log.Fatalf("Failed to format address: %v", err)
	}

	// Print the resulting Bech32 address
	fmt.Printf("Bech32 address: %s\n", encodedAddress)
	fmt.Printf("Formatted address: %s\n", formattedAddress)

	// Parse the formatted address to check for separator
	_, _, _, err = Parse(formattedAddress)
	if err != nil {
		if err == ErrNoSeparator {
			fmt.Println("The formatted address does not contain a separator.")
		} else {
			log.Fatalf("Failed to parse address: %v", err)
		}
	} else {
		fmt.Println("The address is successfully derived from the private key!")
	}
}
