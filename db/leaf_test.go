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

/* BEGIN EXPORTED METHODS HELPERS */

// TestNewLeaf test the functionality of the NewLeaf helper.
func TestNewLeaf(t *testing.T) {
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

	transaction, err := types.NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	leaf, err := NewLeaf(transaction) // Initialize leaf

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if leaf == nil { // Check for nil leaf
		t.Fatal("nil leaf") // panic
	}
}

// TestIsOnlyChild tests the functionality of the IsOnlyChild helper method.
func TestIsOnlyChild(t *testing.T) {
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

	transaction, err := types.NewTransaction(0, nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	leaf, err := NewLeaf(transaction) // Initialize leaf

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if leaf.IsOnlyChild() { // Check is only child
		t.Fatal("should not be only child") // Panic
	}
}

/* END EXPORTED METHODS HELPERS */
