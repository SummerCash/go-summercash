package types

import (
	"github.com/space55/summertech-blockchain/common"
)

// Transaction - primitive transaction type
type Transaction struct {
	AccountNonce uint64 `json:"nonce"` // Nonce in set of account transactions

	Amount float64 `json:"amount"` // Amount of coins sent in transaction

	Sender    *common.Address `json:"sender"`    // Transaction sender
	Recipient *common.Address `json:"recipient"` // Transaction recipient

	Signature *signature `json:"signature"` // Transaction signature meta
}

/* BEGIN EXPORTED METHODS */

/* END EXPORTED METHODS */
