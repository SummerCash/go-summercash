// Package db defines the standard go-summercash transaction database.
package db

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/SummerCash/go-summercash/accounts"
	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/types"
)

var (
	// ErrNoCommonLeaf defines an error describing a lack of common leaves.
	ErrNoCommonLeaf = errors.New("no common leaf exists")

	// ErrNoParents defines an error describing a lack of parent leaves.
	ErrNoParents = errors.New("leaf has no parents")

	// ErrDuplicateLeaf defines an error describing a duplicate dag leaf.
	ErrDuplicateLeaf = errors.New("leaf already exists in dag")
)

// Dag implements the standard directed acyclic
// graph global chain.
type Dag struct {
	Root *Leaf `json:"root"` // Root leaf
}

// Flattened implements a flattened representation of the
// standard dag.
type Flattened struct {
	Transactions []*types.Transaction `json:"transactions"` // Flattened dag transactions
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
	if dag.Root == nil && len(leaf.Parents) == 0 { // Check no root
		dag.Root = leaf // Set root

		return nil // Return
	}

	if len(leaf.Parents) == 0 { // Check no parents
		return ErrNoParents // Return error
	}

	for _, parent := range leaf.Parents { // Iterate through leaf parents
		for _, child := range parent.Children { // Iterate through children
			if bytes.Equal(child.Hash[:], leaf.Hash[:]) { // Check equal hashes
				return ErrDuplicateLeaf // Return error
			}
		}

		parent.Children = append(parent.Children, leaf) // Append child
	}

	return nil // Return
}

// AddTransaction adds the given transaction to the working dag.
func (dag *Dag) AddTransaction(transaction *types.Transaction) error {
	if (transaction.ParentTx == nil || bytes.Equal(transaction.ParentTx[:], new(common.Hash)[:])) && dag.Root == nil { // Check no root
		leaf, err := NewLeaf(transaction, nil) // Initialize leaf
		if err != nil {                        // Check for errors
			return err // Return found error
		}

		dag.Root = leaf // Set dag root

		return nil // Return
	}

	if transaction.ParentTx == nil { // Check no parent
		return ErrNoParents // Return error
	}

	parent, err := dag.QueryLeafWithHash(*transaction.ParentTx) // Query parent
	if err != nil {                                             // Check for errors
		return err // Return found error
	}

	leaf, err := NewLeaf(transaction, parent) // Initialize leaf
	if err != nil {                           // Check for errors
		return err // Return found error
	}

	parent.Children = append(parent.Children, leaf) // Append leaf as child

	return nil // Return nil
}

// END LEAF HELPERS

// BEGIN QUERY HELPERS

