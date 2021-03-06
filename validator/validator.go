// Package validator represents a collection of helper methods useful for validators in the SummerCash network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
// NOTE: Ripped from another one of my projects; https://github.com/polaris-project/go-polaris/tree/master/validator
package validator

import (
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/types"
)

// Validator represents any generic validator.
type Validator interface {
	ValidateTransaction(transaction *types.Transaction) error // Validate a given transaction

	ValidateTransactionHash(transaction *types.Transaction) bool // Validate a given transaction's hash

	ValidateTransactionState(transaction *types.Transaction) bool // Validate a given transaction's state

	ValidateTransactionTimestamp(transaction *types.Transaction) bool // Validate a given transaction's timestamp

	ValidateTransactionSignature(transaction *types.Transaction) bool // Validate a given transaction's signature

	ValidateTransactionSenderBalance(transaction *types.Transaction) bool // Validate a given transaction's sender has

	ValidateTransactionIsNotDuplicate(transaction *types.Transaction) bool // Validate that a given transaction does not already exist in the working chain

	// ValidateTransactionReward(transaction *types.Transaction) bool // Validate a given transaction reward

	ValidateTransactionNonce(transaction *types.Transaction) bool // Validate that a given transaction's nonce is equivalent to the current account index + 1

	ValidationProtocol() string // Get the current validator's validation protocol

	GetWorkingConfig() *config.ChainConfig // Get current validator's working config
}
