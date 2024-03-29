// Package validator represents a collection of helper methods useful for validators in the SummerCash network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
// NOTE: Ripped from another one of my projects; https://github.com/polaris-project/go-polaris/tree/master/validator
package validator

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/SummerCash/ursa/compiler"

	"github.com/SummerCash/go-summercash/common"

	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/crypto"
	"github.com/SummerCash/go-summercash/types"
)

// StandardValidatorValidationProtocol represents the validation protocol of the standard validator.
const StandardValidatorValidationProtocol = "standard_sig_ver"

var (
	// ErrInvalidTransactionHash is an error definition representing a transaction hash of invalid value.
	ErrInvalidTransactionHash = errors.New("transaction hash is invalid")

	// ErrInvalidTransactionState is an error definition representing a transaction state of invalid value.
	ErrInvalidTransactionState = errors.New("transaction state is invalid")

	// ErrInvalidTransactionTimestamp is an error definition representing a transaction timestamp of invalid value.
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
}

/* BEGIN EXPORTED METHODS */

// NewStandardValidator initializes a new beacon dag with a given config and working chain.
func NewStandardValidator(config *config.ChainConfig) *StandardValidator {
	return &StandardValidator{
		Config: config, // Set config
	}
}

// ValidateTransaction validates the given transaction via the standard validator.
// Each validation issue is returned as an error.
func (validator *StandardValidator) ValidateTransaction(transaction *types.Transaction) error {
	err := validator.PerformChainSafetyChecks(transaction) // Perform safety checks
	if err != nil {                                        // Check for errors
		return err // Return found error
	}

	if !validator.ValidateTransactionHash(transaction) { // Check invalid hash
		return ErrInvalidTransactionHash // Invalid hash
	}

	if !validator.ValidateTransactionState(transaction) { // Check invalid state
		return ErrInvalidTransactionState // Invalid state
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

// PerformChainSafetyChecks loads a given transaction's sender chain, requests it if it doesn't exist,
// and makes one if it cannot request it from its peers.
func (validator *StandardValidator) PerformChainSafetyChecks(transaction *types.Transaction) error {
	_, err := types.ReadChainFromMemory(*transaction.Sender) // Read sender chain
	if err != nil {                                          // Check for errors
		_, err := types.NewChain(*transaction.Sender) // Initialize chain
		if err != nil {                               // Check for errors
			return err // Return found error
		}
	}

	_, err = types.ReadChainFromMemory(*transaction.Recipient) // Read sender chain

	if err != nil { // Check for errors
		_, err := types.NewChain(*transaction.Recipient) // Initialize chain
		if err != nil {                                  // Check for errors
			return err // Return found error
		}
	}

	return nil // No error occurred, return nil
}

// ValidateTransactionHash checks that a given transaction's hash is equivalent to the calculated hash of that given transaction.
func (validator *StandardValidator) ValidateTransactionHash(transaction *types.Transaction) bool {
	unsignedTx := *transaction // Init unsigned buffer

	unsignedTx.Signature = nil // Set signature to nil
	unsignedTx.Hash = nil      // Set hash to nil

	return bytes.Equal(transaction.Hash.Bytes(), common.NewHash(crypto.Sha3(unsignedTx.Bytes())).Bytes()) // Return hashes equivalent
}

// ValidateTransactionState checks that a given transaction's working state is equivalent to that of the parent, with the given
// transition applied.
func (validator *StandardValidator) ValidateTransactionState(transaction *types.Transaction) bool {
	if transaction.State == nil { // Check no state
		return true // Valid
	}

	gasPolicy := compiler.GasPolicy(common.GasPolicy) // Get gas policy

	localState, err := transaction.EvaluateNewState(&gasPolicy) // Calculate local state
	if err != nil {                                             // Check for errors
		return false // Invalid
	}

	if !bytes.Equal(localState.Bytes(), transaction.State.Bytes()) { // Check invalid state
		return false // Invalid
	}

	return true // Valid
}

// ValidateTransactionTimestamp validates the given transaction's timestamp against that of its parents.
// If the timestamp of any one of the given transaction's parents is after the given transaction's timestamp, false is returned.
// If any one of the transaction's parent transactions cannot be found in the working dag, false is returned.
func (validator *StandardValidator) ValidateTransactionTimestamp(transaction *types.Transaction) bool {
	senderChain, err := types.ReadChainFromMemory(*transaction.Sender) // Read sender chain
	if err != nil {                                                    // Check for errors
		return false // Invalid
	}

	if len(senderChain.Transactions) == 0 { // Check cannot compare timestamps
		return true // Valid timestamp
	}

	if senderChain.Transactions[len(senderChain.Transactions)-1].Timestamp.After(transaction.Timestamp) { // Check invalid timestamp
		return false // Invalid timestamp
	}

	return true // Valid timestamp
}

// ValidateTransactionSignature validates the given transaction's signature against the transaction sender's public key.
// If the transaction's signature is nil, false is returned.
func (validator *StandardValidator) ValidateTransactionSignature(transaction *types.Transaction) bool {
	if transaction.Signature == nil { // Check has no signature
		return false // Nil signature
	}

	valid, err := types.VerifyTransactionSignature(transaction) // Check valid
	if err != nil {                                             // Check for errors
		return false // Invalid
	}

	return valid // Return signature validity
}

// ValidateTransactionSenderBalance checks that a given transaction's sender has a balance greater than or equal to the transaction's total value (including gas costs).
func (validator *StandardValidator) ValidateTransactionSenderBalance(transaction *types.Transaction) bool {
	if transaction.Amount.Cmp(big.NewFloat(0)) == -1 { // Check is negative number
		return false // Invalid tx amount
	}

	chain, err := types.ReadChainFromMemory(*transaction.Sender) // Read sender chain
	if err != nil {
		return false // Invalid
	}

	balance := chain.CalculateBalance() // Calculate balance

	return balance.Cmp(transaction.Amount) == 0 || balance.Cmp(transaction.Amount) == 1 // Return sender balance adequate
}

// ValidateTransactionIsNotDuplicate checks that a given transaction does not already exist in the working dag.
func (validator *StandardValidator) ValidateTransactionIsNotDuplicate(transaction *types.Transaction) bool {
	chain, err := types.ReadChainFromMemory(*transaction.Sender) // Read sender chain
	if err != nil {
		return false // Invalid
	}

	_, err = chain.QueryTransaction(*transaction.Hash) // Attempt to get tx by hash

	if err == nil { // Check transaction exists
		return false // Transaction is duplicate
	}

	return true // Transaction is unique
}

// ValidateTransactionNonce checks that a given transaction's nonce is equivalent to the sending account's last nonce + 1.
func (validator *StandardValidator) ValidateTransactionNonce(transaction *types.Transaction) bool {
	chain, err := types.ReadChainFromMemory(*transaction.Sender) // Read sender chain
	if err != nil {
		return false // Invalid
	}

	return chain.CalculateTargetNonce() == transaction.AccountNonce // Return nonce valid
}

// ValidationProtocol fetches the current validator's validation protocol.
func (validator *StandardValidator) ValidationProtocol() string {
	return StandardValidatorValidationProtocol // Return validation protocol
}

// GetWorkingConfig attempts to fetch the working config instance.
func (validator *StandardValidator) GetWorkingConfig() *config.ChainConfig {
	return validator.Config // Return working config
}

/* END EXPORTED METHODS */
