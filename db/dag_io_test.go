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

/* BEGIN EXPORTED METHODS */

// TestWriteToMemory tests the functionality of the WriteToMemory helper method.
func TestWriteToMemory(t *testing.T) {
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

	err = dag.WriteToMemory("test_net") // Write dag to persistent memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestReadFromMemory tests the functionality of the ReadFromMemory() helper method.
func TestReadFromMemory(t *testing.T) {
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
		newTransaction, err := types.NewTransaction(uint64(i+1), lastLeaf.Transaction, &sender, &sender, big.NewFloat(0), []byte("test")) // Initialize transaction
		if err != nil {                                                                                                                   // Check for errors
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

	err = dag.WriteToMemory("test_net") // Write dag to persistent memory
	if err != nil {                     // Check for errors
		t.Fatal(err) // Panic
	}

	readDag, err := ReadDagFromMemory("test_net") // Read dag from persistent memory
	if err != nil {                               // Check for errors
		t.Fatal(err) // Panic
	}

	if !bytes.Equal(readDag.Root.Hash[:], dag.Root.Hash[:]) { // Check invalid roots
		t.Fatal("should have same root") // Panic
	}

	flattenedDag, err := dag.Flatten() // Flatten dag
	if err != nil {                    // Check for errors
		t.Fatal(err) // Panic
	}

	flattenedReadDag, err := readDag.Flatten() // Flatten read dag
	if err != nil {                            // Check for errors
		t.Fatal(err) // Panic
	}

	if len(flattenedDag.Transactions) != len(flattenedReadDag.Transactions) { // Check for errors
		t.Fatalf("expected %d txs in read dag; found %d", len(flattenedDag.Transactions), len(flattenedReadDag.Transactions)) // Panic
	}

	for i, transaction := range flattenedDag.Transactions { // Iterate through txs
		if flattenedReadDag.Transactions[i].String() != transaction.String() { // Check not same tx
			t.Fatal("should be same tx") // Panic
		}
	}
}

/* END EXPORTED METHODS */
