// Package validator represents a collection of helper methods useful for validators in the SummerCash network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
// NOTE: Ripped from another one of my projects; https://github.com/polaris-project/go-polaris/tree/master/validator
package validator

import (
	"bytes"
	"errors"

	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"
	"github.com/SummerCash/go-summercash/types"
)

var (
	// ErrInvalidTransactionHash is an error definition representing a transaction hash of invalid value.
	ErrInvalidTransactionHash = errors.New("transaction hash is invalid")

	//ErrInvalidTransactionTimestamp is an error definition representing a transaction timestamp of invalid value.
	ErrInvalidTransactionTimestamp = errors.New("invalid transaction timestamp")

	// ErrInvalidTransactionSignature is an error definition representing a transaction signature of invalid value.
	ErrInvalidTransactionSignature = errors.New("invalid transaction signature")

	// ErrInsufficientSenderBalance is an error definition representing a sender balance of insufficient value.
	ErrInsufficientSenderBalance = errors.New("insufficient sender balance")

	// ErrDuplicateTransaction is an error definition representing a transaction of duplicate value in the working chain.
	ErrDuplicateTransaction = errors.New("transaction already exists in the working chain (duplicate)")

	// ErrInvalidNonce is an error definition representing a transaction of invalid nonce value.
	ErrInvalidNonce = errors.New("invalid transaction nonce")
)

// StandardValidator represents a standard validator implementing the validator interface.
type StandardValidator struct {
	Config *config.ChainConfig `json:"config"` // Chain configuration reference

	WorkingChain *types.Chain `json:"work_chain"` // Working chain reference
}

/* BEGIN EXPORTED METHODS */

// NewStandardValidator initializes a new beacon dag with a given config and working chain.
func NewStandardValidator(config *config.ChainConfig, workingChain *types.Chain) *StandardValidator {
	return &StandardValidator{
		Config:       config,       // Set config
		WorkingChain: workingChain, // Set working chain
	}
}

// ValidateTransaction validates the given transaction via the standard validator.
// Each validation issue is returned as an error.
func (validator *StandardValidator) ValidateTransaction(transaction *types.Transaction) error {
	if !validator.ValidateTransactionHash(transaction) { // Check invalid hash
		return ErrInvalidTransactionHash // Invalid hash
	}

	if !validator.ValidateTransactionTimestamp(transaction) { // Check invalid timestamp
		return ErrInvalidTransactionTimestamp // Invalid timestamp
	}

	if !validator.ValidateTransactionSignature(transaction) { // Check invalid signature
		return ErrInvalidTransactionSignature // Invalid signature
	}

	if !validator.ValidateTransactionSenderBalance(transaction) { // Check invalid value
		return ErrInsufficientSenderBalance // Invalid value
	}

	if !validator.ValidateTransactionIsNotDuplicate(transaction) { // Check duplicate
		return ErrDuplicateTransaction // Duplicate
	}

	if !validator.ValidateTransactionNonce(transaction) { // Check valid nonce
		return ErrInvalidNonce // Invalid nonce
	}

	return nil // Transaction is valid
}

// ValidateTransactionHash checks that a given transaction's hash is equivalent to the calculated hash of that given transaction.
func (validator *StandardValidator) ValidateTransactionHash(transaction *types.Transaction) bool {
	unsignedTx := transaction // Init unsigned buffer

	(*unsignedTx).Signature = nil // Set signature to nil

	return bytes.Equal(transaction.Hash.Bytes(), crypto.Sha3(unsignedTx.Bytes())) // Return hashes equivalent
}

// ValidateTransactionTimestamp validates the given transaction's timestamp against that of its parents.
// If the timestamp of any one of the given transaction's parents is after the given transaction's timestamp, false is returned.
// If any one of the transaction's parent transactions cannot be found in the working dag, false is returned.
func (validator *StandardValidator) ValidateTransactionTimestamp(transaction *types.Transaction) bool {
	for _, parentHash := range transaction.ParentTransactions { // Iterate through parent hashes
		parentTransaction, err := validator.WorkingDag.GetTransactionByHash(parentHash) // Get parent transaction pointer

		if err != nil { // Check for errors
			return false // Invalid parent
		}

		if parentTransaction.Timestamp.After(transaction.Timestamp) {
			return false // Invalid timestamp
		}
	}

	return true // Valid timestamp
}