// QueryTransactionWithHash queries the dag for a transaction with the corresponding hash.
func (dag *Dag) QueryTransactionWithHash(hash common.Hash) (*types.Transaction, error) {
	if bytes.Equal(dag.Root.Hash[:], hash[:]) { // Check root is match
		return dag.Root.Transaction, nil // Return root
	}

	leaf, err := dag.Root.GetChildByHash(hash) // Get child
	if err != nil {                            // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	return leaf.Transaction, nil // Return transaction
}

// QueryTransactionsWithSender queries the dag for a list of transactions with the corresponding sender.
func (dag *Dag) QueryTransactionsWithSender(sender common.Address) ([]*types.Transaction, error) {
	leaves, err := dag.Root.GetChildrenBySender(sender) // Get child
	if err != nil {                                     // Check for errors
		return []*types.Transaction{}, err // Return found error
	}

	if dag.Root.Transaction != nil && dag.Root.Transaction.Sender != nil && bytes.Equal(dag.Root.Transaction.Sender[:], sender[:]) { // Check root matches
		leaves = append(leaves, dag.Root) // Append root to leaves
	}

	transactions := []*types.Transaction{} // Init tx list buffer

	for _, leaf := range leaves { // Iterate through leaves
		if len(transactions) != 0 { // Check matches already exist
			alreadyExists := false // Init duplicate match buffer

			for _, match := range transactions { // Iterate through matches
				if leaf.Transaction.Hash != nil { // Check has hash
					if bytes.Equal(match.Hash[:], leaf.Hash[:]) { // Check duplicate match
						alreadyExists = true // Set true

						break // Break
					}
				}
			}

			if alreadyExists { // Check already exists
				continue // Continue
			}
		}

		transactions = append(transactions, leaf.Transaction) // Append tx
	}

	return transactions, nil // Return leaves
}

// QueryTransactionsWithRecipient queries the dag for a list of transactions with the corresponding recipient.
func (dag *Dag) QueryTransactionsWithRecipient(recipient common.Address) ([]*types.Transaction, error) {
	leaves, err := dag.Root.GetChildrenByRecipient(recipient) // Get child
	if err != nil {                                           // Check for errors
		return []*types.Transaction{}, err // Return found error
	}

	if dag.Root.Transaction != nil && dag.Root.Transaction.Recipient != nil && bytes.Equal(dag.Root.Transaction.Recipient[:], recipient[:]) { // Check root matches
		leaves = append(leaves, dag.Root) // Append root to leaves
	}

	transactions := []*types.Transaction{} // Init tx list buffer

	for _, leaf := range leaves { // Iterate through leaves
		if len(transactions) != 0 { // Check matches already exist
			alreadyExists := false // Init duplicate match buffer

			for _, match := range transactions { // Iterate through matches
				if leaf.Transaction.Hash != nil { // Check has hash
					if bytes.Equal(match.Hash[:], leaf.Hash[:]) { // Check duplicate match
						alreadyExists = true // Set true

						break // Break
					}
				}
			}

			if alreadyExists { // Check already exists
				continue // Continue
			}
		}

		transactions = append(transactions, leaf.Transaction) // Append tx
	}

	return transactions, nil // Return leaves
}

// QueryLeafWithHash queries the dag for a leaf with the corresponding hash.
func (dag *Dag) QueryLeafWithHash(hash common.Hash) (*Leaf, error) {
	return dag.Root.GetChildByHash(hash) // Return leaf
}

// QueryNextCommonLeaf attempts to find the next common leaf in the dag.
// A common leaf is defined as a leaf that has no siblings.
// If no common leaf exists, an error is returned.
func (dag *Dag) QueryNextCommonLeaf(lastCommonLeaf *Leaf) (*Leaf, error) {
	return lastCommonLeaf.GetNextCommonLeaf() // Get next common leaf
}

// END QUERY HELPERS

// BEGIN NETWORK HELPERS

// MakeGenesis constructs a set of genesis transactions from the given config
// and adds them to the working dag.
func (dag *Dag) MakeGenesis(config *config.ChainConfig) error {
	totalGenesisValue := big.NewFloat(0) // Init total genesis value

	for _, allocAddress := range config.AllocAddresses { // Iterate through alloc addresses
		totalGenesisValue.Add(totalGenesisValue, config.Alloc[allocAddress.String()]) // Add value
	}

	genesisAccount, err := accounts.NewAccount() // Initialize new account
	if err != nil {                              // Check for errors
		return err // Return found error
	}

	genesisTransaction, err := types.NewTransaction(0, nil, nil, &genesisAccount.Address, totalGenesisValue, []byte("genesis")) // Initialize genesis transaction
	if err != nil {                                                                                                             // Check for errors
		return err // Return found error
	}

	err = dag.AddTransaction(genesisTransaction) // Add genesis tx
	if err != nil {                              // Check for errors
		return err // Return found error
	}

	lastTransaction := genesisTransaction // Set last tx pointer

	for _, address := range config.AllocAddresses { // Iterate through alloc addresses
		transaction, err := types.NewTransaction(0, lastTransaction, &genesisAccount.Address, &address, config.Alloc[address.String()], []byte("genesis_child")) // Init tx
		if err != nil {                                                                                                                                          // Check for errors
			return err // Return found error
		}

		err = types.SignTransaction(transaction, genesisAccount.PrivateKey) // Sign tx
		if err != nil {                                                     // Check for errors
			return err // Return found error
		}

		err = dag.AddTransaction(transaction) // Add genesis child tx
		if err != nil {                       // Check for errors
			return err // Return found error
		}

		lastTransaction = transaction // Set last tx
	}

	return nil // No error occurred, return nil
}

// END NETWORK HELPERS

/*
	END HELPER METHODS
*/

/* END EXPORTED METHODS */
