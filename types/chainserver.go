package types

import (
	"encoding/hex"
	"errors"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/ursa/vm"
)

var (
	// ErrBadChain - error describing input chain with tx length shorter than current
	ErrBadChain = errors.New("chain out of date")
)

/* BEGIN EXPORTED METHODS */

// HandleReceivedChainRequest - handle chain request
func HandleReceivedChainRequest(b []byte) (*Chain, error) {
	var address common.Address // Init buffer

	copy(address[:], b[12:][:]) // Copy read address

	chain, err := ReadChainFromMemory(address) // Read chain from memory

	if err != nil { // Check for errors
		return &Chain{}, err // Return found error
	}

	return chain, nil // Return read chain
}

// HandleReceivedTransaction - handle received transaction
func HandleReceivedTransaction(b []byte) error {
	tx, err := TransactionFromBytes(b) // Marshal bytes to transaction

	if err != nil { // Check for errors
		return err // Return error
	}

	tx.RecoverSafeEncoding() // Recover from safe encoding

	chain, err := ReadChainFromMemory(*tx.Recipient) // Read chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	oldNonce := tx.AccountNonce // Set old nonce

	tx.AccountNonce = uint64(len(chain.Transactions)) // Reset nonce

	common.Logf("== CHAIN == adding transaction %s to chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx

	err = chain.AddTransaction(tx) // Append tx

	if err != nil { // Check for errors
		common.Logf("== ERROR == error adding transaction to chain %s\n", err.Error()) // Log error

		return err // Return found error
	}

	common.Logf("== SUCCESS == added transaction %s to chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx

	if tx.Sender != nil { // Check has sender
		chain, err = ReadChainFromMemory(*tx.Sender) // Read tx sender chain

		if err != nil { // Check for errors
			return err // Return found error
		}

		common.Logf("== CHAIN == adding transaction %s to sender chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx

		tx.AccountNonce = oldNonce // Set to old nonce

		err = chain.AddTransaction(tx) // Append tx

		if err != nil { // Check for errors
			common.Logf("== ERROR == error adding transaction to sender chain %s\n", err.Error()) // Log error

			return err // Return found error
		}

		common.Logf("== SUCCESS == added transaction %s to sender chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx
	}

	return nil // No error occurred, return nil
}

// HandleReceivedContractCreation - handle received transaction creating contract
func HandleReceivedContractCreation(b []byte) error {
	tx, err := TransactionFromBytes(b) // Marshal bytes to transaction

	if err != nil { // Check for errors
		return err // Return error
	}

	tx.RecoverSafeEncoding() // Recover from safe encoding

	chain, err := ReadChainFromMemory(*tx.Recipient) // Read chain

	if err != nil { // Check for errors
		return err // Return found error
	}

	oldNonce := tx.AccountNonce // Set old nonce

	tx.AccountNonce = uint64(len(chain.Transactions)) // Reset nonce

	common.Logf("== CHAIN == adding transaction %s to chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx

	err = chain.AddTransaction(tx) // Append tx

	if err != nil { // Check for errors
		common.Logf("== ERROR == error adding transaction to chain %s\n", err.Error()) // Log error

		return err // Return found error
	}

	common.Logf("== SUCCESS == added transaction %s to chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx

	if tx.Sender != nil { // Check has sender
		chain, err = ReadChainFromMemory(*tx.Sender) // Read tx sender chain

		if err != nil { // Check for errors
			return err // Return found error
		}

		common.Logf("== CHAIN == adding transaction %s to sender chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx

		tx.AccountNonce = oldNonce // Set to old nonce

		err = chain.AddTransaction(tx) // Append tx

		if err != nil { // Check for errors
			common.Logf("== ERROR == error adding transaction to sender chain %s\n", err.Error()) // Log error

			return err // Return found error
		}

		common.Logf("== SUCCESS == added transaction %s to sender chain %s\n", tx.Hash.String(), chain.ID.String()) // Log add tx
	}

	return nil // No error occurred, return nil
}

// HandleReceivedStateRequest - handle received contract state request
func HandleReceivedStateRequest(b []byte) ([]byte, error) {
	var address common.Address // Init buffer

	copy(address[:], b[9:][:]) // Copy read address

	chain, err := ReadChainFromMemory(address) // Read chain

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	env, err := vm.ReadEnvironmentFromMemory() // Read environment from memory

	if err != nil { // Check for errors
		env = &common.VMConfig // Get VM config

		err = env.WriteToMemory() // Write to persistent memory

		if err != nil { // Check for errors
			return nil, err // Return found error
		}
	}

	virtualMachine, err := vm.NewVirtualMachine(chain.ContractSource, *env, new(vm.Resolver), nil) // Init vm

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	state, err := vm.ReadStateDBFromMemory(hex.EncodeToString(virtualMachine.StateDB.ID)) // Read state DB

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return state.Bytes(), nil // Return read state
}

// HandleReceivedChain - handle received chain
func HandleReceivedChain(b []byte) error {
	chain, err := FromBytes(b) // Marshal bytes to chain

	if err != nil { // Check for errors
		return err // Return error
	}

	oldChain, err := ReadChainFromMemory(chain.Account) // Check for conflicts

	if err == nil { // Check for errors
		if len(oldChain.Transactions) > len(chain.Transactions) { // Check bad chain
			return ErrBadChain // Return error
		}
	}

	err = chain.WriteToMemory() // Write chain to memory

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
