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

// String - convert given address to string
func (address Address) String() string {
	return hex.EncodeToString([]byte(address[:])) // Return string val
}
