// Package p2p outlines helper methods and types for p2p communications.
package p2p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestGetPeerIdentity tests the functionality of the GetPeerIdentity helper method.
func TestGetPeerIdentity(t *testing.T) {
	_, err := GetPeerIdentity() // Get peer identity

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestGetLibp2pPeerIdentity tests the functionality of the GetLibp2pPeerIdentity helper method.
func TestGetLibp2pPeerIdentity(t *testing.T) {
	_, err := GetLibp2pPeerIdentity() // Get peer identity

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
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = WritePeerIdentity(privateKey) // Write identity

	if err != nil && err != ErrIdentityAlreadyExists { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestGetExistingPeerIdentity tests the functionality of the GetExistingPeerIdentity helper method.
func TestGetExistingPeerIdentity(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = WritePeerIdentity(privateKey) // Write identity

	if err != nil && err != ErrIdentityAlreadyExists { // Check for errors
		t.Fatal(err) // Panic
	}

	_, err = GetExistingPeerIdentity() // Get peer identity

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
