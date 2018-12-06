package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
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

// TestWitnessTransaction - test functionality of WitnessTransaction() method
func TestWitnessTransaction(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	witness, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generat private key

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

	err = WitnessTransaction(transaction, witness) // Witness transaction

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

	sender, err := common.NewAddress(privateKey) // Initialize address from private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	validTransaction, err := NewTransaction(0, &sender, &sender, 0, []byte("test")) // Initialize transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	err = SignTransaction(validTransaction, privateKey) // Sign transaction with valid keypair

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	valid, err := VerifyTransactionSignature(validTransaction) // Verify signature

	if err != nil && err != ErrInvalidSignature { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("signature valid: %t", valid) // Log success
}

// TestBytesSignature - test functionality of signature to bytes extension method
func TestBytesSignature(t *testing.T) {
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

	byteVal := transaction.Signature.Bytes() // Get byte val

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte value") // Log found error
		t.FailNow()                    // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringSignature - test functionality of signature to string extension method
func TestStringSignature(t *testing.T) {
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

	stringVal := transaction.Signature.String() // Get string val

	if stringVal == "" { // Check for nil string val
		t.Errorf("invalid string value") // Log found error
		t.FailNow()                      // Panic
	}

	t.Log(stringVal) // Log success
}

// TestBytesWitness - test functionality of Bytes() witness extension method
func TestBytesWitness(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	witness, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generat private key

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

	err = WitnessTransaction(transaction, witness) // Witness transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := transaction.Witness.Bytes() // Get byte val

	if byteVal == nil { // Check for nil byte val
		t.Errorf("invalid byte value") // Log found error
		t.FailNow()                    // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringWitness - test functionality of String() witness extension method
func TestStringWitness(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	witness, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generat private key

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

	err = WitnessTransaction(transaction, witness) // Witness transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := transaction.Witness.String() // Get string val

	if stringVal == "" { // Check for nil string val
		t.Errorf("invalid string value") // Log found error
		t.FailNow()                      // Panic
	}

	t.Log(stringVal) // Log success
}

// TestSelfSignTransaction - test functionality of self tx-signing
func TestSelfSignTransaction(t *testing.T) {
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

	err = selfSignTransaction(transaction, privateKey) // Self-sign transaction

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

// TestSignTransactionWitness - test functionality of witness tx-signing
func TestSignTransactionWitness(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	witnessPrivateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

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

	err = signTransactionWitness(transaction, witnessPrivateKey) // Witness-sign transaction

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

/* END EXPORTED METHODS */
