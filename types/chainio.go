package types

import (
	"fmt"
	"path/filepath"

	"github.com/space55/summertech-blockchain/common"
)

/* BEGIN EXPORTED METHODS */

// WriteToMemory - write given chain to memory
func (chain *Chain) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/db/chain", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	err = chain.MakeEncodingSafe() // Make encoding safe

	if err != nil { // Check for errors
		return err // Return error
	}

	err = common.WriteGob(fmt.Sprintf("%s/db/chain/chain_%s.gob", common.DataDir, chain.Account.String()), chain) // Write gob

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = chain.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadChainFromMemory - read chain from memory
func ReadChainFromMemory(address common.Address) (*Chain, error) {
	chain := &Chain{} // Init buffer

	err := common.ReadGob(filepath.FromSlash(fmt.Sprintf("%s/db/chain/chain_%s.gob", common.DataDir, address.String())), chain) // Read chain

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	err = chain.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	return chain, nil // No error occurred, return read coordinationChain
}

/* END EXPORTED METHODS */
