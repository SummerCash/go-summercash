// Package db defines the standard go-summercash transaction database.
package db

import (
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

// NewDag initializes a new dag.
func NewDag() *Dag {
	return &Dag{} // Return dag
}

// NewDagWithRoot initializes a new dag from a given root.
func NewDagWithRoot(root *Leaf) *Dag {
	return &Dag{
		Root: root, // Set root
	} // Return dag
}

/*
	BEGIN HELPER METHODS
*/

// BEGIN LEAF HELPERS

// AddLeaf adds a given leaf to the working dag
func (dag *Dag) AddLeaf(leaf *Leaf) error {

}

// END LEAF HELPERS

// BEGIN QUERY HELPERS

// QueryTransactionWithHash queries the dag for a transaction with the corresponding hash.
func (dag *Dag) QueryTransactionWithHash(hash common.Hash) (*types.Transaction, error) {
	leaf, err := dag.Root.GetChildByHash(hash) // Get child

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	return leaf.Transaction, nil // Return transaction
}

// QueryNextCommonLeaf attempts to find the next common leaf in the dag.
// A common leaf is defined as a leaf that has no siblings.
// If no common leaf exists, an error is returned.
func (dag *Dag) QueryNextCommonLeaf(lastCommonLeaf *Leaf) (*Leaf, error) {
	return lastCommonLeaf.GetNextCommonLeaf() // Get next common leaf
}

// END QUERY HELPERS

/*
	END HELPER METHODS
*/

/* END EXPORTED METHODS */
