// Package common implements a set of commonly-used types and helper methods.
package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

/*
	These tests are literally all the same lul
*/

// TestNewAddress tests the functionality of the NewAddress() helper method.
func TestNewAddress(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Fatal(err) // Panic
	}

	convertedAddress := NewAddressFromPublicKey(&privateKey.PublicKey) // Convert to address

	address := NewAddress(convertedAddress.Bytes()) // Derive address from public key

	if address.Bytes() == nil || address.String() == "000000000000000000000000000000" { // Check nil value
		t.Fatal("address is nil") // Panic
	}
}

// TestNewAddressFromPublicKey tests the functionality of the
// NewAddressFromPublicKey() helper method.
func TestNewAddressFromPublicKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Fatal(err) // Panic
	}

	address := NewAddressFromPublicKey(&privateKey.PublicKey) // Convert to address

	if address.Bytes() == nil || address.String() == "000000000000000000000000000000" { // Check nil value
		t.Fatal("address is nil") // Panic
	}
}

// TestAddressBytes tests the functionality of the Address Bytes() helper
// method.
func TestAddressBytes(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Fatal(err) // Panic
	}

	address := NewAddressFromPublicKey(&privateKey.PublicKey) // Convert to address

	if address.Bytes() == nil || address.String() == "000000000000000000000000000000" { // Check nil value
		t.Fatal("address is nil") // Panic
	}
}

// TestAddressString tests the functionality of the Address String() helper
// method.
func TestAddressString(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Fatal(err) // Panic
	}

	address := NewAddressFromPublicKey(&privateKey.PublicKey) // Convert to address

	if address.Bytes() == nil || address.String() == "000000000000000000000000000000" { // Check nil value
		t.Fatal("address is nil") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
