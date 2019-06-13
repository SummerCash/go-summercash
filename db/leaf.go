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
func NewLeaf(transaction *types.Transaction, parent *Leaf) (*Leaf, error) {
	if transaction == nil || transaction.Hash == nil { // CHeck for nil transaction
		return &Leaf{}, ErrNilLeafContents // Return error
	}

	if parent == nil { // Check no parent
		return &Leaf{
			Transaction: transaction,       // Set transaction
			Hash:        *transaction.Hash, // Set hash
		}, nil // Return leaf
	}

	return &Leaf{
		Transaction: transaction,       // Set transaction
		Hash:        *transaction.Hash, // Set hash
		Parents:     []*Leaf{parent},   // Set parents
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

// GetChildrenBySender queries the leaf's children for a particular sender.
// If a transaction with the corresponding sender does not exist, an error is returned.
func (leaf *Leaf) GetChildrenBySender(sender common.Address) ([]*Leaf, error) {
	children, _ := leaf.GetDirectChildrenBySender(sender) // Get children

	for i, child := range leaf.Children { // Iterate through children
		matches, err := child.GetChildrenBySender(sender) // Get child

		if err != nil && i == len(leaf.Children)-1 && (children == nil || len(children) == 0) { // Check for errors
			return []*Leaf{}, err // Return found error
		}

		if matches != nil && len(matches) != 0 { // Check for match
			children = append(children, matches...) // Append matches
		}
	}

	return children, nil // Return error
}

// GetDirectChildrenBySender searches the set of immediate children for the
// given sender. If no child exists with the corresponding sender, an
// error is returned.
func (leaf *Leaf) GetDirectChildrenBySender(sender common.Address) ([]*Leaf, error) {
	matchingChildren := []*Leaf{} // Init buffer

	for _, child := range leaf.Children { // Iterate through children
		if child.Transaction.Sender != nil { // Check has sender
			if bytes.Equal(child.Transaction.Sender[:], sender[:]) { // Check for match
				matchingChildren = append(matchingChildren, child) // Append child
			}
		}
	}

	if len(matchingChildren) != 0 { // Check has children
		return matchingChildren, nil // Return matches
	}

	return []*Leaf{}, ErrNoMatchingLeaf // Return error
}

// GetChildrenByRecipient queries the leaf's children for a particular recipient.
// If a transaction with the corresponding recipient does not exist, an error is returned.
func (leaf *Leaf) GetChildrenByRecipient(recipient common.Address) ([]*Leaf, error) {
	children, _ := leaf.GetDirectChildrenByRecipient(recipient) // Get children

	for i, child := range leaf.Children { // Iterate through children
		matches, err := child.GetChildrenByRecipient(recipient) // Get child

		if err != nil && i == len(leaf.Children)-1 && (children == nil || len(children) == 0) { // Check for errors
			return []*Leaf{}, err // Return found error
		}

		if matches != nil && len(matches) != 0 { // Check for match
			children = append(children, matches...) // Append matches
		}
	}

	return children, nil // Return error
}

// GetDirectChildrenByRecipient searches the set of immediate children for the
// given recipient. If no child exists with the corresponding recipient, an
// error is returned.
func (leaf *Leaf) GetDirectChildrenByRecipient(recipient common.Address) ([]*Leaf, error) {
	matchingChildren := []*Leaf{} // Init buffer

	for _, child := range leaf.Children { // Iterate through children
		if child.Transaction.Sender != nil { // Check has sender
			if bytes.Equal(child.Transaction.Recipient[:], recipient[:]) { // Check for match
				matchingChildren = append(matchingChildren, child) // Append child
			}
		}
	}

	if len(matchingChildren) != 0 { // Check has children
		return matchingChildren, nil // Return matches
	}

	return []*Leaf{}, ErrNoMatchingLeaf // Return error
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

// GetChildren gets all children of the given leaf.
func (leaf *Leaf) GetChildren() ([]*Leaf, error) {
	children := leaf.Children // Get direct children

	for _, child := range children { // Iterate through children
		children, err := child.GetChildren() // Get children

		if err != nil { // Check for errors
			return []*Leaf{}, err // Return found error
		}

		children = append(children, children...) // Append children
	}

	return children, nil // Return children
}

/*
	END HELPER METHODS
*/

/* END EXPORTED METHODS */
