package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/SummerCash/go-summercash/common"
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

	json, err := json.MarshalIndent(*chain, "", "  ") // Marshal chain

	if err != nil { // Check for errors
		return err // Return error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/db/chain/chain_%s.json", common.DataDir, chain.Account.String())), json, 0644) // Write JSON

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
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/db/chain/chain_%s.json", common.DataDir, address.String()))) // Read chain

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	buffer := &Chain{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Read json into buffer

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	err = buffer.RecoverSafeEncoding() // Recover

	if err != nil { // Check for errors
		return &Chain{}, err // Return error
	}

	return buffer, nil // No error occurred, return read coordinationChain
}

/* END EXPORTED METHODS */
