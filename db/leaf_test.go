// Package db defines the standard go-summercash transaction database.
package db

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
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

	root, err := NewLeaf(transaction) // Initialize leaf

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if !root.IsOnlyChild() { // Check is only child
		t.Fatal("should be only child") // Panic
	}

	for i := 0; i < 3; i++ { // Add children
		leaf, err := NewLeaf(transaction) // Initialize leaf

		if err != nil { // Check for errors
			t.Fatal(err) // Panic
		}

		leaf.Parents = append(leaf.Parents, root) // Append root to parents

		root.Children = append(root.Children, leaf) // Append leaf to root children
	}

	if root.Children[0].IsOnlyChild() { // Check is only child
		t.Fatal("should not be only child") // Panic
	}
}

// TestGetOnlyChild tests the functionality of the GetOnlyChild helper method.
func TestGetOnlyChild(t *testing.T) {
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

	root, err := NewLeaf(transaction) // Initialize leaf

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	child, err := NewLeaf(transaction) // Initialize leaf

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	root.Children = append(root.Children, child) // Append child

	child.Parents = append(child.Parents, root) // Append root as parent

	onlyChild, err := root.GetOnlyChild() // Get only child

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if onlyChild != child { // Check not same child
		t.Fatal("only child query should lead to same child") // Panic
	}
}

// TestGetChildByHash tests the functionality of the GetChildByHash helper method.
func TestGetChildByHash(t *testing.T) {
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

	root, err := NewLeaf(transaction) // Initialize leaf

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	for i := 0; i < 3; i++ { // Add children
		transaction, err := types.NewTransaction(uint64(i+1), nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction

		if err != nil { // Check for errors
			t.Error(err) // Log found error
			t.FailNow()  // Panic
		}

		child, err := NewLeaf(transaction) // Initialize leaf

		if err != nil { // Check for errors
			t.Fatal(err) // Panic
		}

		child.Parents = append(child.Parents, root) // Append root as parent

		root.Children = append(root.Children, child) // Append child to root

		foundChild, err := root.GetChildByHash(child.Hash) // Get child by hash

		if err != nil { // Check for errors
			t.Fatal(err) // Panic
		}

		if foundChild != child { // Check not same child
			fmt.Println(*foundChild)
			fmt.Println(*child)
			t.Fatal("child by hash query should lead to same child") // Panic
		}
	}
}

/* END EXPORTED METHODS HELPERS */
