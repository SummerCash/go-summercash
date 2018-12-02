package types

import (
	"crypto/ecdsa"
	"math/big"
)

// signature - struct containing signature values
type signature struct {
	PublicKey *ecdsa.PublicKey // Public key

	V []byte   // Hash signature value
	R *big.Int // Signature R
	S *big.Int // Signature S
}
