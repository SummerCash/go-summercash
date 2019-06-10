// Package db defines the standard go-summercash transaction database.
package db

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewDag tests the functionality of the NewDag helper method.
func TestNewDag(t *testing.T) {
	dag := NewDag() // Initialize dag

	if dag == nil { // Check nil pointer
		t.Fatal("nil pointer") // Panic
	}
}

// TestNewDagWithRoot tests the functionality of the NewDagWithRoot helper method.
func TestNewDagWithRoot(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key
	if err != nil {                                                    // Check for errors
		t.Fatal(err) // Panic
	}

	sender, err := common.NewAddress(privateKey) // Initialize address from private key
	if err != nil {                              // Check for errors
		t.Fatal(err) // Panic
	}

	transaction, err := types.NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
	if err != nil {                                                                                     // Check for errors
		t.Fatal(err) // Panic
	}

	root, err := NewLeaf(transaction) // Initialize leaf
	if err != nil {                   // Check for errors
		t.Fatal(err) // Panic
	}

	dag := NewDagWithRoot(root) // Initialize dag

	if dag.Root != root { // Check not same root
		t.Fatal("dag should have same root") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
