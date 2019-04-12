package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/crypto"
)

/* BEGIN EXPORTED METHODS */

// TestWriteChainToMemory - test chain writeToMemory() method
func TestWriteChainToMemory(t *testing.T) {
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

	err = chain.WriteToMemory() // Write chain to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log("success") // Log success
}

// TestReadChainFromMemory - test reading of chain
func TestReadChainFromMemory(t *testing.T) {
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

	chain := &Chain{ // Init chain
		Account:      address,
		Transactions: []*Transaction{},
		NetworkID:    0,
	}

	(*chain).ID = common.NewHash(crypto.Sha3(chain.Bytes())) // Set ID

	transaction, err := NewTransaction(0, nil, &address, &address, big.NewFloat(0), []byte("test")) // Initialize transaction

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

	chain.Transactions = []*Transaction{transaction} // Init arr

	err = chain.WriteToMemory() // Write chain to memory

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	chain, err = ReadChainFromMemory(address) // Read chain

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log(chain.String()) // Log success
}

/* END EXPORTED METHODS */
