// Package types implements many of the core go-summercash types.
package types

import (
	"github.com/SummerCash/go-summercash/crypto"
	"math/big"
)

// Transaction defines a generic transaction.
type Transaction struct {
	Nonce uint64 `json:"nonce"` // Index of transaction in account tx set

	GasLimit uint64 `json:"gas_limit"` // Number of finks willing to spend on gas

	Amount *big.Int `json:"amount"` // Number of finks transaction is worth (1 SMC = 1,000,000,000,000,000,000 finks)

	Reward *big.Int `json:"reward_amount"` // Reward amount (in finks, of course)

	Payload []byte `json:"payload"` // Transaction payload

	Signature *crypto.Signature `json:"signature"` // Transaction signature

	Hash crypto.Hash `json:"hash"` // Transaction hash

	PreviousStateHash crypto.Hash `json:"pre_state_hash"` // Hash of contract state at previous transaction

	Parents []crypto.Hash `json:"parents"` // Transaction parents
}

/* BEGIN EXPORTED METHODS */

// NewTransaction initializes a new transaction with a set of parameters.
func NewTransaction(nonce uint64, gasLimit uint64, amount *big.Int, payload []byte, lastTransaction *Transaction) (*Transaction, error) {
}

/* END EXPORTED METHODS */