// ValidateTransactionSignature validates the given transaction's signature against the transaction sender's public key.
// If the transaction's signature is nil, false is returned.
func (validator *StandardValidator) ValidateTransactionSignature(transaction *types.Transaction) bool {
	if transaction.Signature == nil { // Check has no signature
		return false // Nil signature
	}

	return transaction.Signature.Verify(transaction.Sender) // Return signature validity
}

// ValidateTransactionSenderBalance checks that a given transaction's sender has a balance greater than or equal to the transaction's total value (including gas costs).
func (validator *StandardValidator) ValidateTransactionSenderBalance(transaction *types.Transaction) bool {
	balance, err := validator.WorkingDag.CalculateAddressBalance(transaction.Sender) // Calculate balance

	if err != nil { // Check for errors
		return false // Invalid
	}

	return balance.Cmp(transaction.CalculateTotalValue()) == 0 || balance.Cmp(transaction.CalculateTotalValue()) == 1 // Return sender balance adequate
}

// ValidateTransactionIsNotDuplicate checks that a given transaction does not already exist in the working dag.
func (validator *StandardValidator) ValidateTransactionIsNotDuplicate(transaction *types.Transaction) bool {
	transaction, err := validator.WorkingDag.GetTransactionByHash(transaction.Hash) // Attempt to get tx by hash

	if err == nil && !transaction.Hash.IsNil() { // Check transaction exists
		return false // Transaction is duplicate
	}

	return true // Transaction is unique
}

// ValidateTransactionDepth checks that a given transaction's parent hash is a member of the last edge.
func (validator *StandardValidator) ValidateTransactionDepth(transaction *types.Transaction) bool {
	for _, parentHash := range transaction.ParentTransactions { // Iterate through parent hashes
		if bytes.Equal(parentHash.Bytes(), transaction.Hash.Bytes()) { // Check self in parent hashes
			return false // Invalid
		}

		children, err := validator.WorkingDag.GetTransactionChildren(parentHash) // Get children of transaction

		if err != nil { // Check for errors
			return false // Invalid
		}

		for _, child := range children { // Iterate through children
			currentChildren, err := validator.WorkingDag.GetTransactionChildren(child.Hash) // Get children of current child

			if err != nil { // Check for errors
				return false // Invalid
			}

			if len(currentChildren) != 0 { // Check child has children
				return false // Invalid depth
			}
		}
	}

	return true // Valid
}

// ValidateTransactionNonce checks that a given transaction's nonce is equivalent to the sending account's last nonce + 1.
func (validator *StandardValidator) ValidateTransactionNonce(transaction *types.Transaction) bool {
	senderTransactions, err := validator.WorkingDag.GetTransactionsBySender(transaction.Sender) // Get sender txs

	if err != nil { // Check for errors
		return false // Invalid
	}

	if len(senderTransactions) == 0 { // Check is genesis
		if transaction.AccountNonce != 0 { // Check nonce is not 0
			return false // Invalid nonce
		}

		return true // Valid nonce
	}

	lastNonce := uint64(0) // Init nonce buffer

	for _, currentTransaction := range senderTransactions { // Iterate through sender txs
		if currentTransaction.AccountNonce > lastNonce { // Check greater than last nonce
			lastNonce = currentTransaction.AccountNonce // Set last nonce
		}
	}

	if transaction.AccountNonce != lastNonce+1 { // Check invalid nonce
		return false // Invalid nonce
	}

	return true // Valid nonce
}

// ValidationProtocol fetches the current validator's validation protocol.
func (validator *StandardValidator) ValidationProtocol() string {
	return StandardValidatorValidationProtocol // Return validation protocol
}

// GetWorkingDag attempts to fetch the working dag instance.
func (validator *StandardValidator) GetWorkingDag() *types.Dag {
	return validator.WorkingDag // Return working dag
}

// GetWorkingConfig attempts to fetch the working config instance.
func (validator *StandardValidator) GetWorkingConfig() *config.DagConfig {
	return validator.Config // Return working config
}

/* END EXPORTED METHODS */
