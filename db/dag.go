// Package db defines the standard go-summercash transaction database.
package db

import (
	"bytes"
	"errors"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

var (
	// ErrNoCommonLeaf defines an error describing a lack of common leaves.
	ErrNoCommonLeaf = errors.New("no common leaf exists")
)

// Dag implements the standard directed acyclic
// graph global chain.
type Dag struct {
	Root *Leaf `json:"root"` // Root leaf
}

/* BEGIN EXPORTED METHODS */

/*
	BEGIN HELPER METHODS
*/

// QueryTransactionWithHash queries the dag for a transaction with the corresponding hash.
func (dag *Dag) QueryTransactionWithHash(hash common.Hash) (*types.Transaction, error) {
	lastLeaf := dag.Root // Get first leaf (root leaf)

	for { // Wait for match
		if bytes.Equal(lastLeaf.Hash[:], hash[:]) { // Check matching hash
			return lastLeaf.Transaction, nil // Return match
		}

		for _, child := range lastLeaf.Children { // Iterate through children
			if bytes.Equal(child.Hash[:], hash[:]) { // Check matching hash
				return child.Transaction, nil // Return match
			}
		}
	}
}

// GetNextCommonLeaf attempts to find the next common leaf in the dag.
// A common leaf is defined as a leaf that has no siblings.
// If no common leaf exists, an error is returned.
func (dag *Dag) GetNextCommonLeaf(lastCommonLeaf *Leaf) (*Leaf, error) {
	return lastCommonLeaf.GetNextCommonLeaf() // Get next common leaf
}

/*
	END HELPER METHODS
*/

/* END EXPORTED METHODS */
