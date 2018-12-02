package common

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/space55/summertech-blockchain/crypto"
)

// Address - []byte wrapper for addresses
type Address [AddressLength]byte

// Hash - []byte wrapper for hashes
type Hash [HashLength]byte

const (
	// AddressLength - max addr length
	AddressLength = 20

	// HashLength - max hash length
	HashLength = 32
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN ADDRESS METHODS
*/

// NewAddress - initialize new address
func NewAddress(privateKey *ecdsa.PrivateKey) (Address, error) {
	var address Address // Init buffer

	marhsaledPublicKey, err := x509.MarshalPKIXPublicKey(&(*privateKey).PublicKey) // Marshal public key

	if err != nil { // Check for errors
		return Address{}, err // Return found error
	}

	fmt.Println(marhsaledPublicKey)

	marhsaledPublicKey = append([]byte("0x"), marhsaledPublicKey...) // Prepend prefix

	copy(address[:], marhsaledPublicKey) // Copy marshaled

	return address, nil // Return public key
}

// PublicKeyToAddress - initialize new address with given public key
func PublicKeyToAddress(publicKey *ecdsa.PublicKey) (Address, error) {
	var address Address // Init buffer

	marhsaledPublicKey, err := x509.MarshalPKIXPublicKey(publicKey) // Marshal public key

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

/*
	END ADDRESS METHODS
*/

/*
	BEGIN HASH METHODS
*/

// NewHash - initialize hash from byte array
func NewHash(b []byte) Hash {
	var hash Hash // Init buffer

	if string(crypto.Sha3(b[0:2])) != string(crypto.Sha3(MemPrefix)) { // Check no mem prefix
		b = append([]byte("0x"), b...) // Prepend 0x
	}

	copy(hash[:], b[:HashLength]) // Copy byte val, trim to max hash length

	return hash // Return init hash
}

// StringToHash - convert string to hash
func StringToHash(s string) (Hash, error) {
	var hash Hash // Init buffer

	decoded, err := hex.DecodeString(s[2:]) // Decode string

	if err != nil { // Check for errors
		return Hash{}, err // Return found error
	}

	copy(hash[:], decoded) // Copy decoded

	return hash, nil // Return hash
}

// Bytes - convert given hash to bytes
func (hash Hash) Bytes() []byte {
	return hash[:] // Return bytes representation
}

// String - convert given hash to string
func (hash Hash) String() string {
	enc := make([]byte, len(hash)*2+2) // Init encoder buffer

	copy(enc, "0x") // Copy prefix

	hex.Encode(enc[2:], hash[:]) // Encode given byte array

	return string(enc) // Return string val
}

/*
	END HASH METHODS
*/

/* END EXPORTED METHODS */
