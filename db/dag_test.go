// Package db defines the standard go-summercash transaction database.
package db

import (
	"bytes"
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

	root, err := NewLeaf(transaction, nil) // Initialize leaf
	if err != nil {                        // Check for errors
		t.Fatal(err) // Panic
	}

	dag := NewDagWithRoot(root) // Initialize dag

	if dag.Root != root { // Check not same root
		t.Fatal("dag should have same root") // Panic
	}
}

// TestAddLeaf tests the functionality of the AddLeaf helper method.
func TestAddLeaf(t *testing.T) {
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

	root, err := NewLeaf(transaction, nil) // Initialize leaf
	if err != nil {                        // Check for errors
		t.Fatal(err) // Panic
	}

	dag := NewDagWithRoot(root) // Initialize dag

	if dag.Root != root { // Check not same root
		t.Fatal("dag should have same root") // Panic
	}

	lastLeaf := root // Set last leaf

	for i := 0; i < 1000; i++ { // Lol
		newTransaction, err := types.NewTransaction(uint64(i+1), nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
		if err != nil {                                                                                                  // Check for errors
			t.Fatal(err) // Panic
		}

		leaf, err := NewLeaf(newTransaction, lastLeaf) // Initialize leaf
		if err != nil {                                // Check for errors
			t.Fatal(err) // Panic
		}

		lastLeaf = leaf // Set last leaf

		err = dag.AddLeaf(leaf) // Add leaf to dag
		if err != nil {         // Check for errors
			t.Fatal(err) // Panic
		}
	}

	foundLastLeaf, err := dag.QueryLeafWithHash(lastLeaf.Hash) // Query last leaf hash
	if err != nil {                                            // Check for errors
		t.Fatal(err) // Panic
	}

	if foundLastLeaf != lastLeaf { // Check not same leaf
		t.Fatal("found last leaf should be last leaf") // Panic
	}
}

// TestQueryTransactionsWithSender tests the functionality of the QueryTransactionsWithSender helper method.
func TestQueryTransactionsWithSender(t *testing.T) {
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

	root, err := NewLeaf(transaction, nil) // Initialize leaf
	if err != nil {                        // Check for errors
		t.Fatal(err) // Panic
	}

	dag := NewDagWithRoot(root) // Initialize dag

	if dag.Root != root { // Check not same root
		t.Fatal("dag should have same root") // Panic
	}

	lastLeaf := root // Set last leaf

	for i := 0; i < 1000; i++ { // Lol
		newTransaction, err := types.NewTransaction(uint64(i+1), nil, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
		if err != nil {                                                                                                  // Check for errors
			t.Fatal(err) // Panic
		}

		leaf, err := NewLeaf(newTransaction, lastLeaf) // Initialize leaf
		if err != nil {                                // Check for errors
			t.Fatal(err) // Panic
		}

		lastLeaf = leaf // Set last leaf

		err = dag.AddLeaf(leaf) // Add leaf to dag
		if err != nil {         // Check for errors
			t.Fatal(err) // Panic
		}
	}

	foundTransactions, err := dag.QueryTransactionsWithSender(sender) // Query transactions

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if len(foundTransactions) == 0 { // Check no transactions
		t.Fatal("zero results in dag") // Panic
	}

	lastLeaf = root // Set last leaf

	for i := 0; i < 1001; i++ { // Iterate through txs
		found := false // Init found buffer

		for _, tx := range foundTransactions { // Iterate through queried txs
			if bytes.Equal(tx.Hash[:], lastLeaf.Hash[:]) { // Check match
				found = true // Set found

				break // Break
			}
		}

		if !found { // Check not found
			t.Fatalf("leaf with sender %s not found in dag", lastLeaf.Transaction.Sender.String()) // Panic
		}

		if i != 1000 { // Check is not last tx
			lastLeaf = lastLeaf.Children[0] // Set last leaf
		}
	}
}

/* END EXPORTED METHODS TESTS */
