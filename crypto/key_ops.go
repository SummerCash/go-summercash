// Package crypto implements a set of abstractions enabling the creation,
// verification, and validation of ECDSA signatures. Additionally, the
// aforementioned package implements several helper types and methods for
// hashing via Blake2s.
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/SummerCash/go-summercash/common"
)

/* BEGIN EXPORTED METHODS */

// MarshalPublicKey marshals a given public key to a byte slice via
// elliptic.Marshal. Note that this method assumes that the P521
// curve is being used.
func MarshalPublicKey(publicKey *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(common.PreferredCurve(), publicKey.X, publicKey.Y) // Return marshaled
}

// UnmarshalPublicKey unmarshals a given byte slice, b, into an ecdsa public
// key. Note that this method assumes that the P521 curve is being used.
func UnmarshalPublicKey(b []byte) *ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(common.PreferredCurve(), b) // Get x, y from data

	return &ecdsa.PublicKey{
		Curve: common.PreferredCurve(), // Use the commonly-preferred curve
		X:     x,                       // Set x
		Y:     y,                       // Set y
	} // Return unmarshalled public key
}

/* END EXPORTED METHODS */
