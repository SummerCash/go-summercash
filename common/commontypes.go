package common

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
)

// Address - []byte wrapper for addresses
type Address [AddressLength]byte

// AddressLength - max addr length
const AddressLength = 20

/* BEGIN EXPORTED METHODS */

// NewAddress - initialize new address
func NewAddress(privateKey *ecdsa.PrivateKey) (Address, error) {
	var address Address // Init buffer

	marhsaledPublicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey) // Marshal public key

	if err != nil { // Check for errors
		return Address{}, err // Return found error
	}

	marhsaledPublicKey = append([]byte("0x"), marhsaledPublicKey...) // Prepend prefix

	copy(address[:], marhsaledPublicKey) // Copy marshaled

	return address, nil // Return public key
}

// StringToAddress - convert string to address
func StringToAddress(s string) (Address, error) {
	var address Address // Init buffer

	decoded, err := hex.DecodeString(s[2:]) // Decode string

	if err != nil { // Check for errors
		return Address{}, err // Return found error
	}

	copy(address[:], decoded) // Copy decoded

	return address, nil // Return address
}

// Bytes - convert given address to bytes
func (address Address) Bytes() []byte {
	return address[:] // Return byte val
}

// String - convert given address to string
func (address Address) String() string {
	enc := make([]byte, len(address)*2+2) // Init encoder buffer

	copy(enc, "0x") // Copy prefix

	hex.Encode(enc[2:], address[:]) // Encode given byte array

	return string(enc) // Return string val
}

/* END EXPORTED METHODS */
