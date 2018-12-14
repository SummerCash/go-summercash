package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

/* BEGIN EXPORTED METHODS */

// TestWriteToMemory - test address serialization, persistence
func TestWriteToMemory(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = address.WriteToMemory(privateKey) // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log("success") // Log success
}

// TestReadFromMemory - test reading of account from persistent storage
func TestReadFromMemory(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = address.WriteToMemory(privateKey) // Write to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	readAddress, privateKey, err := ReadAccountFromMemory(&address) // Read account from memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(privateKey)           // Log success
	t.Log(readAddress.String()) // Log success
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

/* END INTERNAL METHODS */
