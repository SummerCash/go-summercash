// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"crypto/rand"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestGetPeerIdentity tests the functionality of the GetPeerIdentity helper method.
func TestGetPeerIdentity(t *testing.T) {
	_, err := GetPeerIdentity() // Get peer identity

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestNewPeerIdentity tests the functionality of the NewPeerIdentity helper method.
func TestNewPeerIdentity(t *testing.T) {
	_, err := NewPeerIdentity() // Get peer identity

	if err != nil && err != ErrIdentityAlreadyExists { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestWritePeerIdentity tests the functionality of the WritePeerIdentity helper method.
func TestWritePeerIdentity(t *testing.T) {
	privateKey, _, err := crypto.GenerateRSAKeyPair(2048, rand.Reader) // Generate RSA key pair

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = WritePeerIdentity(&privateKey) // Write identity

	if err != nil && err != ErrIdentityAlreadyExists { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestGetExistingPeerIdentity tests the functionality of the GetExistingPeerIdentity helper method.
func TestGetExistingPeerIdentity(t *testing.T) {
	privateKey, _, err := crypto.GenerateRSAKeyPair(2048, rand.Reader) // Generate RSA key pair

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = WritePeerIdentity(&privateKey) // Write identity

	if err != nil && err != ErrIdentityAlreadyExists { // Check for errors
		t.Fatal(err) // Panic
	}

	_, err = GetExistingPeerIdentity() // Get peer identity

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
