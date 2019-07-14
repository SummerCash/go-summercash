package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"

	"golang.org/x/crypto/sha3"
)

// Signature - struct containing signature values
type Signature struct {
	PublicKey *ecdsa.PublicKey `json:"pubKey"` // Public key

	Time time.Time `json:"time"` // Signature time

	V []byte   `json:"V"` // Hash signature value
	R *big.Int `json:"R"` // Signature R
	S *big.Int `json:"s"` // Signature S
}

// Bytes - convert given signature to byte array
func (signature *Signature) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*signature) // Serialize signature

	return buffer.Bytes() // Return serialized
}

// String - convert given signature to string
func (signature *Signature) String() string {
	marshaled, _ := json.Marshal(*signature) // Marshal signature

	return string(marshaled) // Return marshaled
}

// Sha3 - hash specified byte array
func Sha3(b []byte) []byte {
	hash := sha3.New256() // Init hasher

	hash.Write(b) // Write

	return hash.Sum(nil) // Return final hash
}

// Sha3String - hash specified byte array to string
func Sha3String(b []byte) string {
	b = Sha3(b) // Hash

	return hex.EncodeToString(b) // Return string
}

// Sha3n - hash specified byte array n times
func Sha3n(b []byte, n uint) []byte {
	for x := uint(0); x != n; x++ { // Hash n times
		b = Sha3(b) // Hash
	}

	return b // Return hashed
}

// Sha3nString - hash specified byte array n times to string
func Sha3nString(b []byte, n uint) string {
	b = Sha3n(b, n) // Hash

	return hex.EncodeToString(b) // Return string
}

// Sha3d - hash specified byte array using sha3d algorithm
func Sha3d(b []byte) []byte {
	return Sha3(Sha3(b)) // Return sha3d result
}

// Sha3dString - hash specified byte array to string using sha3d algorithm
func Sha3dString(b []byte) string {
	b = Sha3d(b) // Hash

	return hex.EncodeToString(b) // Return string
}
