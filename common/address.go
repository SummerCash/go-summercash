// Package common implements a set of commonly-used types and helper methods.
package common

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/SummerCash/go-summercash/crypto"
)

// Address represents a fixed-length blake2s hash of an account's public key.
type Address [32]byte

/* BEGIN EXPORTED METHODS */

// NewAddress converts a given dynamically-sized slice, b, into a fixed-length
// address.
func NewAddress(b []byte) Address {
	var buffer Address // Initialize buffer

	copy(buffer[:], b) // Copy value into buffer

	return buffer // Return contents of buffer
}

// NewAddressFromPublicKey converts a given ECDSA public key to an address.
func NewAddressFromPublicKey(publicKey *ecdsa.PublicKey) Address {
	marshalledPublicKey := crypto.MarshalPublicKey(publicKey) // Marshal public key

	return Address(crypto.Blake2s(marshalledPublicKey)) // Return hashed
}

// Bytes converts a given fixed-length address into a dynamically-sized byte
// slice.
func (address *Address) Bytes() []byte {
	return address[:] // Slice
}

// String converts a given address into a hex-encoded string.
func (address *Address) String() string {
	return hex.EncodeToString(address.Bytes()) // Return hex-encoded
}

/* END EXPORTED METHODS */
