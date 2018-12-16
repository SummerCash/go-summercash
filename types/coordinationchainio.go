package types

import (
	"fmt"

	"github.com/space55/summertech-blockchain/common"
)

// WriteToMemory - write given coordination chain to memory
func (coordinationChain *CoordinationChain) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/db/coordination_chain", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	err = common.WriteGob(fmt.Sprintf("%s/db/coordination_chain/chain.gob", common.DataDir), *coordinationChain) // Write gob

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ReadCoordinationChainFromMemory - read coordinationChain from memory
func ReadCoordinationChainFromMemory() (*CoordinationChain, error) {
	coordinationChain := &CoordinationChain{} // Init buffer

	err := common.ReadGob(fmt.Sprintf("%s/db/coordination_chain/chain.gob", common.DataDir), coordinationChain) // Read chain

	if err != nil { // Check for errors
		return &CoordinationChain{}, err // Return error
	}

	return coordinationChain, nil // No error occurred, return read coordinationChain
}