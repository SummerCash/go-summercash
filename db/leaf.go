// Package db defines the standard go-summercash transaction database.
package db

import (
	"bytes"
	"errors"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
)

var (
	// ErrNoOnlyChild represents an error describing a request for an only child where the leaf has multiple or no children.
	ErrNoOnlyChild = errors.New("leaf has no only child")

	// ErrNoMatchingLeaf represents an error describing a request for a child with a particular hash--that of which
	// is nonexistent.
	ErrNoMatchingLeaf = errors.New("no matching leaf")

	// ErrNilLeafContents represents an error describing a leaf value of nil.
	ErrNilLeafContents = errors.New("nil leaf contents")
)

// Leaf defines a leaf in the go-summercash DAG.
type Leaf struct {
	Transaction *types.Transaction `json:"transaction"` // Transaction at leaf

	Parents  []*Leaf `json:"parents"`  // Transaction parents
	Children []*Leaf `json:"children"` // Leaf children

	Hash common.Hash `json:"hash"` // Transaction hash at leaf
}

/* BEGIN EXPORTED METHODS */

// NewLeaf initializes a new leaf with the given transaction.
func NewLeaf(transaction *types.Transaction) (*Leaf, error) {
	if transaction == nil { // CHeck for nil transaction
		return &Leaf{}, ErrNilLeafContents // Return error
	}

	return &Leaf{
		Transaction: transaction,
	}, nil // Return leaf
}

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

// GetChildByHash queries the leaf's children for a particular hash.
// If the hash does not exist, an error is returned.
func (leaf *Leaf) GetChildByHash(hash common.Hash) (*Leaf, error) {
	if child, err := leaf.GetDirectChildByHash(hash); err == nil { // Check has child
		return child, nil // Return child
	}

	for i, child := range leaf.Children { // Iterate through children
		match, err := child.GetChildByHash(hash) // Get child

		if err != nil && i == len(leaf.Children)-1 { // Check for errors
			return &Leaf{}, err // Return found error
		}

		if match != nil { // Check for match
			return match, nil // Return match
		}
	}

	return &Leaf{}, ErrNoMatchingLeaf // Return error
}

// GetDirectChildByHash searches the set of immediate children for the
// given hash. If no child exists with the corresponding hash, an
// error is returned.
func (leaf *Leaf) GetDirectChildByHash(hash common.Hash) (*Leaf, error) {
	for _, child := range leaf.Children { // Iterate through children
		if bytes.Equal(child.Hash[:], hash[:]) { // Check for match
			return child, nil // Return child
		}
	}

	return &Leaf{}, ErrNoMatchingLeaf // Return error
}

// GetNextCommonLeaf implements the functionality of the DAG
// QueryNextCommonLeaf method.
func (leaf *Leaf) GetNextCommonLeaf() (*Leaf, error) {
	if child, err := leaf.GetOnlyChild(); err == nil { // Check has child
		return child, nil // Return child
	}

	for i, child := range leaf.Children { // Iterate through children
		commonLeaf, err := child.GetNextCommonLeaf() // Get next common leaf

		if err != nil && i == len(leaf.Children)-1 { // Check for errors
			return &Leaf{}, err // Return found error
		}

		if commonLeaf != nil { // Check has leaf
			return commonLeaf, nil // Return common leaf
		}
	}

	return &Leaf{}, ErrNoCommonLeaf // Return error
}

/*
	END HELPER METHODS
*/

/* END EXPORTED METHODS */
