package common

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
)

// Address - []byte wrapper for addresses
type Address []byte

// AddressLength - max addr length
const AddressLength = 32

// NewAddress - initialize new address
func NewAddress(privateKey *ecdsa.PrivateKey) (Address, error) {
	marhsaledPublicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey) // Marshal public key

	if err != nil { // Check for errors
		return Address{}, err // Return found error
	}

	return marhsaledPublicKey[0:32], nil // Return public key
}

// StringToAddress - convert string to address
func StringToAddress(s string) (Address, error) {
	decoded, err := hex.DecodeString(s) // Decode string

	if err != nil { // Check for errors
		return Address{}, err // Return found error
	}

	return decoded, nil // Return address
}

// Bytes - convert given address to bytes
func (address Address) Bytes() []byte {
	return address // Return byte val
}

// String - convert given address to string
func (address Address) String() string {
	return hex.EncodeToString([]byte(address[:])) // Return string val
}
