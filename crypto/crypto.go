package crypto

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

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
