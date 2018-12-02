package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/space55/summertech-blockchain/common"
)

/* BEGIN EXPORTED METHODS */

// TestSignTransaction - test functionality of SignTransaction() method
func TestSignTransaction(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	sender, err := common.NewAddress(privateKey) // Initialize address from private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	transaction, err := NewTransaction(0, &sender, &sender, 0, []byte("test")) // Initialize transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	marshaledVal, err := json.MarshalIndent(*transaction, "", "  ") // Marshal tx

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(string(marshaledVal)) // Log success
}

// TestVerifyTransactionSignature - test functionality of VerifyTransactionSignature() method
func TestVerifyTransactionSignature(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	invalidPrivateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	sender, err := common.NewAddress(privateKey) // Initialize address from private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	sender2, err := common.NewAddress(invalidPrivateKey) // Init address from private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	fmt.Println(sender)
	fmt.Println(sender2)

	transaction, err := NewTransaction(0, &sender, &sender, 0, []byte("test")) // Initialize transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(transaction, invalidPrivateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	valid, err := VerifyTransactionSignature(transaction) // Verify signature

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("signature valid: %t", valid) // Log success
}

/* END EXPORTED METHODS */
