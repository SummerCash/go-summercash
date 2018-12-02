package types

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/space55/summertech-blockchain/common"
	"github.com/space55/summertech-blockchain/crypto"
)

var (
	// ErrAlreadySigned - error definition stating transaction has already been signed
	ErrAlreadySigned = errors.New("transaction already signed")

	// ErrNilSignature - error definition describing nil tx signature
	ErrNilSignature = errors.New("nil signature")

	// ErrInvalidSignature - error definition describing invalid tx signature (doesn't match public key)
	ErrInvalidSignature = errors.New("invalid signature")
)

// Transaction - primitive transaction type
type Transaction struct {
	AccountNonce uint64 `json:"nonce"` // Nonce in set of account transactions

	Sender    *common.Address `json:"sender"`    // Transaction sender
	Recipient *common.Address `json:"recipient"` // Transaction recipient

	Amount float64 `json:"amount"` // Amount of coins sent in transaction

	Payload []byte `json:"payload"` // Misc. data transported with transaction

	Signature *signature `json:"signature"` // Transaction signature meta

	Hash *common.Hash `json:"hash"` // Transaction hash
}

/* BEGIN EXPORTED METHODS */

// NewTransaction - attempt to initialize transaction primitive
func NewTransaction(nonce uint64, sender *common.Address, destination *common.Address, amount float64, payload []byte) (*Transaction, error) {
	transaction := Transaction{ // Init tx
		AccountNonce: nonce,       // Set nonce
		Sender:       sender,      // Set sender
		Recipient:    destination, // Set recipient
		Amount:       amount,      // Set amount
		Payload:      payload,     // Set tx payload
	}

	hash := common.NewHash(crypto.Sha3(transaction.Bytes())) // Hash transaction

	transaction.Hash = &hash // Set hash

	return &transaction, nil // Return initialized transaction
}

// Bytes - convert given transaction to byte array
func (tx *Transaction) Bytes() []byte {
	buffer := new(bytes.Buffer) // Init buffer

	json.NewEncoder(buffer).Encode(*tx) // Serialize tx

	return buffer.Bytes() // Return serialized
}

/* END EXPORTED METHODS */
