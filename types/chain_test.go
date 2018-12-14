package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/space55/summertech-blockchain/common"
)

// TestNewChain - test chain initializer
func TestNewChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(*chain) // Log initialized chain
}

// TestAddTransaction - test functionality of transaction appending
func TestAddTransaction(t *testing.T) {
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

	transaction, err := NewTransaction(0, nil, &sender, &sender, 0, []byte("test")) // Initialize transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created transaction: %s", transaction.Hash.String()) // Log issued tx

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("signed transaction: %s", transaction.Signature.String()) // Log signed

	chain, err := NewChain(sender) // Initialize chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("created chain: %s", chain.ID.String()) // Log init

	err = chain.AddTransaction(transaction) // Add transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Logf("added transaction: %s", transaction.Hash.String()) // Log signed

	t.Log("success") // Log success
}

// TestBytesChain - test chain to bytes conversion
func TestBytesChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	byteVal := chain.Bytes() // Get byte val

	if byteVal == nil { // Check nil byte val
		t.Errorf("invalid byte val") // Log error
		t.FailNow()                  // Panic
	}

	t.Log(byteVal) // Log success
}

// TestStringChain - test chain to string conversion
func TestStringChain(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	address, err := common.NewAddress(privateKey) // Generate address

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err := NewChain(address) // Initialize chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	stringVal := chain.String() // Get string val

	if stringVal == "" { // Check nil string val
		t.Errorf("invalid string val") // Log error
		t.FailNow()                    // Panic
	}

	t.Log(stringVal) // Log success
}
