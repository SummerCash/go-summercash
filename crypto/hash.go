// Package crypto implements a set of abstractions enabling the creation,
// verification, and validation of ECDSA signatures. Additionally, the
// aforementioned package implements several helper types and methods for
// hashing via Blake2s.
package crypto

import (
	"encoding/hex"
	"golang.org/x/crypto/blake2s"
)

// Hash represents a fixed-length blake2s hash.
type Hash [32]byte

/* BEGIN EXPORTED METHODS */

// Blake2s hashes a given input, b, via blake2s.
func Blake2s(b []byte) Hash {
	return Hash(blake2s.Sum256(b)) // Return hashed
}

// Bytes converts a given fixed-length hash into a dynamically-sized byte
// slice.
func (hash *Hash) Bytes() []byte {
	return hash[:] // Slice
}

// String converts a given hash into a hex-encoded string.
func (hash *Hash) String() string {
	return hex.EncodeToString(hash.Bytes()) // Return hex-encoded
}

/* END EXPORTED METHODS */
