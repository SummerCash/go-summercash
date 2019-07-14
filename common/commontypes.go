package common

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/SummerCash/go-summercash/crypto"
)

// Address - []byte wrapper for addresses
type Address [AddressLength]byte

// Hash - []byte wrapper for hashes
type Hash [HashLength]byte

// AddressSpace - given address-space (e.g. 0x000-0x123)
type AddressSpace struct {
	Addresses []Address `json:"addresses"` // Addresses in address space
	ID        Hash      `json:"ID"`        // AddressSpace identifier
}

// ErrDuplicateAddress - error definition describing two addresses of equal value
var ErrDuplicateAddress = errors.New("duplicate address")

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

	marshaledPublicKey := elliptic.Marshal(privateKey.PublicKey, privateKey.PublicKey.X, privateKey.PublicKey.Y) // Marshal public key

	marshaledPublicKey = append([]byte("0x"), marshaledPublicKey...) // Prepend prefix

	copy(address[:], marshaledPublicKey) // Copy marshaled

	return address, nil // Return public key
}

// PublicKeyToAddress - initialize new address with given public key
func PublicKeyToAddress(publicKey *ecdsa.PublicKey) Address {
	var address Address // Init buffer

	marshaledPublicKey := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y) // Marshal public key

	marshaledPublicKey = append([]byte("0x"), marshaledPublicKey...) // Prepend prefix

	copy(address[:], marshaledPublicKey) // Copy marshaled

	return address // Return public key
}

// StringToAddress - convert string to address
func StringToAddress(s string) (Address, error) {
	var address Address // Init buffer

	decoded, err := hex.DecodeString(s[2:]) // Decode string
	if err != nil {                         // Check for errors
		return Address{}, err // Return found error
	}

	copy(address[:], append([]byte("0x"), decoded...)) // Copy decoded

	return address, nil // Return address
}

// Bytes - convert given address to bytes
func (address Address) Bytes() []byte {
	return address[:] // Return byte val
}

// String - convert given address to string
func (address Address) String() string {
	noPrefix := address[2:] // Remove duplicate 0x prefix

	enc := make([]byte, len(noPrefix)*2+2) // Init encoder buffer

	copy(enc, "0x") // Copy prefix

	hex.Encode(enc[2:], noPrefix[:]) // Encode given byte array

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
	if err != nil {                         // Check for errors
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

/*
	BEGIN ADDRESS-SPACE METHODS
*/

// NewAddressSpace - initialize address space
func NewAddressSpace(originAddresses []Address) (*AddressSpace, error) {
	for x, address := range originAddresses { // Iterate through addresses
		if x != 0 && address == originAddresses[x-1] { // Check for duplicate
			return &AddressSpace{}, ErrDuplicateAddress // Return error
		}
	}

	addressSpace := &AddressSpace{ // Init space
		Addresses: originAddresses,
	}

	hash := crypto.Sha3(addressSpace.Bytes()) // Hash

	(*addressSpace).ID = NewHash(hash) // Set ID

	return addressSpace, nil // Return initialized
}

// Bytes - convert given address-space to bytes
func (addressSpace *AddressSpace) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*addressSpace) // Serialize space

	return buffer.Bytes() // Return serialized
}

// String - convert given address-space to string
func (addressSpace *AddressSpace) String() string {
	marshaled, _ := json.Marshal(*addressSpace) // Marshal address-space

	return string(marshaled) // Return marshaled
}

/*
	END ADDRESS-SPACE METHODS
*/

/* END EXPORTED METHODS */
