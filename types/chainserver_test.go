package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/space55/summertech-blockchain/common"
)

/* BEGIN EXPORTED METHODS */

// TestHandleReceivedChainRequest - test received chain request handler
func TestHandleReceivedChainRequest(t *testing.T) {
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

	err = chain.WriteToMemory() // Write chain to persistent storage

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	t.Log("success") // Log success
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

/* END INTERNAL METHODS */
