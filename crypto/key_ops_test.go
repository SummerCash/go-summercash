// Package crypto implements a set of abstractions enabling the creation,
// verification, and validation of ECDSA signatures. Additionally, the
// aforementioned package implements several helper types and methods for
// hashing via Blake2s.
package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestMarshalPublicKey tests the functionality of the MarshalPublicKey()
// helper method.
func TestMarshalPublicKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(PreferredCurve(), rand.Reader) // Initialize private key
	if err != nil {                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	if MarshalPublicKey(&privateKey.PublicKey) == nil || len(MarshalPublicKey(&privateKey.PublicKey)) == 0 { // Check couldn't encoode
		t.Fatal("invalid public key marshalling") // Panic
	}
}

// TestUnmarshalPublicKey tests the functionality of the UnmarshalPublicKey()
// helper method.
func TestUnmarshalPublicKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(PreferredCurve(), rand.Reader) // Initialize private key
	if err != nil {                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	publicKey := UnmarshalPublicKey(MarshalPublicKey(&privateKey.PublicKey)) // Unmarshal

	if publicKey.X.String() != privateKey.PublicKey.X.String() || publicKey.Y.String() != privateKey.PublicKey.Y.String() { // Check either x or y not equal
		t.Fatal("invalid unmarshalled public key") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
