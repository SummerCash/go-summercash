package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/space55/summertech-blockchain/common"
)

// WriteToMemory - write given coordination chain to memory
func (coordinationChain *CoordinationChain) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/db/coordination_chain", common.DataDir)) // Create dir if necessary

	if err != nil { // Check for errors
		return err // Return error
	}

	json, err := json.MarshalIndent(*coordinationChain, "", "  ") // Marshal coordination chain

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/db/coordination_chain/chain.gob", common.DataDir), json, 0644) // Write json

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// ReadCoordinationChainFromMemory - read coordinationChain from memory
func ReadCoordinationChainFromMemory() (*CoordinationChain, error) {
	coordinationChain := &CoordinationChain{} // Init buffer

	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/db/coordination_chain/chain.gob", common.DataDir))) // Read file

	if err != nil { // Check for errors
		return &CoordinationChain{}, err // Return error
	}

	err = json.Unmarshal(data, coordinationChain) // Read json into buffer

	if err != nil { // Check for errors
		return &CoordinationChain{}, err // Return error
	}

	return coordinationChain, nil // No error occurred, return read coordinationChain
}
