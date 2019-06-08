// Package db defines the standard go-summercash transaction database.
package db

import (
	"errors"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

var (
	// ErrNoOnlyChild represents an error describing a request for an only child where the leaf has multiple or no children.
	ErrNoOnlyChild = errors.New("leaf has no only child")
)

// Leaf defines a leaf in the go-summercash DAG.
type Leaf struct {
	Transaction *types.Transaction `json:"transaction"` // Transaction at leaf

	Parents  []*Leaf `json:"parents"`  // Transaction parents
	Children []*Leaf `json:"children"` // Leaf children

	Hash common.Hash `json:"hash"` // Transaction hash at leaf
}

/* BEGIN EXPORTED METHODS */

/*
	BEGIN HELPER METHODS
*/

// IsOnlyChild determines whether or not
// the current leaf has any siblings.
func (leaf *Leaf) IsOnlyChild() bool {
	numWithLastLeafAsOnlyChild := 0 // Init num parents with only child

	for _, parent := range leaf.Parents { // Iterate through parents
		if len(parent.Children) == 1 { // Check only child
			numWithLastLeafAsOnlyChild++ // Increment
		}
	}

	return numWithLastLeafAsOnlyChild == len(leaf.Parents) // Return is only child
}

// GetOnlyChild returns the only child of the current leaf
// (if applicable).
func (leaf *Leaf) GetOnlyChild() (*Leaf, error) {
	if len(leaf.Children) > 1 || len(leaf.Children) == 0 { // Check no only child
		return &Leaf{}, ErrNoOnlyChild // Return error
	}

	return leaf.Children[0], nil // Return child
}

/*
	END HELPER METHODS
*/

/* END EXPORTED METHODS */
